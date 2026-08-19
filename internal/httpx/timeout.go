package httpx

import "time"

// DefaultTimeout 默认 HTTP 超时。
const DefaultTimeout = 10 * time.Second

// WithTimeout 返回带超时的 Client（便捷构造）。
func WithTimeout(d time.Duration, ua string) *Client {
	return New(d, nil, ua)
}
