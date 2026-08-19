package runner

import (
	"context"
	"math"
	"time"
)

// Delay 计算指数退避延迟。
func Delay(base, cap time.Duration, factor float64, attempts int) time.Duration {
	if attempts < 1 {
		attempts = 1
	}
	if base <= 0 {
		base = time.Millisecond
	}
	if cap <= 0 {
		cap = time.Second
	}
	if factor < 1 {
		factor = 2
	}
	d := float64(base) * math.Pow(factor, float64(attempts-1))
	if d > float64(cap) {
		d = float64(cap)
	}
	return time.Duration(d)
}

// Wait 等待退避，尊重 ctx 取消。
func Wait(ctx context.Context, base, cap time.Duration, factor float64, attempts int) error {
	if ctx == nil {
		ctx = context.Background()
	}
	d := Delay(base, cap, factor, attempts)
	if d <= 0 {
		return ctx.Err()
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
