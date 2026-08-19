package queue

import "time"

// DueBefore 返回 nextRun<=t 的 ID 列表但不弹出（快照）。
func (q *Queue) DueBefore(t time.Time) []string {
	q.mu.Lock()
	defer q.mu.Unlock()
	var out []string
	for _, it := range q.h {
		if !it.NextRun.After(t) {
			out = append(out, it.ID)
		}
	}
	return out
}

// Contains 是否包含 id。
func (q *Queue) Contains(id string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	_, ok := q.byID[id]
	return ok
}
