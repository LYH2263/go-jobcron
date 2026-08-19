package jobcron_test

import (
	"testing"
	"time"

	"github.com/LYH2263/go-jobcron"
)

func TestBug02_ListJobsPayloadSliceAlias(t *testing.T) {
	s := jobcron.New(jobcron.WithDefaultURL("http://127.0.0.1:9/hook"))
	defer s.Close()
	id, err := s.Enqueue("once", []byte(`{"id":1}`), time.Now())
	if err != nil {
		t.Fatal(err)
	}
	list, err := s.ListJobs(10)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("len=%d", len(list))
	}
	list[0].Payload[2] = 'Z'
	job, err := s.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if string(job.Payload) != `{"id":1}` {
		t.Fatalf("payload alias into store: %q", job.Payload)
	}
}
