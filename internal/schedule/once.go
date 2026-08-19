package schedule

import "time"

// NextOnce 一次性任务的下一次运行即 runAt。
func NextOnce(runAt time.Time) time.Time {
	return runAt
}

// IsDue 是否到期（含等于）。
func IsDue(next, now time.Time) bool {
	return !next.IsZero() && !next.After(now)
}
