package store

import "time"

// Record 内部任务记录。
type Record struct {
	ID          string
	Name        string
	Kind        string
	Payload     []byte
	Status      string
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

// Store 存储抽象。
type Store interface {
	Insert(Record) (string, error)
	Get(id string) (Record, error)
	List(limit int) ([]Record, error)
	ListDue(now time.Time, limit int) ([]Record, error)
	Update(Record) error
	CountByStatus(status string) (int, error)
	Flush() error
	Close() error
}
