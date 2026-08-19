package jobcron

import (
	"context"
	"time"
)

// Shutdown 带超时的优雅关闭：先 Flush 再 Close。
func (s *Scheduler) Shutdown(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	done := make(chan error, 1)
	go func() {
		done <- s.Close()
	}()
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return mapCtxErr(ctx.Err())
	case <-time.After(30 * time.Second):
		return ErrTimeout
	}
}
