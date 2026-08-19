package jobcron_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/LYH2263/go-jobcron"
)

func TestBug10_CloseFlushesStoreBeforeRelease(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "jobs.json")
	s := jobcron.New(
		jobcron.WithDefaultURL("http://127.0.0.1:9/hook"),
		jobcron.WithPersistPath(path),
	)
	id, err := s.Enqueue("once", []byte(`{"n":1}`), time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(raw) == 0 {
		t.Fatal("persist file empty after Close")
	}
	var snap struct {
		Items []struct {
			ID string `json:"ID"`
		} `json:"items"`
	}
	if err := json.Unmarshal(raw, &snap); err != nil {
		t.Fatal(err)
	}
	if len(snap.Items) != 1 || snap.Items[0].ID != id {
		t.Fatalf("Close released store before Flush; loaded=%+v", snap.Items)
	}
}
