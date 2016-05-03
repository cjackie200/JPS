package jps

import (
	"sync"
)

func abs(a int) int {
	if a < 0 {
		return -a
	} else {
		return a
	}
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

const (
	COST_STRAIGHT = 1000
	COST_DIAGONAL = 1414

	MapWidth = 100
)

const (
	DirStart = iota
	DirUp
	DirDown
	DirLeft
	DirRight
	DirLeftUp
	DirLeftDown
	DirRightUp
	DirRightDown
)

var adjecentDirs = [][3]int{
	DirUp:        [3]int{-1, 0, COST_STRAIGHT},
	DirDown:      [3]int{1, 0, COST_STRAIGHT},
	DirLeft:      [3]int{0, -1, COST_STRAIGHT},
	DirRight:     [3]int{0, 1, COST_STRAIGHT},
	DirLeftUp:    [3]int{-1, -1, COST_DIAGONAL},
	DirLeftDown:  [3]int{1, -1, COST_DIAGONAL},
	DirRightUp:   [3]int{-1, 1, COST_DIAGONAL},
	DirRightDown: [3]int{1, 1, COST_DIAGONAL},
}

func str_map(data MapData, path []int) string {
	var result string
	max := len(data)
	for i := 0; i < max; i++ {
		var notPath = true
		for _, p := range path {
			if p == i {
				result += "o"
				notPath = false
				break
			}
		}
		if notPath {
			if data[i] {
				result += "."
			} else {
				result += "#"
			}
		}
		if i%MapWidth == MapWidth-1 {
			result += "\n"
		}
	}
	return result
}

type Node struct {
	pos        int
	row        int
	col        int
	parent     *Node
	dir        int
	f, g, h    int
	heap_index int
}

func NewNode(pos int) *Node {
	node := new(Node)
	node.pos = pos
	node.row = pos / MapWidth
	node.col = pos % MapWidth
	return node
}

func (n *Node) newNext(dir int) *Node {
	dirValue := adjecentDirs[dir]
	next := new(Node)
	next.row = n.row + dirValue[0]
	next.col = n.col + dirValue[1]
	next.pos = next.row*MapWidth + next.col
	next.g = n.g + dirValue[2]
	next.dir = dir
	next.parent = n
	return next
}

func (n *Node) next() {
	dirValue := adjecentDirs[n.dir]
	n.row = n.row + dirValue[0]
	n.col = n.col + dirValue[1]
	n.pos = n.row*MapWidth + n.col
	n.g = n.g + dirValue[2]
}

type nodeList map[int]*Node

func newNodeList() nodeList {
	return make(nodeList)
}

func (n nodeList) addNode(node *Node) {
	n[node.pos] = node
}

func (n nodeList) getNode(pos int) *Node {
	return n[pos]
}

func (n nodeList) removeNode(pos int) {
	delete(n, pos)
}

func (n nodeList) hasNode(pos int) bool {
	_, ok := n[pos]
	return ok
}

func (n nodeList) len() int {
	return len(n)
}

func retracePath(current_node *Node) []int {
	var path []int
	path = append(path, current_node.pos)
	for current_node.parent != nil {
		path = append(path, current_node.parent.pos)
		current_node = current_node.parent
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

type World struct {
	rwLock *sync.RWMutex
	pass   map[int]bool
}

func newWorld() *World {
	w := new(World)
	w.rwLock = new(sync.RWMutex)
	w.pass = make(map[int]bool)
	return w
}

func (w World) isPass(id int) bool {
	w.rwLock.RLock()
	b, ok := w.pass[id]
	w.rwLock.RUnlock()
	if !ok {
		b = false
	}
	return b
}

func (w World) isPassRowCol(row int, col int) bool {
	if isOutWorld(row, col) {
		return false
	} else {
		return w.isPass(row*MapWidth + col)
	}
}

func isOutWorld(row int, col int) bool {
	return row < 0 || col < 0 || row >= MapWidth || col >= MapWidth
}

func (w World) setPass(id int, b bool) {
	w.rwLock.Lock()
	w.pass[id] = b
	w.rwLock.Unlock()
}

func (w World) Jps(start int, stop int) []int {
	closedSet := newNodeList()
	pq := newPriorityQueue()
	startNode := NewNode(start)
	startNode.dir = DirStart
	pq.PushNode(startNode)
	for pq.Len() != 0 {
		current := pq.PopNode()
		if !closedSet.hasNode(current.pos) {
			closedSet.addNode(current)
			jump := w.searchJumpPoint(current, stop)
			for _, n := range jump {
				if n.pos == stop {
					return retracePath(n)
				} else {
					pq.PushNode(n)
				}
			}
		}
	}
	return []int{}
}

var jumpSearchDir = map[int]([]int){
	DirStart: []int{DirLeft, DirLeftUp, DirUp, DirRightUp, DirRight, DirRightDown, DirDown, DirLeftDown},
	DirUp:    []int{DirLeft, DirLeftUp, DirUp, DirRightUp, DirRight},

	DirDown: []int{DirLeft, DirLeftDown, DirDown, DirRightDown, DirRight},

	DirLeft: []int{DirUp, DirLeftUp, DirLeft, DirLeftDown, DirDown},

	DirRight: []int{DirUp, DirRightUp, DirRight, DirRightDown, DirDown},

	DirLeftUp: []int{DirRightUp, DirUp, DirLeftUp, DirLeft, DirLeftDown},

	DirLeftDown: []int{DirLeftUp, DirLeft, DirLeftDown, DirDown, DirRightDown},

	DirRightUp: []int{DirLeftUp, DirUp, DirRightUp, DirRight, DirRightDown},

	DirRightDown: []int{DirRightUp, DirRight, DirRightDown, DirDown, DirLeftDown},
}

func (w World) searchJumpPoint(node *Node, stop int) (jump []*Node) {
	for _, dir := range jumpSearchDir[node.dir] {
		res, ok := w.searchJumpPointDir(node, dir, stop)
		if ok {
			for _, n := range res {
				jump = append(jump, n)
			}
		}
	}
	return
}

func (w World) searchJumpPointDir(node *Node, dir int, stop int) (res []*Node, isFind bool) {
	next := node.newNext(dir)
	for {
		if next.pos == stop {
			next.h = 0
			next.f = next.g
			res = append(res, next)
			isFind = true
			return
		} else {
			if w.isPassRowCol(next.row, next.col) {
				nextJump, ok := w.nextJumpPoint(next, dir, stop)
				if ok {
					for _, n := range nextJump {
						n.h = heuristicDistance(n.pos, stop)
						n.f = n.g + n.h
						res = append(res, n)
					}
					isFind = true
					return
				}
			} else {
				return
			}
		}
		next.next()
	}
	return
}

func heuristicDistance(cur int, stop int) int {
	row := abs(cur/MapWidth - stop/MapWidth)
	col := abs(cur%MapWidth - stop%MapWidth)
	h_dia := min(row, col)
	h_str := abs(row - col)
	return COST_DIAGONAL*h_dia + COST_STRAIGHT*h_str
}

func gradeDistance(parent int, cur int, dir int) int {
	cost := adjecentDirs[dir][2]
	return (cur - parent) / dir * cost
}

func (w World) nextJumpPoint(node *Node, dir int, stop int) ([]*Node, bool) {
	switch dir {
	case DirUp:
		return w.nextJumpPointUp(node, stop)
	case DirDown:
		return w.nextJumpPointDown(node, stop)
	case DirLeft:
		return w.nextJumpPointLeft(node, stop)
	case DirRight:
		return w.nextJumpPointRight(node, stop)
	case DirLeftUp:
		return w.nextJumpPointLeftUp(node, stop)
	case DirLeftDown:
		return w.nextJumpPointLeftDown(node, stop)
	case DirRightUp:
		return w.nextJumpPointRightUp(node, stop)
	case DirRightDown:
		return w.nextJumpPointRightDown(node, stop)
	default:
		return []*Node{}, false
	}
}

func (w World) nextJumpPointUp(node *Node, stop int) (res []*Node, ok bool) {
	if w.isJumpPointUp(node) {
		res = append(res, node)
		ok = true
	}
	return
}

func (w World) nextJumpPointDown(node *Node, stop int) (res []*Node, ok bool) {
	if w.isJumpPointDown(node) {
		res = append(res, node)
		ok = true
	}
	return
}

func (w World) nextJumpPointLeft(node *Node, stop int) (res []*Node, ok bool) {
	if w.isJumpPointLeft(node) {
		res = append(res, node)
		ok = true
	}
	return
}

func (w World) nextJumpPointRight(node *Node, stop int) (res []*Node, ok bool) {
	if w.isJumpPointRight(node) {
		res = append(res, node)
		ok = true
	}
	return
}

func (w World) isJumpPointUp(cur *Node) bool {
	return w.isPassRowCol(cur.row-1, cur.col-1) && (!w.isPassRowCol(cur.row, cur.col-1)) ||
		w.isPassRowCol(cur.row-1, cur.col+1) && (!w.isPassRowCol(cur.row, cur.col+1))
}

func (w World) isJumpPointDown(cur *Node) bool {
	return w.isPassRowCol(cur.row+1, cur.col-1) && (!w.isPassRowCol(cur.row, cur.col-1)) ||
		w.isPassRowCol(cur.row+1, cur.col+1) && (!w.isPassRowCol(cur.row, cur.col+1))
}

func (w World) isJumpPointRight(cur *Node) bool {
	return w.isPassRowCol(cur.row+1, cur.col+1) && (!w.isPassRowCol(cur.row+1, cur.col)) ||
		w.isPassRowCol(cur.row-1, cur.col+1) && (!w.isPassRowCol(cur.row-1, cur.col))
}

func (w World) isJumpPointLeft(cur *Node) bool {
	return w.isPassRowCol(cur.row+1, cur.col-1) && (!w.isPassRowCol(cur.row+1, cur.col)) ||
		w.isPassRowCol(cur.row-1, cur.col-1) && (!w.isPassRowCol(cur.row-1, cur.col))
}

func (w World) isJumpPointRightUp(cur *Node) bool {
	return w.isJumpPointRight(cur) || w.isJumpPointUp(cur)
}

func (w World) isJumpPointRightDown(cur *Node) bool {
	return w.isJumpPointRight(cur) || w.isJumpPointDown(cur)
}

func (w World) isJumpPointLeftUp(cur *Node) bool {
	return w.isJumpPointLeft(cur) || w.isJumpPointUp(cur)
}

func (w World) isJumpPointLeftDown(cur *Node) bool {
	return w.isJumpPointLeft(cur) || w.isJumpPointDown(cur)
}

func (w World) nextJumpPointLeftUp(cur *Node, stop int) (res []*Node, isFind bool) {
	jumpLeft, ok1 := w.searchJumpPointDir(cur, DirLeft, stop)
	if ok1 {
		res = append(res, jumpLeft[0])
	}
	jumpUp, ok2 := w.searchJumpPointDir(cur, DirUp, stop)
	if ok2 {
		res = append(res, jumpUp[0])
	}
	if ok1 || ok2 || w.isJumpPointLeftUp(cur) {
		isFind = true
		res = append(res, cur)
	}
	return
}

func (w World) nextJumpPointLeftDown(cur *Node, stop int) (res []*Node, isFind bool) {
	jumpLeft, ok1 := w.searchJumpPointDir(cur, DirLeft, stop)
	if ok1 {
		res = append(res, jumpLeft[0])
	}
	jumpDown, ok2 := w.searchJumpPointDir(cur, DirDown, stop)
	if ok2 {
		res = append(res, jumpDown[0])
	}
	if ok1 || ok2 || w.isJumpPointLeftDown(cur) {
		isFind = true
		res = append(res, cur)
	}
	return
}

func (w World) nextJumpPointRightUp(cur *Node, stop int) (res []*Node, isFind bool) {
	jumpRight, ok1 := w.searchJumpPointDir(cur, DirRight, stop)
	if ok1 {
		res = append(res, jumpRight[0])
	}
	jumpUp, ok2 := w.searchJumpPointDir(cur, DirUp, stop)
	if ok2 {
		res = append(res, jumpUp[0])
	}
	if ok1 || ok2 || w.isJumpPointRightUp(cur) {
		isFind = true
		res = append(res, cur)
	}
	return
}

func (w World) nextJumpPointRightDown(cur *Node, stop int) (res []*Node, isFind bool) {
	jumpRight, ok1 := w.searchJumpPointDir(cur, DirRight, stop)
	if ok1 {
		res = append(res, jumpRight[0])
	}
	jumpDown, ok2 := w.searchJumpPointDir(cur, DirDown, stop)
	if ok2 {
		res = append(res, jumpDown[0])
	}
	if ok1 || ok2 || w.isJumpPointRightDown(cur) {
		isFind = true
		res = append(res, cur)
	}
	return
}
