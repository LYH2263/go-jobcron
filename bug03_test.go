package jobcron_test

import (
	"errors"
	"testing"
	"time"

	"github.com/LYH2263/go-jobcron"
)

func TestBug03_TickAfterCloseNoPanic(t *testing.T) {
	s := jobcron.New(jobcron.WithDefaultURL("http://127.0.0.1:9/hook"))
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Tick panicked after Close: %v", r)
		}
	}()
	_, err := s.TickAt(time.Now())
	if !errors.Is(err, jobcron.ErrClosed) {
		t.Fatalf("want ErrClosed, got %v", err)
	}
}
