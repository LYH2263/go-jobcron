package cronexpr

import (
	"fmt"
	"strconv"
	"strings"
)

// Field cron 字段匹配器。
type Field struct {
	all   bool
	items []fieldItem
}

type fieldItem struct {
	from, to, step int
}

func parseField(part string, min, max int) (Field, error) {
	part = strings.TrimSpace(part)
	if part == "" {
		return Field{}, fmt.Errorf("empty field")
	}
	if part == "*" {
		return Field{all: true}, nil
	}
	var items []fieldItem
	for _, seg := range strings.Split(part, ",") {
		seg = strings.TrimSpace(seg)
		step := 1
		if i := strings.Index(seg, "/"); i >= 0 {
			st, err := strconv.Atoi(strings.TrimSpace(seg[i+1:]))
			if err != nil || st < 1 {
				return Field{}, fmt.Errorf("bad step")
			}
			step = st
			seg = strings.TrimSpace(seg[:i])
		}
		from, to := min, max
		if seg == "*" {
			// */step already handled
		} else if strings.Contains(seg, "-") {
			bits := strings.SplitN(seg, "-", 2)
			a, err1 := strconv.Atoi(strings.TrimSpace(bits[0]))
			b, err2 := strconv.Atoi(strings.TrimSpace(bits[1]))
			if err1 != nil || err2 != nil || a > b {
				return Field{}, fmt.Errorf("bad range")
			}
			from, to = a, b
		} else {
			v, err := strconv.Atoi(seg)
			if err != nil {
				return Field{}, err
			}
			from, to = v, v
		}
		if from < min || to > max {
			return Field{}, fmt.Errorf("out of range")
		}
		items = append(items, fieldItem{from: from, to: to, step: step})
	}
	return Field{items: items}, nil
}

func (f Field) match(v int) bool {
	if f.all {
		return true
	}
	for _, it := range f.items {
		for x := it.from; x <= it.to; x += it.step {
			if x == v {
				return true
			}
		}
	}
	return false
}

// Values 返回字段在 [min,max] 内的所有匹配值（测试/调试）。
func (f Field) Values(min, max int) []int {
	if f.all {
		out := make([]int, 0, max-min+1)
		for i := min; i <= max; i++ {
			out = append(out, i)
		}
		return out
	}
	seen := map[int]bool{}
	for _, it := range f.items {
		for x := it.from; x <= it.to; x += it.step {
			if x >= min && x <= max {
				seen[x] = true
			}
		}
	}
	out := make([]int, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j] < out[i] {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
}
