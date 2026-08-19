package store

import (
	"encoding/json"
	"os"
)

func (m *Memory) loadLocked() error {
	data, err := os.ReadFile(m.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var snap snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return err
	}
	m.byID = make(map[string]*Record)
	m.order = snap.Order
	for i := range snap.Items {
		cp := CloneRecord(snap.Items[i])
		m.byID[cp.ID] = &cp
	}
	m.dirty = false
	return nil
}
