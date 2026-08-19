package api

import "net/http"

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/api/stats", s.handleStats)
	mux.HandleFunc("/api/jobs", s.handleJobs)
	mux.HandleFunc("/api/jobs/", s.handleJobByID)
	mux.HandleFunc("/api/run", s.handleRun)
	if s.opt.WebDir != "" {
		mux.Handle("/", staticDir(s.opt.WebDir))
	}
	return mux
}
