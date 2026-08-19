package httpx

// OK 是否 2xx。
func OK(code int) bool {
	return code >= 200 && code < 300
}

// Retryable 是否可重试状态码。
func Retryable(code int) bool {
	return code == 429 || code >= 500
}
