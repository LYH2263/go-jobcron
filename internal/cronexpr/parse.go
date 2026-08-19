package cronexpr

import (
	"fmt"
	"strings"
)

// Parse 解析五段 cron 表达式。
func Parse(expr string) (*Expr, error) {
	parts := strings.Fields(strings.TrimSpace(expr))
	if len(parts) != 5 {
		return nil, fmt.Errorf("cronexpr: want 5 fields")
	}
	minute, err := parseField(parts[0], 0, 59)
	if err != nil {
		return nil, fmt.Errorf("cronexpr: minute: %w", err)
	}
	hour, err := parseField(parts[1], 0, 23)
	if err != nil {
		return nil, fmt.Errorf("cronexpr: hour: %w", err)
	}
	dom, err := parseField(parts[2], 1, 31)
	if err != nil {
		return nil, fmt.Errorf("cronexpr: dom: %w", err)
	}
	month, err := parseField(parts[3], 1, 12)
	if err != nil {
		return nil, fmt.Errorf("cronexpr: month: %w", err)
	}
	dow, err := parseField(parts[4], 0, 6)
	if err != nil {
		return nil, fmt.Errorf("cronexpr: dow: %w", err)
	}
	return &Expr{
		Minute: minute,
		Hour:   hour,
		Dom:    dom,
		Month:  month,
		Dow:    dow,
		raw:    strings.TrimSpace(expr),
	}, nil
}
