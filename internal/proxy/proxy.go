package proxy

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gobenpark/go_thought/internal/log"
	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type LLMParser interface {
	ParseRequest(r *http.Request) error
	ParseResponse(resp *http.Response, writer http.ResponseWriter) error
}

type Config struct {
	Port           int    `yaml:"port"`
	TargetURL      string `yaml:"target_url"`
	LogRequests    bool   `yaml:"log_requests"`
	RequestTimeout string `yaml:"request_timeout"`
	debug          bool   `yaml:"debug"`
}

// ProxyServer represents the proxy server
type ProxyServer struct {
	config     Config
	httpServer *http.Server
	client     *http.Client
	logger     log.Logger
}

func NewProxyServer(config Config) *ProxyServer {

	if config.Port == 0 {
		config.Port = 8080
	}

	server := &ProxyServer{
		config: config,
		httpServer: &http.Server{
			Addr: fmt.Sprintf(":%d", config.Port),
		},
	}

	if config.RequestTimeout == "" {
		duration, err := time.ParseDuration(config.RequestTimeout)
		if err != nil {
			fmt.Println(err)
			duration = time.Second * 30
		}

		server.client = &http.Client{
			Timeout: duration,
		}
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

	p.logger.Debug("receive url", "url", r.URL.String())
	if strings.Contains(r.Header.Get("User-Agent"), "OpenAI") {
		r.URL.Scheme = "https"
		r.URL.Host = "api.openai.com"
		r.URL.Path = "/v1" + r.URL.Path
		p.logger.Debug(r.URL.String())
	}

	//TODO: Request Parse

	proxyReq, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		p.logger.Error("Error creating proxy request", err)
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
			p.logger.Debug("request header", name, value)
			proxyReq.Header.Add(name, value)
		}
	}

	resp, err := p.client.Do(proxyReq)
	if err != nil {
		p.logger.Error("proxy client error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
			p.logger.Debug("response headers", name, value)
		}
	}

	p.logger.Debug("response", "url", r.URL.String(), "statusCode", resp.StatusCode, "statusText", http.StatusText(resp.StatusCode))
	w.WriteHeader(resp.StatusCode)

	//TODO: Response parse
}

func (p *ProxyServer) Start() error {
	p.httpServer.Handler = p

	return p.httpServer.ListenAndServe()
}

func (p *ProxyServer) Shutdown(ctx context.Context) error {
	return p.httpServer.Shutdown(ctx)
}
