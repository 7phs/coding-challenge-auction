package models

type bidHeap []BidRecI

func (h bidHeap) Len() int { return len(h) }

func (h bidHeap) Less(i, j int) bool {
	if h[i].Bid() < h[j].Bid() {
		return true
	}

	if h[i].Bid() > h[j].Bid() {
		return false
	}

	return h[i].Updated()>h[i].Updated()
}

func (h bidHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *bidHeap) Push(x interface{}) {
	*h = append(*h, x.(BidRecI))
}

func (h *bidHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
