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
)

func TestBug07_RunDueContextHonorsCancel(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		defer r.Body.Close()
		_, _ = io.Copy(io.Discard, r.Body)
		w.WriteHeader(200)
	}))
	defer srv.Close()

	s := jobcron.New(
		jobcron.WithDefaultURL(srv.URL),
		jobcron.WithHTTPTimeout(2*time.Second),
	)
	defer s.Close()
	for i := 0; i < 5; i++ {
		if _, err := s.EnqueueHTTP("t", []byte(`{}`), time.Now(), srv.URL); err != nil {
			t.Fatal(err)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	start := time.Now()
	res, err := s.RunDueContext(ctx, 5)
	if err == nil {
		t.Fatal("expected cancel error")
	}
	if !(errors.Is(err, jobcron.ErrCanceled) || errors.Is(err, context.Canceled)) {
		t.Fatalf("want canceled, got %v", err)
	}
	if hits.Load() != 0 {
		t.Fatalf("RunDueContext ignored cancel and delivered hits=%d results=%d", hits.Load(), len(res))
	}
	if time.Since(start) > time.Second {
		t.Fatalf("RunDueContext ignored cancel, took %s", time.Since(start))
	}
}
