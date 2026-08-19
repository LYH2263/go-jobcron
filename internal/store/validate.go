package store

import "errors"

var (
	ErrClosed    = errors.New("store: closed")
	ErrNotFound  = errors.New("store: not found")
	ErrCapacity  = errors.New("store: capacity exceeded")
	ErrBadRecord = errors.New("store: bad record")
)

func validateInsert(r Record) error {
	if r.Name == "" {
		return ErrBadRecord
	}
	if len(r.Payload) == 0 {
		return ErrBadRecord
	}
	return nil
}
