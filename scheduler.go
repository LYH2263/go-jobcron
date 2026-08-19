package jobcron

import (
	"net/http"
	"sync"
	"time"

	"github.com/LYH2263/go-jobcron/internal/clock"
	"github.com/LYH2263/go-jobcron/internal/httpx"
	"github.com/LYH2263/go-jobcron/internal/queue"
	"github.com/LYH2263/go-jobcron/internal/runner"
	"github.com/LYH2263/go-jobcron/internal/store"
)

const (
	defaultHTTPTimeout = 10 * time.Second
	defaultMaxAttempts = 5
	defaultBackoffBase = 50 * time.Millisecond
	defaultBackoffCap  = 2 * time.Second
	defaultFactor      = 2.0
	defaultCapacity    = 4096
	DefaultUA          = "go-jobcron/1.0"
)

// Scheduler 任务调度门面。零值不可用，须经 New 构造。
type Scheduler struct {
	mu sync.Mutex

	closed bool
	clk    clock.Clock
	st     store.Store
	q      *queue.Queue
	eng    *runner.Engine
	handler runner.Handler

	httpTimeout   time.Duration
	transport     http.RoundTripper
	defaultURL    string
	maxAttempts   int
	backoffBase   time.Duration
	backoffCap    time.Duration
	backoffFactor float64
	capacity      int
	persistPath   string
	userAgent     string

	enqueued  uint64
	executed  uint64
	succeeded uint64
	failed    uint64
	dead      uint64
	lastName  string
	lastRun   time.Time
}

// New 构造 Scheduler。
func New(opts ...Option) *Scheduler {
	s := &Scheduler{
		clk:           clock.Real{},
		httpTimeout:   defaultHTTPTimeout,
		maxAttempts:   defaultMaxAttempts,
		backoffBase:   defaultBackoffBase,
		backoffCap:    defaultBackoffCap,
		backoffFactor: defaultFactor,
		capacity:      defaultCapacity,
		userAgent:     DefaultUA,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(s)
		}
	}
	if s.clk == nil {
		s.clk = clock.Real{}
	}
	if s.maxAttempts < 1 {
		s.maxAttempts = 1
	}
	if s.capacity < 8 {
		s.capacity = 8
	}
	s.st = store.NewMemory(s.capacity, s.clk, s.persistPath)
	s.q = queue.New()
	client := httpx.New(s.httpTimeout, s.transport, s.userAgent)
	s.eng = runner.NewEngine(runner.Config{
		Store:         s.st,
		Queue:         s.q,
		Clock:         s.clk,
		Handler:       s.handler,
		Client:        client,
		MaxAttempts:   s.maxAttempts,
		BackoffBase:   s.backoffBase,
		BackoffCap:    s.backoffCap,
		BackoffFactor: s.backoffFactor,
		DefaultURL:    s.defaultURL,
		UserAgent:     s.userAgent,
	})
	return s
}

func (s *Scheduler) checkOpenLocked() error {
	if s == nil {
		return ErrNilScheduler
	}
	if s.closed {
		return ErrClosed
	}
	if s.st == nil || s.q == nil || s.eng == nil {
		return ErrClosed
	}
	return nil
}

// Store 返回底层存储（Close 后可能为 nil）。
func (s *Scheduler) Store() store.Store {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.st
}

// Queue 返回到期队列（Close 后可能为 nil）。
func (s *Scheduler) Queue() *queue.Queue {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.q
}

// Closed 是否已关闭。
func (s *Scheduler) Closed() bool {
	if s == nil {
		return true
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.closed
}
