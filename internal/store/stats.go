package store

// Stats 存储统计。
type Stats struct {
	Total   int
	Pending int
	Paused  int
	Done    int
	Failed  int
	Dead    int
}

// Snapshot 返回状态计数。
func (m *Memory) Snapshot() (Stats, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return Stats{}, ErrClosed
	}
	var st Stats
	st.Total = len(m.byID)
	for _, r := range m.byID {
		switch r.Status {
		case "pending":
			st.Pending++
		case "paused":
			st.Paused++
		case "done":
			st.Done++
		case "failed":
			st.Failed++
		case "dead":
			st.Dead++
		}
	}
	return st, nil
}
