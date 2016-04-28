package astar

import (
	"fmt"
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

	MapWidth = 10
)

const (
	DirStart     = 0
	DirUp        = -MapWidth
	DirDown      = MapWidth
	DirLeft      = -1
	DirRight     = 1
	DirLeftUp    = -MapWidth - 1
	DirLeftDown  = MapWidth - 1
	DirRightUp   = -MapWidth + 1
	DirRightDown = MapWidth + 1
)

var adjecentDirs = initDirCost()

func initDirCost() map[int]int {
	return map[int]int{
		DirUp:        COST_STRAIGHT,
		DirDown:      COST_STRAIGHT,
		DirLeft:      COST_STRAIGHT,
		DirRight:     COST_STRAIGHT,
		DirLeftUp:    COST_DIAGONAL,
		DirLeftDown:  COST_DIAGONAL,
		DirRightUp:   COST_DIAGONAL,
		DirRightDown: COST_DIAGONAL,
	}
}

func str_map(data MapData) string {
	var result string
	max := len(data)
	for i := 0; i < max; i++ {
		if i%10 == 0 {
			result += "\n"
		}
		if data[i] {
			result += "."
		} else {
			result += "#"
		}
	}
	return result
}

type Node struct {
	pos        int
	parent     *Node
	dir        int
	f, g, h    int
	heap_index int
}

func NewNode(pos int) *Node {
	node := new(Node)
	node.pos = pos
	return node
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

func retracePath(current_node *Node) []*Node {
	var path []*Node
	path = append(path, current_node)
	for current_node.parent != nil {
		path = append(path, current_node.parent)
		current_node = current_node.parent
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

// func Heuristic(tile, stop *Node) (h int) {
// 	h_diag := min(abs(tile.X-stop.X), abs(tile.Y-stop.Y))
// 	h_stra := abs(tile.X-stop.X) + abs(tile.Y-stop.Y)
// 	h = COST_DIAGONAL*h_diag + COST_STRAIGHT*(h_stra-2*h_diag)
// 	return
// }

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

func (w World) setPass(id int, b bool) {
	w.rwLock.Lock()
	w.pass[id] = b
	w.rwLock.Unlock()
}

func (w World) Jps(start, stop int) []*Node {
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
	return []*Node{}
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
		fmt.Println(dir, res, ok)
		if ok {
			for _, n := range res {
				jump = append(jump, n)
			}
		}
	}
	return
}

func (w World) searchJumpPointDir(node *Node, dir int, stop int) (res []*Node, isFind bool) {
	gVale := adjecentDirs[dir]
	next := NewNode(node.pos + dir)
	next.g = node.g + gVale
	next.dir = dir
	next.parent = node
	for {
		if next.pos == stop {
			next.h = heuristicDistance(next.pos, stop)
			next.f = next.g + next.f
			res = append(res, next)
			isFind = true
			return
		} else {
			if w.isPass(next.pos) {
				nextJump, ok := w.nextJumpPoint(next, dir, stop)
				if ok {
					for _, n := range nextJump {
						res = append(res, n)
					}
					isFind = true
					return
				}
			} else {
				return
			}
		}
		next.pos += dir
		next.g += gVale
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
	cost := adjecentDirs[dir]
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
	if w.isJumpPointUp(node.pos) {
		res = append(res, node)
		ok = true
	}
	return
}

func (w World) nextJumpPointDown(node *Node, stop int) (res []*Node, ok bool) {
	if w.isJumpPointDown(node.pos) {
		res = append(res, node)
		ok = true
	}
	return
}

func (w World) nextJumpPointLeft(node *Node, stop int) (res []*Node, ok bool) {
	if w.isJumpPointLeft(node.pos) {
		res = append(res, node)
		ok = true
	}
	return
}

func (w World) nextJumpPointRight(node *Node, stop int) (res []*Node, ok bool) {
	if w.isJumpPointRight(node.pos) {
		res = append(res, node)
		ok = true
	}
	return
}

func (w World) isJumpPointUp(cur int, stop int) bool {
	return cur == stop ||
		w.isPass(cur+DirLeftUp) && (!w.isPass(cur+DirLeft)) ||
		w.isPass(cur+DirRightUp) && (!w.isPass(cur+DirRight))
}

func (w World) isJumpPointDown(cur int, stop int) bool {
	return cur == stop ||
		w.isPass(cur+DirLeftDown) && (!w.isPass(cur+DirLeft)) ||
		w.isPass(cur+DirRightDown) && (!w.isPass(cur+DirRight))
}

func (w World) isJumpPointRight(cur int, stop int) bool {
	return cur == stop ||
		w.isPass(cur+DirRightUp) && (!w.isPass(cur+DirUp)) ||
		w.isPass(cur+DirRightDown) && (!w.isPass(cur+DirDown))
}

func (w World) isJumpPointLeft(cur int, stop int) bool {
	return cur == stop ||
		w.isPass(cur+DirLeftUp) && (!w.isPass(cur+DirUp)) ||
		w.isPass(cur+DirLeftDown) && (!w.isPass(cur+DirDown))
}

func (w World) isJumpPointRightUp(cur *Node, stop int) bool {
	return w.isJumpPointRight(cur, stop) || w.isJumpPointUp(cur, stop)
}

func (w World) isJumpPointRightDown(cur *Node, stop int) bool {
	return w.isJumpPointRight(cur, stop) || w.isJumpPointDown(cur, stop)
}

func (w World) isJumpPointLeftUp(cur *Node, stop int) bool {
	return w.isJumpPointRight(cur, stop) || w.isJumpPointUp(cur, stop)
}
func (w World) isJumpPointLeftDown(cur *Node, stop int) bool {
	return w.isJumpPointLeft(cur, stop) || w.isJumpPointDown(cur, stop)
}

func (w World) nextJumpPointLeftDown(cur *Node, stop int) (res []*Node, isFind bool) {
	jumpRight, ok1 := w.searchJumpPointDir(cur, DirRight, stop)
	if ok1 {
		isFind = true
		res = append(res, jr[0])

	}
	jumpUp, ok2 := w.searchJumpPointDir(cur, DirUp, stop)
	if ok2 {
		isFind = true
		res = append(res, ju[0])
	}
	if w.isJumpPointRightUp(cur, stop) {
		isFind = true
		res = append(res, cur)
	}
	return
}

func (w World) nextJumpPointRightDown(cur *Node, stop int) (res []*Node, isFind bool) {
	jumpRight, ok1 := w.searchJumpPointDir(cur, DirRight, stop)
	if ok1 {
		isFind = true
		res = append(res, jr[0])
	}

	jumpDown, ok2 := w.searchJumpPointDir(cur, DirDown, stop)
	if ok2 {
		isFind = true
		res = append(res, jd[0])
	}
	if w.isJumpPointRightDown(cur, stop) {
		isFind = true
		res = append(res, cur)
	}
	return
}

func (w World) nextJumpPointLeftUp(cur *Node, stop int) (res []*Node, ok bool) {
	jumpLeft, ok1 := w.searchJumpPointDir(cur, DirLeft, stop)
	if ok {
		for _, jl := range jumpLeft {
			res = append(res, jl)
		}
	}

	jumpUp, ok2 := w.searchJumpPointDir(cur, DirUp, stop)
	if ok {
		for _, ju := range jumpUp {
			res = append(res, ju)
		}
	}
	ok = ok1 || ok2
	return
}

func (w World) nextJumpPointLeftDown(cur *Node, stop int) (res []*Node, ok bool) {
	jumpLeft, ok1 := w.searchJumpPointDir(cur, DirLeft, stop)
	if ok {
		for _, jl := range jumpLeft {
			res = append(res, jl)
		}
	}

	jumpDown, ok2 := w.searchJumpPointDir(cur, DirDown, stop)
	if ok {
		for _, jd := range jumpDown {
			res = append(res, jd)
		}
	}
	ok = ok1 || ok2
	return
}

// func makeJumpPoint(cur c, parent int, dir int, stop int) *Node {
// 	jump := NewNode(cur)
// 	jump.parent = parent
// 	jump.dir = dir
// 	jump.g = gradeDistance(cur, parent, dir)
// 	jump.h = heuristicDistance(cur, stop)
// 	jump.f = jump.g + jump.h
// 	return jump
// }
