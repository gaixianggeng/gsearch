package recall

type docCountSort []*queryTokenHash

func (q docCountSort) Less(i, j int) bool {
	return q[i].invertedIndex.DocsCount < q[j].invertedIndex.DocsCount
}

func (q docCountSort) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q docCountSort) Len() int {
	return len(q)
}
