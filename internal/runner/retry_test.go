package runner_test

import (
	"context"
	"testing"
	"time"

	"github.com/LYH2263/go-jobcron/internal/runner"
)

func TestWaitHonorsCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	start := time.Now()
	err := runner.Wait(ctx, time.Millisecond, time.Second, 2, 1)
	if err == nil {
		t.Fatal("expected err")
	}
	if time.Since(start) > time.Second {
		t.Fatal("slow")
	}
}
