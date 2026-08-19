package jobcron

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/LYH2263/go-jobcron/internal/runner"
)

func TestBug05_RunErrorWrapsSentinel(t *testing.T) {
	s := New()
	defer s.Close()
	if _, err := s.EnqueueHTTP("cb", []byte(`{}`), time.Now(), "http://127.0.0.1:1/dead"); err != nil {
		t.Fatal(err)
	}
	_, err := s.RunDue(context.Background())
	if err == nil {
		t.Fatal("expected RunDue error")
	}
	if !errors.Is(err, runner.ErrHTTP) && !errors.Is(err, runner.ErrRun) {
		t.Fatalf("RunDue must errors.Is ErrHTTP/ErrRun, got %v", err)
	}
}
