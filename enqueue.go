package jobcron

import (
	"time"

	"github.com/LYH2263/go-jobcron/internal/payload"
	"github.com/LYH2263/go-jobcron/internal/schedule"
	"github.com/LYH2263/go-jobcron/internal/store"
)

// Enqueue 入队一次性延迟任务。
func (s *Scheduler) Enqueue(name string, pl []byte, runAt time.Time) (string, error) {
	return s.enqueueKind(name, payload.CloneBytes(pl), KindOnce, "", "", runAt)
}

// EnqueueCron 入队 cron 周期任务。
func (s *Scheduler) EnqueueCron(name, expr string, pl []byte) (string, error) {
	if expr == "" {
		return "", ErrEmptyCron
	}
	next, err := schedule.NextCron(expr, s.now())
	if err != nil {
		return "", err
	}
	return s.enqueueKind(name, payload.CloneBytes(pl), KindCron, expr, "", next)
}

// EnqueueHTTP 入队 HTTP 回调任务。
func (s *Scheduler) EnqueueHTTP(name string, pl []byte, runAt time.Time, url string) (string, error) {
	if url == "" {
		url = s.defaultURLLocked()
	}
	if url == "" {
		return "", ErrEmptyURL
	}
	return s.enqueueKind(name, payload.CloneBytes(pl), KindHTTPCallback, "", url, runAt)
}

// EnqueueHandler 入队 handler 型任务（需 WithHandler）。
func (s *Scheduler) EnqueueHandler(name string, pl []byte, runAt time.Time) (string, error) {
	return s.enqueueKind(name, payload.CloneBytes(pl), KindHandler, "", "", runAt)
}

func (s *Scheduler) enqueueKind(name string, pl []byte, kind Kind, cronExpr, url string, runAt time.Time) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.checkOpenLocked(); err != nil {
		return "", err
	}
	if name == "" {
		return "", ErrEmptyName
	}
	if len(pl) == 0 {
		return "", ErrEmptyPayload
	}
	if runAt.IsZero() {
		runAt = s.clk.Now()
	}
	rec := store.Record{
		Name:        name,
		Kind:        string(kind),
		Payload:     pl,
		Status:      string(StatusPending),
		CronExpr:    cronExpr,
		URL:         url,
		RunAt:       runAt,
		NextRun:     runAt,
		MaxAttempts: s.maxAttempts,
	}
	id, err := s.st.Insert(rec)
	if err != nil {
		return "", err
	}
	s.q.Push(id, runAt)
	s.enqueued++
	return id, nil
}

func (s *Scheduler) now() time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.clk == nil {
		return time.Now()
	}
	return s.clk.Now()
}

func (s *Scheduler) defaultURLLocked() string {
	return s.defaultURL
}
