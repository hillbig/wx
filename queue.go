package wx

// used in builder.go Bredth First Search
type node struct {
	beg   int
	end   int
	depth int
}

// FIFO queue
type nodeQueue struct {
	beg   int
	num   int
	nodes []*node
}

func newNodeQueue(size int) *nodeQueue {
	return &nodeQueue{
		beg:   0,
		num:   0,
		nodes: make([]*node, size),
	}
}

func (q *nodeQueue) push(n *node) {
	if q.num == len(q.nodes) {
		newSize := len(q.nodes) * 2
		if newSize == 0 {
			newSize = 1
		}
		newNodes := make([]*node, newSize)
		copy(newNodes, q.nodes[q.beg:])
		copy(newNodes[len(q.nodes[q.beg:]):], q.nodes[0:q.beg])
		q.nodes = newNodes
		q.beg = 0
	}
	q.nodes[(q.beg+q.num)%len(q.nodes)] = n
	q.num++
}

func (q *nodeQueue) pop() *node {
	// no shrinkage for simplification
	n := q.nodes[q.beg]

	q.beg = (q.beg + 1) % len(q.nodes)
	q.num--
	return n
}
