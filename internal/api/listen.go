package api

import (
	"context"
	"net/http"
	"time"
)

// Listen 启动 HTTP 服务。
func Listen(addr string, h http.Handler) (*http.Server, error) {
	srv := &http.Server{
		Addr:              addr,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		_ = srv.ListenAndServe()
	}()
	return srv, nil
}

// Shutdown 优雅关闭。
func Shutdown(ctx context.Context, srv *http.Server) error {
	if srv == nil {
		return nil
	}
	return srv.Shutdown(ctx)
}
