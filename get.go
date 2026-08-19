package jobcron

// Get 按 ID 获取任务（载荷深拷贝）。
func (s *Scheduler) Get(id string) (Job, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.checkOpenLocked(); err != nil {
		return Job{}, err
	}
	rec, err := s.st.Get(id)
	if err != nil {
		return Job{}, err
	}
	return toJob(rec), nil
}

// CountByStatus 统计指定状态任务数。
func (s *Scheduler) CountByStatus(st Status) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.checkOpenLocked(); err != nil {
		return 0, err
	}
	return s.st.CountByStatus(string(st))
}
