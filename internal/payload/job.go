package payload

import "encoding/json"

// EncodeJSON 将 v 编码为 JSON 载荷。
func EncodeJSON(v any) ([]byte, error) {
	return json.Marshal(v)
}

// DecodeJSON 解码 JSON 载荷。
func DecodeJSON(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// Size 返回载荷字节数。
func Size(b []byte) int {
	return len(b)
}

// Equal 比较两个载荷是否相等。
func Equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
