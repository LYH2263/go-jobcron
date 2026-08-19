package cronexpr_test

import (
	"testing"
	"time"

	"github.com/LYH2263/go-jobcron/internal/cronexpr"
)

func TestNextEveryMinute(t *testing.T) {
	after := time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)
	next, err := cronexpr.Next("* * * * *", after)
	if err != nil {
		t.Fatal(err)
	}
	if !next.After(after) {
		t.Fatalf("next=%v after=%v", next, after)
	}
}
