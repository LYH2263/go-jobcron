package jobcron

import "time"

// Tick 用当前时钟推进到期队列（不执行）。
func (s *Scheduler) Tick() (int, error) {
	return s.TickAt(s.now())
}

// TickAt 在指定时刻标记到期（返回到期 ID 数，仅用于观测/测试）。
func (s *Scheduler) TickAt(t time.Time) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	due := s.q.PopDue(t)
	return len(due), nil
}
