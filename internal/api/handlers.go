package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "closed": s.box.Closed()})
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.box.Stats())
}

func (s *Server) handleJobs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		jobs, err := s.box.ListJobs(limit)
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, http.StatusOK, jobs)
	case http.MethodPost:
		s.handleCreateJob(w, r)
	default:
		http.Error(w, "method", http.StatusMethodNotAllowed)
	}
}

type createBody struct {
	Name     string          `json:"name"`
	Kind     string          `json:"kind"`
	Payload  json.RawMessage `json:"payload"`
	RunAt    string          `json:"run_at"`
	CronExpr string          `json:"cron_expr"`
	URL      string          `json:"url"`
}

func (s *Server) handleCreateJob(w http.ResponseWriter, r *http.Request) {
	var body createBody
	if err := readJSON(r, &body); err != nil {
		writeErr(w, err)
		return
	}
	pl := []byte(body.Payload)
	var (
		id  string
		err error
	)
	runAt := time.Now()
	if body.RunAt != "" {
		runAt, err = time.Parse(time.RFC3339, body.RunAt)
		if err != nil {
			writeErr(w, err)
			return
		}
	}
	switch body.Kind {
	case "cron":
		id, err = s.box.EnqueueCron(body.Name, body.CronExpr, pl)
	case "http_callback":
		id, err = s.box.EnqueueHTTP(body.Name, pl, runAt, body.URL)
	case "handler":
		id, err = s.box.EnqueueHandler(body.Name, pl, runAt)
	default:
		id, err = s.box.Enqueue(body.Name, pl, runAt)
	}
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (s *Server) handleJobByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/jobs/")
	if id == "" {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		job, err := s.box.Get(id)
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, http.StatusOK, job)
	case http.MethodPost:
		action := r.URL.Query().Get("action")
		var err error
		switch action {
		case "pause":
			err = s.box.Pause(id)
		case "resume":
			err = s.box.Resume(id)
		default:
			http.Error(w, "unknown action", http.StatusBadRequest)
			return
		}
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
	default:
		http.Error(w, "method", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method", http.StatusMethodNotAllowed)
		return
	}
	max, _ := strconv.Atoi(r.URL.Query().Get("max"))
	if max <= 0 {
		max = 1
	}
	_, _ = s.box.Tick()
	results, err := s.box.RunDueContext(r.Context(), max)
	if err != nil && len(results) == 0 {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"results": results, "err": errString(err)})
}

func readJSON(r *http.Request, v any) error {
	defer r.Body.Close()
	data, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
