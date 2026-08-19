package jobcron

import "errors"

var (
	ErrClosed       = errors.New("jobcron: closed")
	ErrNilScheduler = errors.New("jobcron: nil scheduler")
	ErrEmptyName    = errors.New("jobcron: empty name")
	ErrEmptyPayload = errors.New("jobcron: empty payload")
	ErrEmptyURL     = errors.New("jobcron: empty url")
	ErrEmptyCron    = errors.New("jobcron: empty cron expr")
	ErrNotFound     = errors.New("jobcron: job not found")
	ErrInvalidKind  = errors.New("jobcron: invalid kind")
	ErrInvalidStatus = errors.New("jobcron: invalid status")
	ErrPersist      = errors.New("jobcron: persist failed")
	ErrRun          = errors.New("jobcron: run failed")
	ErrHTTP         = errors.New("jobcron: http error")
	ErrTimeout      = errors.New("jobcron: timeout")
	ErrCanceled     = errors.New("jobcron: canceled")
	ErrNilHandler   = errors.New("jobcron: nil handler")
	ErrNilClient    = errors.New("jobcron: nil http client")
	ErrCapacity     = errors.New("jobcron: capacity exceeded")
	ErrBadCron      = errors.New("jobcron: bad cron expr")
)
