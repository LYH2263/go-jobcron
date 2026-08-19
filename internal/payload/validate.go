package payload

import "errors"

var ErrEmpty = errors.New("payload: empty")

// ValidateNonEmpty 要求非空载荷。
func ValidateNonEmpty(b []byte) error {
	if len(b) == 0 {
		return ErrEmpty
	}
	return nil
}
