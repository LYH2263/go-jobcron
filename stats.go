package jobcron

// Stats 返回运行统计快照。
func (s *Scheduler) Stats() Stats {
	s.mu.Lock()
	defer s.mu.Unlock()
	st := Stats{
		Enqueued:  s.enqueued,
		Executed:  s.executed,
		Succeeded: s.succeeded,
		Failed:    s.failed,
		Dead:      s.dead,
		Closed:    s.closed,
		LastName:  s.lastName,
		LastRun:   s.lastRun,
	}
	if s.st != nil {
		if n, err := s.st.CountByStatus(string(StatusPending)); err == nil {
			st.Pending = n
		}
		if n, err := s.st.CountByStatus(string(StatusPaused)); err == nil {
			st.Paused = n
		}
	}
	return st
}
