package queue

type WordCnt struct {
	Word string
	Cnt  int64
}

type PQ []*WordCnt

func (pq PQ) Len() int {
	return len(pq)
}

func (pq PQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

// min-heap
func (pq PQ) Less(i, j int) bool {
	if pq[i].Cnt == pq[j].Cnt {
		return pq[i].Word > pq[j].Word
	}
	return pq[i].Cnt < pq[j].Cnt
}

func (pq *PQ) Push(x interface{}) {
	tmp := x.(*WordCnt)
	*pq = append(*pq, tmp)
}

func (pq *PQ) Pop() interface{} {
	n := len(*pq)
	tmp := (*pq)[n-1]
	*pq = (*pq)[:n-1]
	return tmp
}
