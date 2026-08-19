package queue

import (
	"container/heap"
	"sync"
	"time"
)

// Queue 按 NextRun 排序的到期小顶堆。
type Queue struct {
	mu   sync.Mutex
	h    jobHeap
	byID map[string]*Item
}

// New 构造空队列。
func New() *Queue {
	q := &Queue{byID: make(map[string]*Item)}
	heap.Init(&q.h)
	return q
}

// Push 加入或更新任务到期时间。
func (q *Queue) Push(id string, at time.Time) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if cur, ok := q.byID[id]; ok {
		cur.NextRun = at
		heap.Fix(&q.h, cur.index)
		return
	}
	it := &Item{ID: id, NextRun: at}
	heap.Push(&q.h, it)
	q.byID[id] = it
}

// Remove 从队列移除。
func (q *Queue) Remove(id string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	cur, ok := q.byID[id]
	if !ok {
		return
	}
	heap.Remove(&q.h, cur.index)
	delete(q.byID, id)
}

// PopDue 弹出所有 nextRun<=now 的 ID（不从 byID 删除，由调用方决定后续）。
func (q *Queue) PopDue(now time.Time) []string {
	q.mu.Lock()
	defer q.mu.Unlock()
	var ids []string
	for q.h.Len() > 0 && !q.h[0].NextRun.After(now) {
		it := heap.Pop(&q.h).(*Item)
		delete(q.byID, it.ID)
		ids = append(ids, it.ID)
	}
	return ids
}

// Peek 查看堆顶（可能已过期）。
func (q *Queue) Peek() (string, time.Time, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.h.Len() == 0 {
		return "", time.Time{}, false
	}
	it := q.h[0]
	return it.ID, it.NextRun, true
}

// Len 队列长度。
func (q *Queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.h.Len()
}
