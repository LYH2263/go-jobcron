package jobcron

import "time"

// Kind 任务类型。
type Kind string

const (
	KindOnce         Kind = "once"
	KindCron         Kind = "cron"
	KindHTTPCallback Kind = "http_callback"
	KindHandler      Kind = "handler"
)

// Status 任务生命周期。
type Status string

const (
	StatusPending Status = "pending"
	StatusRunning Status = "running"
	StatusDone    Status = "done"
	StatusFailed  Status = "failed"
	StatusPaused  Status = "paused"
	StatusDead    Status = "dead"
)

// Job 对外可见任务视图（载荷已深拷贝）。
type Job struct {
	ID          string
	Name        string
	Kind        Kind
	Payload     []byte
	Status      Status
	CronExpr    string
	URL         string
	RunAt       time.Time
	NextRun     time.Time
	Attempts    int
	MaxAttempts int
	LastError   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// RunResult 单次执行结果。
type RunResult struct {
	JobID      string
	Name       string
	OK         bool
	Attempts   int
	StatusCode int
	Duration   time.Duration
	Err        string
	FinishedAt time.Time
}

// Stats 运行统计。
type Stats struct {
	Enqueued  uint64
	Executed  uint64
	Succeeded uint64
	Failed    uint64
	Dead      uint64
	Pending   int
	Paused    int
	Closed    bool
	LastName  string
	LastRun   time.Time
}
