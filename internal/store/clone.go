package store

import "github.com/LYH2263/go-jobcron/internal/payload"

// CloneRecord 深拷贝记录（含载荷）。
func CloneRecord(r Record) Record {
	r.Payload = payload.CloneBytes(r.Payload)
	return r
}

// CloneRecords 深拷贝切片。
func CloneRecords(in []Record) []Record {
	if in == nil {
		return nil
	}
	out := make([]Record, len(in))
	for i := range in {
		out[i] = CloneRecord(in[i])
	}
	return out
}
