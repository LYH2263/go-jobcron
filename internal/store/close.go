package store

// Closed 是否已关闭。
func (m *Memory) Closed() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.closed
}
