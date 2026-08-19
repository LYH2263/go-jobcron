package runner

import "time"

// Result 单次执行结果。
type Result struct {
	JobID      string
	Name       string
	OK         bool
	Attempts   int
	StatusCode int
	Duration   time.Duration
	Err        string
	FinishedAt time.Time
}
