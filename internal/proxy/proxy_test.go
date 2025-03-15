package proxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createMockTargetServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/test":
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Target server test response"))
		case "/headers":
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(r.Header.Get("X-Test-Header")))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func TestProxyServer_ServeHTTP(t *testing.T) {
	targetServer := createMockTargetServer()
	defer targetServer.Close()

	proxyPort := "8081"
	proxy := NewProxyServer(Config{})

	server := &http.Server{
		Addr:    ":" + proxyPort,
		Handler: proxy,
	}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			t.Logf("Stop Proxy server: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	defer func() {
		server.Close()
	}()

	proxyURL, err := url.Parse("http://localhost:" + proxyPort)
	require.NoError(t, err)

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	t.Run("Basic proxy request", func(t *testing.T) {
		resp, err := client.Get(targetServer.URL + "/test")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, "Target server test response", string(body))
	})

	t.Run("Header forwarding", func(t *testing.T) {
		req, err := http.NewRequest("GET", targetServer.URL+"/headers", nil)
		require.NoError(t, err)

		// 커스텀 헤더 추가
		req.Header.Set("X-Test-Header", "test-value")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Header.Get("X-Test-Header")

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, "test-value", string(body))
	})

	t.Run("fail test", func(t *testing.T) {
		req, err := http.NewRequest("GET", "http://localhost:8080"+"/headers", nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.Nil(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		resp.Header.Get("X-Test-Header")

		body, err := io.ReadAll(resp.Body)
		require.Nil(t, err)
		assert.Contains(t, string(body), "connection refused")
	})

	t.Run("404 error handling", func(t *testing.T) {
		resp, err := client.Get(targetServer.URL + "/nonexistent")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
