package runner

import "context"

// RunDueBatch 批量执行（内部辅助）。
func RunDueBatch(ctx context.Context, eng *Engine, max int) ([]Result, error) {
	var out []Result
	for i := 0; max <= 0 || i < max; i++ {
		if err := ctx.Err(); err != nil {
			return out, err
		}
		res, ok, err := eng.RunOne(ctx, eng.clk.Now())
		if !ok {
			return out, nil
		}
		out = append(out, res)
		if err != nil {
			return out, err
		}
	}
	return out, nil
}
