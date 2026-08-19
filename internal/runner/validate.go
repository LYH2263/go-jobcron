package runner

import "github.com/LYH2263/go-jobcron/internal/store"

// ValidateRunnable 检查任务是否可运行。
func ValidateRunnable(rec store.Record) error {
	if rec.Status != "pending" {
		return ErrRun
	}
	return nil
}
