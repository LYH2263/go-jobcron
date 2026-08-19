package store

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// NewID 生成任务 ID。
func NewID() string {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return fmt.Sprintf("job_%d_%s", time.Now().UnixNano(), hex.EncodeToString(b[:]))
}
