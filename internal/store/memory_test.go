package store_test

import (
	"testing"

	"github.com/LYH2263/go-jobcron/internal/clock"
	"github.com/LYH2263/go-jobcron/internal/store"
)

func TestMemoryCloneGet(t *testing.T) {
	m := store.NewMemory(8, clock.Real{}, "")
	id, err := m.Insert(store.Record{
		Name: "a", Kind: "once", Payload: []byte("p"), Status: "pending", MaxAttempts: 2,
	})
	if err != nil {
		t.Fatal(err)
	}
	rec, _ := m.Get(id)
	rec.Payload[0] = 'Z'
	again, _ := m.Get(id)
	if again.Payload[0] == 'Z' {
		t.Fatal("alias")
	}
}

func TestFlushBeforeClose(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/snap.json"
	m := store.NewMemory(8, clock.Real{}, path)
	_, _ = m.Insert(store.Record{
		Name: "a", Kind: "once", Payload: []byte("p"), Status: "pending", MaxAttempts: 2,
	})
	if err := m.Flush(); err != nil {
		t.Fatal(err)
	}
	if err := m.Close(); err != nil {
		t.Fatal(err)
	}
}
