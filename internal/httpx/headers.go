package httpx

import (
	"errors"
	"net/http"
)

var ErrNilClient = errors.New("httpx: nil client")

// CloneHeaders 深拷贝 HTTP 头。
func CloneHeaders(h http.Header) http.Header {
	if h == nil {
		return nil
	}
	out := make(http.Header, len(h))
	for k, vv := range h {
		cp := make([]string, len(vv))
		copy(cp, vv)
		out[k] = cp
	}
	return out
}

// SetJSON 设置 JSON 内容类型。
func SetJSON(h http.Header) {
	if h == nil {
		return
	}
	h.Set("Content-Type", "application/json")
}
