package payload_test

import (
	"testing"

	"github.com/LYH2263/go-jobcron/internal/payload"
)

func TestCloneBytes(t *testing.T) {
	in := []byte("abc")
	out := payload.CloneBytes(in)
	in[0] = 'Z'
	if out[0] == 'Z' {
		t.Fatal("alias")
	}
}
