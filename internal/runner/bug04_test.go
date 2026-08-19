package runner_test

import (
	"context"
	"errors"
	"testing"

	"github.com/LYH2263/go-jobcron/internal/runner"
	"github.com/LYH2263/go-jobcron/internal/store"
)

func TestBug04_NilHandlerNoPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("nil handler panicked: %v", r)
		}
	}()
	err := runner.InvokeHandler(nil, context.Background(), store.Record{ID: "1", Name: "n"})
	if !errors.Is(err, runner.ErrNilHandler) {
		t.Fatalf("want ErrNilHandler, got %v", err)
	}
	_, err = runner.HTTPCallback(context.Background(), nil, store.Record{URL: "http://127.0.0.1:9/"}, "")
	if !errors.Is(err, runner.ErrNilClient) {
		t.Fatalf("want ErrNilClient, got %v", err)
	}
}
