package jps

import (
	"container/heap"
)

type PriorityQueue struct {
	pos  int
	node map[int]*Node
}

func newPriorityQueue() *PriorityQueue {
	p := new(PriorityQueue)
	p.node = make(map[int]*Node)
	return p
}

func (p PriorityQueue) Len() int {
	return len(p.node)
}

func (p PriorityQueue) Less(i, j int) bool {
	return p.node[i].f < p.node[j].f
}

func (p PriorityQueue) Swap(i, j int) {
	p.node[i], p.node[j] = p.node[j], p.node[i]
	p.node[i].heap_index = i
	p.node[j].heap_index = j
}

func (p *PriorityQueue) Push(x interface{}) {
	item, ok := x.(*Node)
	if ok {
		item.heap_index = p.pos
		p.node[p.pos] = item
		p.pos++
	}
}

func (p *PriorityQueue) Pop() interface{} {
	p.pos--
	item := p.node[p.pos]
	delete(p.node, p.pos)
	return item
}

func (p *PriorityQueue) PushNode(n *Node) {
	heap.Push(p, n)
}

func (p *PriorityQueue) PopNode() *Node {
	return heap.Pop(p).(*Node)
}

func (p *PriorityQueue) RemoveNode(n *Node) {
	heap.Remove(p, n.heap_index)
}
