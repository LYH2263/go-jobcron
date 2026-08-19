package queue

import "time"

// Item 到期队列元素。
type Item struct {
	ID      string
	NextRun time.Time
	index   int
}
