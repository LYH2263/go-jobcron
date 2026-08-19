package runner

import (
	"context"
	"errors"
	"time"

	"github.com/LYH2263/go-jobcron/internal/schedule"
	"github.com/LYH2263/go-jobcron/internal/store"
)

// Execute 执行单条任务并更新状态。
func Execute(ctx context.Context, eng *Engine, rec store.Record) (Result, error) {
	start := time.Now()
	res := Result{
		JobID:    rec.ID,
		Name:     rec.Name,
		Attempts: rec.Attempts + 1,
	}
	rec.Status = "running"
	rec.Attempts++
	_ = eng.st.Update(rec)

	var runErr error
	var code int
	switch rec.Kind {
	case "http_callback":
		code, runErr = HTTPCallback(ctx, eng.client, rec, eng.userAgent)
		res.StatusCode = code
	case "handler":
		runErr = InvokeHandler(eng.handler, ctx, rec)
	default:
		runErr = InvokeHandler(eng.handler, ctx, rec)
	}
	res.Duration = time.Since(start)
	res.FinishedAt = eng.clk.Now()

	if runErr == nil {
		res.OK = true
		if err := MarkDone(eng, rec); err != nil {
			res.OK = false
			res.Err = err.Error()
			return res, err
		}
		return res, nil
	}
	res.Err = runErr.Error()
	dead, err := MarkFail(eng, rec, runErr)
	if err != nil {
		return res, err
	}
	if dead {
		res.Err = "dead: " + res.Err
	}
	return res, WrapErr(ErrRun, runErr)
}

// MarkDone 标记完成；persist 失败时回滚内存状态。
func MarkDone(eng *Engine, rec store.Record) error {
	prevStatus := rec.Status
	prevAttempts := rec.Attempts
	rec.Status = "done"
	rec.LastError = ""
	rec.UpdatedAt = eng.clk.Now()
	if rec.Kind == "cron" {
		next, err := schedule.RescheduleCron(rec.CronExpr, eng.clk.Now())
		if err != nil {
			rec.Status = prevStatus
			rec.Attempts = prevAttempts
			return err
		}
		rec.Status = "pending"
		rec.NextRun = next
		rec.Attempts = 0
		if err := eng.st.Update(rec); err != nil {
			rec.Status = prevStatus
			rec.Attempts = prevAttempts
			return err
		}
		eng.q.Push(rec.ID, next)
		return nil
	}
	if err := eng.st.Update(rec); err != nil {
		rec.Status = prevStatus
		rec.Attempts = prevAttempts
		return err
	}
	eng.q.Remove(rec.ID)
	return nil
}

func mapCtx(err error) error {
	if errors.Is(err, context.Canceled) {
		return ErrCanceled
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return ErrTimeout
	}
	return err
}
