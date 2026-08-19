package store

// Delete 删除任务（管理用）。
func (m *Memory) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return ErrClosed
	}
	if _, ok := m.byID[id]; !ok {
		return ErrNotFound
	}
	delete(m.byID, id)
	for i, oid := range m.order {
		if oid == id {
			m.order = append(m.order[:i], m.order[i+1:]...)
			break
		}
	}
	m.dirty = true
	if m.path != "" {
		return m.persistLocked()
	}
	return nil
}
