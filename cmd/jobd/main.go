package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LYH2263/go-jobcron"
	"github.com/LYH2263/go-jobcron/internal/api"
)

func main() {
	addr := flag.String("addr", ":8093", "HTTP 监听地址")
	web := flag.String("web", "web", "静态管理页目录")
	defaultURL := flag.String("url", "http://127.0.0.1:9999/hook", "默认 HTTP 回调 URL")
	persist := flag.String("persist", "", "任务快照 JSON 路径（可选）")
	timeout := flag.Duration("http-timeout", 10*time.Second, "HTTP 回调超时")
	maxTry := flag.Int("max-attempts", 5, "最大尝试次数")
	flag.Parse()

	opts := []jobcron.Option{
		jobcron.WithDefaultURL(*defaultURL),
		jobcron.WithHTTPTimeout(*timeout),
		jobcron.WithMaxAttempts(*maxTry),
	}
	if *persist != "" {
		opts = append(opts, jobcron.WithPersistPath(*persist))
	}

	sched := jobcron.New(opts...)
	defer sched.Close()

	srv := api.New(sched, api.Options{WebDir: *web, AllowCORS: true})
	httpSrv := &http.Server{
		Addr:              *addr,
		Handler:           srv,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("jobd 监听 %s", *addr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = httpSrv.Shutdown(ctx)
	_ = sched.Shutdown(ctx)
}
