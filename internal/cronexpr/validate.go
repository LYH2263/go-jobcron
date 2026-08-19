package cronexpr

import "strings"

// Validate 检查表达式语法。
func Validate(expr string) error {
	_, err := Parse(expr)
	return err
}

// IsSimpleEveryMinute 是否每分钟触发（测试辅助）。
func IsSimpleEveryMinute(expr string) bool {
	return strings.TrimSpace(expr) == "* * * * *"
}

// Describe 返回表达式摘要（调试）。
func Describe(expr string) (string, error) {
	e, err := Parse(expr)
	if err != nil {
		return "", err
	}
	return e.Raw(), nil
}
