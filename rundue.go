package jobcron

import (
	"context"
	"errors"

	"github.com/LYH2263/go-jobcron/internal/runner"
)

// RunDue 执行一条到期任务。
func (s *Scheduler) RunDue(ctx context.Context) (RunResult, error) {
	results, err := s.RunDueContext(ctx, 1)
	if err != nil && len(results) == 0 {
		return RunResult{}, err
	}
	if len(results) == 0 {
		return RunResult{}, nil
	}
	return results[0], err
}

// RunDueContext 在 ctx 取消前最多执行 max 条到期任务（max<=0 不限但仍尊重取消）。
func (s *Scheduler) RunDueContext(ctx context.Context, max int) ([]RunResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out []RunResult
	for {
		// ctx 已取消（如 SIGTERM）时立即退出，不再 claim 到期任务。
		if err := ctx.Err(); err != nil {
			return out, mapCtxErr(err)
		}
		if max > 0 && len(out) >= max {
			return out, nil
		}
		s.mu.Lock()
		if err := s.checkOpenLocked(); err != nil {
			s.mu.Unlock()
			if len(out) == 0 {
				return out, err
			}
			return out, err
		}
		eng := s.eng
		clk := s.clk
		s.mu.Unlock()

		res, ok, err := eng.RunOne(ctx, clk.Now())
		if err != nil {
			if errors.Is(err, ErrClosed) || errors.Is(err, ErrCanceled) || errors.Is(err, context.Canceled) {
				return out, mapCtxErr(err)
			}
			if !ok {
				return out, err
			}
		}
		if !ok {
			return out, nil
		}
		rr := toRunResult(res)
		out = append(out, rr)
		s.mu.Lock()
		s.executed++
		s.lastRun = rr.FinishedAt
		s.lastName = rr.Name
		if rr.OK {
			s.succeeded++
		} else {
			s.failed++
		}
		s.mu.Unlock()
		if err != nil {
			return out, err
		}
		if !rr.OK && rr.Attempts > 0 {
			s.mu.Lock()
			base, cap, factor := s.backoffBase, s.backoffCap, s.backoffFactor
			s.mu.Unlock()
			if err := runner.Wait(ctx, base, cap, factor, rr.Attempts); err != nil {
				return out, mapCtxErr(err)
			}
		}
	}
}

func toRunResult(r runner.Result) RunResult {
	return RunResult{
		JobID:      r.JobID,
		Name:       r.Name,
		OK:         r.OK,
		Attempts:   r.Attempts,
		StatusCode: r.StatusCode,
		Duration:   r.Duration,
		Err:        r.Err,
		FinishedAt: r.FinishedAt,
	}
}

func mapCtxErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.Canceled) {
		return ErrCanceled
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return ErrTimeout
	}
	return err
}
