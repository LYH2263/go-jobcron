package store

import "time"

// ListByStatus 按状态列出。
func (m *Memory) ListByStatus(status string, limit int) ([]Record, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return nil, ErrClosed
	}
	var out []Record
	for _, id := range m.order {
		r := m.byID[id]
		if r == nil || r.Status != status {
			continue
		}
		out = append(out, CloneRecord(*r))
		if limit > 0 && len(out) >= limit {
			break
		}
	}
	return out, nil
}

// ListDue 列出 pending 且 nextRun<=now 的任务。
func (m *Memory) ListDue(now time.Time, limit int) ([]Record, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return nil, ErrClosed
	}
	var out []Record
	for _, id := range m.order {
		r := m.byID[id]
		if r == nil || r.Status != "pending" {
			continue
		}
		if r.NextRun.After(now) {
			continue
		}
		out = append(out, CloneRecord(*r))
		if limit > 0 && len(out) >= limit {
			break
		}
	}
	return out, nil
}
