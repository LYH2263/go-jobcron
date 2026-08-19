package jobcron

// Pause 暂停任务。
func (s *Scheduler) Pause(id string) error {
	return s.setPaused(id, true)
}

// Resume 恢复任务。
func (s *Scheduler) Resume(id string) error {
	return s.setPaused(id, false)
}

func (s *Scheduler) setPaused(id string, pause bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.checkOpenLocked(); err != nil {
		return err
	}
	rec, err := s.st.Get(id)
	if err != nil {
		return err
	}
	if pause {
		if rec.Status == string(StatusDone) || rec.Status == string(StatusDead) {
			return ErrInvalidStatus
		}
		rec.Status = string(StatusPaused)
		s.q.Remove(id)
	} else {
		if rec.Status != string(StatusPaused) {
			return ErrInvalidStatus
		}
		rec.Status = string(StatusPending)
		s.q.Push(id, rec.NextRun)
	}
	return s.st.Update(rec)
}
