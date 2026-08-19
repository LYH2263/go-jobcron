package schedule

import (
	"time"

	"github.com/LYH2263/go-jobcron/internal/cronexpr"
)

// NextCron 计算 cron 下次触发。
func NextCron(expr string, after time.Time) (time.Time, error) {
	return cronexpr.Next(expr, after)
}

// RescheduleCron 任务执行后计算下一次（cron 专用）。
func RescheduleCron(expr string, finished time.Time) (time.Time, error) {
	return cronexpr.Next(expr, finished)
}
