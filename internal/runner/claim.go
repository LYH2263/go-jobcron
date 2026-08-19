package runner

import (
	"time"

	"github.com/LYH2263/go-jobcron/internal/store"
)

// ClaimNext 从队列/store 认领下一条到期 pending 任务。
func ClaimNext(st store.Store, q interface {
	PopDue(time.Time) []string
}, now time.Time) (store.Record, bool, error) {
	ids := q.PopDue(now)
	for _, id := range ids {
		rec, err := st.Get(id)
		if err != nil {
			continue
		}
		if rec.Status != "pending" {
			continue
		}
		if rec.NextRun.After(now) {
			continue
		}
		return rec, true, nil
	}
	recs, err := st.ListDue(now, 1)
	if err != nil {
		return store.Record{}, false, err
	}
	if len(recs) == 0 {
		return store.Record{}, false, nil
	}
	return recs[0], true, nil
}
