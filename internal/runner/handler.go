package runner

import (
	"context"

	"github.com/LYH2263/go-jobcron/internal/store"
)

// Handler 用户注册的任务处理函数。
type Handler func(ctx context.Context, rec store.Record) error

// InvokeHandler 调用 handler，nil 时返回 ErrNilHandler。
func InvokeHandler(h Handler, ctx context.Context, rec store.Record) error {
	return h(ctx, rec)
}
