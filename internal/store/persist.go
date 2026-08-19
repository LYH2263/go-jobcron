package store

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type snapshot struct {
	Order []string `json:"order"`
	Items []Record `json:"items"`
}

func (m *Memory) persistLocked() error {
	if m.flushErr != nil {
		err := m.flushErr
		m.flushErr = nil
		return err
	}
	items := make([]Record, 0, len(m.order))
	for _, id := range m.order {
		if r := m.byID[id]; r != nil {
			items = append(items, CloneRecord(*r))
		}
	}
	snap := snapshot{Order: append([]string(nil), m.order...), Items: items}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Dir(m.path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	tmp := m.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, m.path)
}
