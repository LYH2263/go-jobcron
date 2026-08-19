package store

import (
	"sync"

	"github.com/LYH2263/go-jobcron/internal/clock"
)

// Memory 线程安全内存存储，可选 JSON 快照。
type Memory struct {
	mu       sync.Mutex
	cap      int
	clk      clock.Clock
	path     string
	closed   bool
	byID     map[string]*Record
	order    []string
	dirty    bool
	flushErr error
}

// NewMemory 构造内存存储。
func NewMemory(capacity int, clk clock.Clock, persistPath string) *Memory {
	if capacity < 1 {
		capacity = 64
	}
	if clk == nil {
		clk = clock.Real{}
	}
	m := &Memory{
		cap:  capacity,
		clk:  clk,
		path: persistPath,
		byID: make(map[string]*Record),
	}
	if persistPath != "" {
		_ = m.loadLocked()
	}
	return m
}

// SetFlushError 测试注入：下次 Flush/persist 返回该错误。
func (m *Memory) SetFlushError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.flushErr = err
}

func (m *Memory) Insert(rec Record) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return "", ErrClosed
	}
	if err := validateInsert(rec); err != nil {
		return "", err
	}
	if len(m.byID) >= m.cap {
		return "", ErrCapacity
	}
	now := m.clk.Now()
	id := rec.ID
	if id == "" {
		id = NewID()
	}
	cp := CloneRecord(rec)
	cp.ID = id
	cp.CreatedAt = now
	cp.UpdatedAt = now
	if cp.NextRun.IsZero() {
		cp.NextRun = cp.RunAt
	}
	m.byID[id] = &cp
	m.order = append(m.order, id)
	m.dirty = true
	if m.path != "" {
		if err := m.persistLocked(); err != nil {
			delete(m.byID, id)
			m.order = m.order[:len(m.order)-1]
			return "", err
		}
	}
	return id, nil
}

func (m *Memory) Get(id string) (Record, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return Record{}, ErrClosed
	}
	r, ok := m.byID[id]
	if !ok {
		return Record{}, ErrNotFound
	}
	return CloneRecord(*r), nil
}

func (m *Memory) List(limit int) ([]Record, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return nil, ErrClosed
	}
	var out []Record
	for _, id := range m.order {
		r := m.byID[id]
		if r == nil {
			continue
		}
		out = append(out, CloneRecord(*r))
		if limit > 0 && len(out) >= limit {
			break
		}
	}
	return out, nil
}

func (m *Memory) Update(rec Record) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return ErrClosed
	}
	cur, ok := m.byID[rec.ID]
	if !ok {
		return ErrNotFound
	}
	cp := CloneRecord(rec)
	cp.CreatedAt = cur.CreatedAt
	cp.UpdatedAt = m.clk.Now()
	*cur = cp
	m.dirty = true
	if m.path != "" {
		if err := m.persistLocked(); err != nil {
			return err
		}
	}
	return nil
}

func (m *Memory) CountByStatus(status string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return 0, ErrClosed
	}
	n := 0
	for _, id := range m.order {
		if r := m.byID[id]; r != nil && r.Status == status {
			n++
		}
	}
	return n, nil
}

func (m *Memory) Flush() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return ErrClosed
	}
	if m.path == "" {
		return nil
	}
	return m.persistLocked()
}

func (m *Memory) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return nil
	}
	m.closed = true
	return nil
}
