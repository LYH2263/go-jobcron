package schedule

import (
	"time"

	"github.com/LYH2263/go-jobcron/internal/cronexpr"
)

// ComputeNext 按 kind 计算下次运行。
func ComputeNext(kind, cronExpr string, runAt, after time.Time) (time.Time, error) {
	switch kind {
	case "once", "http_callback", "handler":
		if IsDue(runAt, after) {
			return runAt, nil
		}
		return runAt, nil
	case "cron":
		if cronExpr == "" {
			return time.Time{}, cronexpr.Validate("")
		}
		return cronexpr.Next(cronExpr, after)
	default:
		return time.Time{}, nil
	}
}
