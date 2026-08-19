package clock

import "time"

// Clock 抽象时间源。
type Clock interface {
	Now() time.Time
}
