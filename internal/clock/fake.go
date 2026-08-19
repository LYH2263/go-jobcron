package clock

import (
	"sync"
	"time"
)

// Fake 可手动推进的测试时钟。
type Fake struct {
	mu sync.Mutex
	t  time.Time
}

// NewFake 构造 Fake，起点为 t（零值则用 Unix(0)）。
func NewFake(t time.Time) *Fake {
	if t.IsZero() {
		t = time.Unix(0, 0)
	}
	return &Fake{t: t}
}

func (f *Fake) Now() time.Time {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.t
}

// Set 设置当前时刻。
func (f *Fake) Set(t time.Time) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.t = t
}

// Advance 向前推进 d。
func (f *Fake) Advance(d time.Duration) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.t = f.t.Add(d)
}
