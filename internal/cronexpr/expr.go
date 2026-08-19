package cronexpr

import "time"

// Expr 五段 cron 表达式（分 时 日 月 周）。
type Expr struct {
	Minute Field
	Hour   Field
	Dom    Field
	Month  Field
	Dow    Field
	raw    string
}

// Raw 返回原始表达式。
func (e *Expr) Raw() string {
	if e == nil {
		return ""
	}
	return e.raw
}

// Match 报告 t 是否匹配（秒与纳秒忽略，按分钟粒度）。
func (e *Expr) Match(t time.Time) bool {
	if e == nil {
		return false
	}
	t = t.Truncate(time.Minute)
	return e.Minute.match(t.Minute()) &&
		e.Hour.match(t.Hour()) &&
		e.Dom.match(t.Day()) &&
		e.Month.match(int(t.Month())) &&
		e.Dow.match(int(t.Weekday()))
}
