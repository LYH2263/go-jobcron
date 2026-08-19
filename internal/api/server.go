package api

import (
	"net/http"

	"github.com/LYH2263/go-jobcron"
)

// Options HTTP 服务配置。
type Options struct {
	WebDir    string
	AllowCORS bool
}

// Server 管理 API + 静态页。
type Server struct {
	box *jobcron.Scheduler
	opt Options
}

// New 构造 Server（sched 为调度器指针）。
func New(sched *jobcron.Scheduler, opt Options) http.Handler {
	s := &Server{box: sched, opt: opt}
	return s.routes()
}
