package cronexpr

import "time"

// Next 返回 strictly after 之后的首个匹配时刻。
func Next(expr string, after time.Time) (time.Time, error) {
	e, err := Parse(expr)
	if err != nil {
		return time.Time{}, err
	}
	return e.NextAfter(after), nil
}

// NextAfter 在 expr 上找 after 之后的下一次触发（最多搜 4 年）。
func (e *Expr) NextAfter(after time.Time) time.Time {
	start := after.Truncate(time.Minute).Add(time.Minute)
	limit := start.AddDate(4, 0, 0)
	for t := start; t.Before(limit); t = t.Add(time.Minute) {
		if e.Match(t) {
			return t
		}
	}
	return time.Time{}
}
