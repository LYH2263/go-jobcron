package httpx

import (
	"net/http"
	"time"
)

// Client 轻量 HTTP 客户端封装。
type Client struct {
	hc *http.Client
	ua string
}

// New 构造 Client。
func New(timeout time.Duration, rt http.RoundTripper, ua string) *Client {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	tr := rt
	if tr == nil {
		tr = http.DefaultTransport
	}
	return &Client{
		hc: &http.Client{Timeout: timeout, Transport: tr},
		ua: ua,
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if c == nil || c.hc == nil {
		return nil, ErrNilClient
	}
	if c.ua != "" && req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", c.ua)
	}
	return c.hc.Do(req)
}

// CloseIdle 关闭空闲连接。
func (c *Client) CloseIdle() {
	if c == nil || c.hc == nil {
		return
	}
	if tr, ok := c.hc.Transport.(interface{ CloseIdleConnections() }); ok {
		tr.CloseIdleConnections()
	}
}
