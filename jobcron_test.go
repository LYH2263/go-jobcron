package jobcron_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/LYH2263/go-jobcron"
	"github.com/LYH2263/go-jobcron/internal/clock"
	"github.com/LYH2263/go-jobcron/internal/runner"
	"github.com/LYH2263/go-jobcron/internal/store"
)

func TestEnqueueRunHTTP(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		defer r.Body.Close()
		_, _ = io.Copy(io.Discard, r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	s := jobcron.New(
		jobcron.WithDefaultURL(srv.URL),
		jobcron.WithMaxAttempts(3),
		jobcron.WithHTTPTimeout(2*time.Second),
	)
	defer s.Close()

	payload := []byte(`{"n":1}`)
	id, err := s.EnqueueHTTP("test", payload, time.Now(), srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	payload[0] = 'X'
	list, err := s.ListJobs(10)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || string(list[0].Payload) != `{"n":1}` {
		t.Fatalf("list=%+v", list)
	}
	list[0].Payload[0] = 'Y'
	again, _ := s.Get(id)
	if again.Payload[0] == 'Y' {
		t.Fatal("payload alias")
	}

	res, err := s.RunDue(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !res.OK || hits.Load() != 1 {
		t.Fatalf("res=%+v hits=%d", res, hits.Load())
	}
}

func TestRunAfterClose(t *testing.T) {
	s := jobcron.New(jobcron.WithDefaultURL("http://127.0.0.1:1/x"))
	_ = s.Close()
	_, err := s.RunDue(context.Background())
	if !errors.Is(err, jobcron.ErrClosed) {
		t.Fatalf("got %v", err)
	}
	_, err = s.Tick()
	if !errors.Is(err, jobcron.ErrClosed) {
		t.Fatalf("tick got %v", err)
	}
}

func TestRunDueContextCancel(t *testing.T) {
	block := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-block
		w.WriteHeader(200)
	}))
	defer srv.Close()
	defer close(block)

	s := jobcron.New(jobcron.WithDefaultURL(srv.URL))
	defer s.Close()
	_, _ = s.EnqueueHTTP("t", []byte("{}"), time.Now(), srv.URL)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := s.RunDueContext(ctx, 10)
	if err == nil {
		t.Fatal("expected cancel")
	}
}

func TestHandlerJob(t *testing.T) {
	var called atomic.Bool
	s := jobcron.New(jobcron.WithHandler(func(ctx context.Context, rec store.Record) error {
		called.Store(true)
		return nil
	}))
	defer s.Close()
	_, err := s.EnqueueHandler("h", []byte("{}"), time.Now())
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.RunDue(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !called.Load() {
		t.Fatal("handler not called")
	}
}

func TestNilHandlerError(t *testing.T) {
	s := jobcron.New()
	defer s.Close()
	_, _ = s.EnqueueHandler("h", []byte("{}"), time.Now())
	_, err := s.RunDue(context.Background())
	if err == nil || !errors.Is(err, runner.ErrRun) {
		t.Fatalf("got %v", err)
	}
}

func TestCronEnqueue(t *testing.T) {
	fc := clock.NewFake(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
	s := jobcron.New(jobcron.WithClock(fc))
	defer s.Close()
	id, err := s.EnqueueCron("c", "* * * * *", []byte("{}"))
	if err != nil || id == "" {
		t.Fatalf("%s %v", id, err)
	}
}
