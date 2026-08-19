package runner_test

import (
	"context"
	"io"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/LYH2263/go-jobcron/internal/httpx"
	"github.com/LYH2263/go-jobcron/internal/runner"
	"github.com/LYH2263/go-jobcron/internal/store"
)

type closeCounter struct {
	n atomic.Int32
}

func (c *closeCounter) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       &countCloser{c: c, r: http.NoBody},
		Header:     make(http.Header),
		Request:    req,
	}, nil
}

type countCloser struct {
	c *closeCounter
	r io.ReadCloser
}

func (c *countCloser) Read(p []byte) (int, error) { return c.r.Read(p) }
func (c *countCloser) Close() error {
	c.c.n.Add(1)
	return c.r.Close()
}

func TestBug09_CallbackBodyClosed(t *testing.T) {
	ctr := &closeCounter{}
	client := httpx.New(0, ctr, "ua")
	_, err := runner.HTTPCallback(context.Background(), client, store.Record{
		ID: "1", Name: "n", URL: "http://example.invalid/hook", Payload: []byte("{}"),
	}, "ua")
	if err != nil {
		t.Fatal(err)
	}
	if ctr.n.Load() < 1 {
		t.Fatalf("response Body Close count=%d, want >=1", ctr.n.Load())
	}
}
