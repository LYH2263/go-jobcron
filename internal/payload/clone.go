package payload

// CloneBytes 深拷贝字节切片，避免外部别名污染存储。
func CloneBytes(b []byte) []byte {
	if b == nil {
		return nil
	}
	out := make([]byte, len(b))
	copy(out, b)
	return out
}
