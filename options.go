package jobcron

import (
	"net/http"
	"time"

	"github.com/LYH2263/go-jobcron/internal/clock"
	"github.com/LYH2263/go-jobcron/internal/runner"
)

// Option 配置 Scheduler。
type Option func(*Scheduler)

// WithClock 注入时钟（测试可用 Fake）。
func WithClock(c clock.Clock) Option {
	return func(s *Scheduler) {
		if c != nil {
			s.clk = c
		}
	}
}

// WithHandler 注册 handler 型任务执行函数。
func WithHandler(h runner.Handler) Option {
	return func(s *Scheduler) {
		s.handler = h
	}
}

// WithPersistPath 可选 JSON 快照路径。
func WithPersistPath(path string) Option {
	return func(s *Scheduler) {
		s.persistPath = path
	}
}

// WithMaxAttempts 单任务最大尝试次数。
func WithMaxAttempts(n int) Option {
	return func(s *Scheduler) {
		if n > 0 {
			s.maxAttempts = n
		}
	}
}

// WithBackoff 指数退避参数。
func WithBackoff(base, cap time.Duration, factor float64) Option {
	return func(s *Scheduler) {
		if base > 0 {
			s.backoffBase = base
		}
		if cap > 0 {
			s.backoffCap = cap
		}
		if factor >= 1 {
			s.backoffFactor = factor
		}
	}
}

// WithHTTPTimeout HTTP 回调超时。
func WithHTTPTimeout(d time.Duration) Option {
	return func(s *Scheduler) {
		if d > 0 {
			s.httpTimeout = d
		}
	}
}

// WithTransport 注入 RoundTripper。
func WithTransport(rt http.RoundTripper) Option {
	return func(s *Scheduler) {
		s.transport = rt
	}
}

// WithCapacity 内存任务容量上限。
func WithCapacity(n int) Option {
	return func(s *Scheduler) {
		if n > 0 {
			s.capacity = n
		}
	}
}

// WithDefaultURL 默认 HTTP 回调 URL。
func WithDefaultURL(u string) Option {
	return func(s *Scheduler) {
		s.defaultURL = u
	}
}
