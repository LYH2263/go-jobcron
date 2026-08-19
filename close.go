package jobcron

// Close 关闭调度器：先 Flush 持久化，再 Close store，最后释放队列与执行器引用。
func (s *Scheduler) Close() error {
	if s == nil {
		return ErrNilScheduler
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil
	}
	s.closed = true

	var first error
	if s.st != nil {
		if err := s.st.Close(); err != nil && first == nil {
			first = err
		}
		if err := s.st.Flush(); err != nil && first == nil {
			first = err
		}
	}
	s.st = nil
	s.q = nil
	s.eng = nil
	return first
}
