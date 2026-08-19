package jobcron_test

import (
	"testing"
	"time"

	"github.com/LYH2263/go-jobcron"
)

func TestBug01_EnqueuePayloadSliceAlias(t *testing.T) {
	s := jobcron.New(jobcron.WithDefaultURL("http://127.0.0.1:9/hook"))
	defer s.Close()
	payload := []byte(`{"id":1}`)
	id, err := s.Enqueue("once", payload, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	payload[2] = 'Z'
	job, err := s.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if string(job.Payload) != `{"id":1}` {
		t.Fatalf("store payload polluted by caller alias: %q", job.Payload)
	}
}
