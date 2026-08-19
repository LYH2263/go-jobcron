package runner

import (
	"github.com/LYH2263/go-jobcron/internal/schedule"
	"github.com/LYH2263/go-jobcron/internal/store"
)

// MarkFail 标记失败并安排重试或 dead。
func MarkFail(eng *Engine, rec store.Record, runErr error) (dead bool, err error) {
	rec.LastError = runErr.Error()
	if rec.Attempts >= rec.MaxAttempts {
		rec.Status = "dead"
		dead = true
		eng.q.Remove(rec.ID)
	} else {
		rec.Status = "pending"
		delay := Delay(eng.backoffBase, eng.backoffCap, eng.backoffFactor, rec.Attempts)
		rec.NextRun = eng.clk.Now().Add(delay)
		eng.q.Push(rec.ID, rec.NextRun)
	}
	rec.UpdatedAt = eng.clk.Now()
	if err := eng.st.Update(rec); err != nil {
		return dead, err
	}
	if rec.Kind == "cron" && dead {
		// cron dead 不再重调度
		return dead, nil
	}
	if rec.Kind == "cron" && !dead {
		next, nerr := schedule.RescheduleCron(rec.CronExpr, eng.clk.Now())
		if nerr == nil && !next.IsZero() {
			rec.NextRun = next
			_ = eng.st.Update(rec)
			eng.q.Push(rec.ID, next)
		}
	}
	return dead, nil
}
