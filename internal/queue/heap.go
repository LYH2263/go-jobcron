package queue

type jobHeap []*Item

func (h jobHeap) Len() int { return len(h) }

func (h jobHeap) Less(i, j int) bool {
	if h[i].NextRun.Equal(h[j].NextRun) {
		return h[i].ID < h[j].ID
	}
	return h[i].NextRun.Before(h[j].NextRun)
}

func (h jobHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *jobHeap) Push(x any) {
	it := x.(*Item)
	it.index = len(*h)
	*h = append(*h, it)
}

func (h *jobHeap) Pop() any {
	old := *h
	n := len(old)
	it := old[n-1]
	old[n-1] = nil
	it.index = -1
	*h = old[:n-1]
	return it
}
