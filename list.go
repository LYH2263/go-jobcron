package jobcron

import "github.com/LYH2263/go-jobcron/internal/store"

// ListJobs 列出任务（载荷深拷贝）。
func (s *Scheduler) ListJobs(limit int) ([]Job, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.checkOpenLocked(); err != nil {
		return nil, err
	}
	recs, err := s.st.List(limit)
	if err != nil {
		return nil, err
	}
	out := make([]Job, len(recs))
	for i, r := range recs {
		out[i] = toJob(r)
	}
	return out, nil
}

func toJob(r store.Record) Job {
	return Job{
		ID:          r.ID,
		Name:        r.Name,
		Kind:        Kind(r.Kind),
		Payload:     r.Payload,
		Status:      Status(r.Status),
		CronExpr:    r.CronExpr,
		URL:         r.URL,
		RunAt:       r.RunAt,
		NextRun:     r.NextRun,
		Attempts:    r.Attempts,
		MaxAttempts: r.MaxAttempts,
		LastError:   r.LastError,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}
