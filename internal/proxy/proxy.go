package proxy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gobenpark/go_thought/internal/log"
)

type Config struct {
	Logger         log.Logger
	Port           int
	TargetURL      string
	LogRequests    bool
	RequestTimeout time.Duration
}

// ProxyServer represents the proxy server
type ProxyServer struct {
	config     Config
	httpServer *http.Server
	client     *http.Client
}

func NewProxyServer(config Config) *ProxyServer {

	if config.Port == 0 {
		config.Port = 8080
	}

	if config.RequestTimeout == 0 {
		config.RequestTimeout = 30 * time.Second
	}

	if config.Logger == nil {
		config.Logger = log.NewZapLogger()
	}

	return &ProxyServer{
		config: config,
		httpServer: &http.Server{
			Addr: fmt.Sprintf(":%d", config.Port),
		},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (p *ProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	p.config.Logger.Debug("receive url", "url", r.URL.String())
	re := regexp.MustCompile("OpenAI*")
	if re.MatchString(r.Header.Get("User-Agent")) {
		r.URL.Scheme = "https"
		r.URL.Host = "api.openai.com"
		r.URL.Path = "/v1" + r.URL.Path
		p.config.Logger.Debug(r.URL.String())
	}

	buf := &bytes.Buffer{}
	_, err := io.Copy(buf, r.Body)
	if err != nil {
		p.config.Logger.Error("copy error", err)
		return
	}

	requestBodyCopy := io.NopCloser(buf)
	r.Body = io.NopCloser(buf)

	proxyReq, err := http.NewRequest(r.Method, r.URL.String(), requestBodyCopy)
	if err != nil {
		p.config.Logger.Error("Error creating proxy request", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if clientIP, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		if prior, ok := proxyReq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		proxyReq.Header.Set("X-Forwarded-For", clientIP)
	}

	for name, values := range r.Header {
		for _, value := range values {
			p.config.Logger.Debug(name, "value", value)
			proxyReq.Header.Add(name, value)
		}
	}

	resp, err := p.client.Do(proxyReq)
	if err != nil {
		p.config.Logger.Error("proxy client error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	p.config.Logger.Debug("response", "url", r.URL.String(), "statusCode", resp.StatusCode, "statusText", http.StatusText(resp.StatusCode))
	for name, values := range resp.Header {
		for _, value := range values {
			p.config.Logger.Debug("headers", name, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		p.config.Logger.Error("proxy client error", err)
		return
	}
}

func (p *ProxyServer) Start() error {
	p.httpServer.Handler = p

	return p.httpServer.ListenAndServe()
}

func (p *ProxyServer) Shutdown(ctx context.Context) error {
	return p.httpServer.Shutdown(ctx)
}
