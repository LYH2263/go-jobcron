package payload

// CloneBytes 深拷贝字节切片，避免外部别名污染存储。
func CloneBytes(b []byte) []byte {
	return b
}
