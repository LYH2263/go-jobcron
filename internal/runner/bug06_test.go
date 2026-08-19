package runner_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/LYH2263/go-jobcron/internal/clock"
	"github.com/LYH2263/go-jobcron/internal/queue"
	"github.com/LYH2263/go-jobcron/internal/runner"
	"github.com/LYH2263/go-jobcron/internal/store"
)

func TestBug06_MarkDonePersistFailureRollsBack(t *testing.T) {
	path := filepath.Join(t.TempDir(), "snap.json")
	st := store.NewMemory(16, clock.Real{}, path)
	q := queue.New()
	eng := runner.NewEngine(runner.Config{Store: st, Queue: q, Clock: clock.Real{}, MaxAttempts: 3})
	id, err := st.Insert(store.Record{
		Name: "t", Kind: "once", Payload: []byte("{}"), Status: "running", MaxAttempts: 3,
	})
	if err != nil {
		t.Fatal(err)
	}
	st.SetFlushError(errors.New("disk full"))
	rec, err := st.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	err = runner.MarkDone(eng, rec)
	if err == nil {
		t.Fatal("expected persist error")
	}
	again, err := st.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if again.Status == "done" {
		t.Fatalf("memory mutated to done despite persist failure: %+v", again)
	}
	if again.Status != "running" {
		t.Fatalf("want running rollback, got %s", again.Status)
	}
}
