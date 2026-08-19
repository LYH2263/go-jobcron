package runner

import (
	"context"
	"time"

	"github.com/LYH2263/go-jobcron/internal/clock"
	"github.com/LYH2263/go-jobcron/internal/httpx"
	"github.com/LYH2263/go-jobcron/internal/queue"
	"github.com/LYH2263/go-jobcron/internal/store"
)

// Config Engine 配置。
type Config struct {
	Store         store.Store
	Queue         *queue.Queue
	Clock         clock.Clock
	Handler       Handler
	Client        *httpx.Client
	MaxAttempts   int
	BackoffBase   time.Duration
	BackoffCap    time.Duration
	BackoffFactor float64
	DefaultURL    string
	UserAgent     string
}

// Engine 任务执行引擎。
type Engine struct {
	st            store.Store
	q             *queue.Queue
	clk           clock.Clock
	handler       Handler
	client        *httpx.Client
	maxAttempts   int
	backoffBase   time.Duration
	backoffCap    time.Duration
	backoffFactor float64
	defaultURL    string
	userAgent     string
}

// NewEngine 构造 Engine。
func NewEngine(cfg Config) *Engine {
	return &Engine{
		st:            cfg.Store,
		q:             cfg.Queue,
		clk:           cfg.Clock,
		handler:       cfg.Handler,
		client:        cfg.Client,
		maxAttempts:   cfg.MaxAttempts,
		backoffBase:   cfg.BackoffBase,
		backoffCap:    cfg.BackoffCap,
		backoffFactor: cfg.BackoffFactor,
		defaultURL:    cfg.DefaultURL,
		userAgent:     cfg.UserAgent,
	}
}

// RunOne 执行一条到期任务。
func (e *Engine) RunOne(ctx context.Context, now time.Time) (Result, bool, error) {
	if e.st == nil || e.q == nil {
		return Result{}, false, ErrRun
	}
	rec, ok, err := ClaimNext(e.st, e.q, now)
	if err != nil {
		return Result{}, false, err
	}
	if !ok {
		return Result{}, false, nil
	}
	res, err := Execute(ctx, e, rec)
	return res, true, err
}
