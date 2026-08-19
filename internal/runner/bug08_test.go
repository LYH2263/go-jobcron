package runner_test

import (
	"context"
	"testing"
	"time"

	"github.com/LYH2263/go-jobcron/internal/runner"
)

func TestBug08_RetryWaitHonorsContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	start := time.Now()
	err := runner.Wait(ctx, 3*time.Second, 3*time.Second, 2, 1)
	if err == nil {
		t.Fatal("expected ctx error")
	}
	if time.Since(start) > time.Second {
		t.Fatalf("Wait used Sleep ignoring ctx, took %s", time.Since(start))
	}
}
