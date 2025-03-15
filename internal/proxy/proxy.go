package proxy

import (
	"fmt"
	"io"
	"net"
	"net/http"
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
	return &ProxyServer{
		config:     config,
		httpServer: nil,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (p *ProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	proxyReq, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		//TODO: Logger
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
			proxyReq.Header.Add(name, value)
		}
	}

	resp, err := p.client.Do(proxyReq)
	if err != nil {
		//TODO: Logger
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	//p.logger.Printf("[RESP:%s] %s %d %s", requestID, r.URL.String(), resp.StatusCode, http.StatusText(resp.StatusCode))
	//for name, values := range resp.Header {
	//	for _, value := range values {
	//		p.logger.Printf("[RESP:%s] Header: %s: %s", requestID, name, value)
	//	}
	//}

	w.WriteHeader(resp.StatusCode)
	bytesWritten, err := io.Copy(w, resp.Body)
	if err != nil {
		//p.logger.Printf("[ERR:%s] Failed to write response: %v", requestID, err)
		return
	}
	fmt.Println(bytesWritten)

}
