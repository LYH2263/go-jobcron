package runner

import (
	"errors"
	"fmt"
)

var (
	ErrRun        = errors.New("runner: run failed")
	ErrHTTP       = errors.New("runner: http error")
	ErrNilHandler = errors.New("runner: nil handler")
	ErrNilClient  = errors.New("runner: nil http client")
	ErrTimeout    = errors.New("runner: timeout")
	ErrCanceled   = errors.New("runner: canceled")
)

// Wrap 用 %w 包裹哨兵，便于 errors.Is。
func Wrap(sentinel error, msg string) error {
	if sentinel == nil {
		return fmt.Errorf("%s", msg)
	}
	return fmt.Errorf("%w: %s", sentinel, msg)
}

// WrapErr 用 %w 链上底层错误。
func WrapErr(sentinel, err error) error {
	if err == nil {
		return sentinel
	}
	if sentinel == nil {
		return err
	}
	return fmt.Errorf("%w: %w", sentinel, err)
}
