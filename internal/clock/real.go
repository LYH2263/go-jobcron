package clock

import "time"

// Real 使用系统时钟。
type Real struct{}

func (Real) Now() time.Time { return time.Now() }
