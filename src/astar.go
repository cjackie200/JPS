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

func str_map(data MapData, nodes []*Node) string {
	var result string
	for i, row := range data {
		for j, cell := range row {
			added := false
			for _, node := range nodes {
				if node.X == i && node.Y == j {
					result += "o"
					added = true
					break
				}
			}
			if !added {
				switch cell {
				case LAND:
					result += "."
				case WALL:
					result += "#"
				default:
					result += "?"
				}
			}
		}
		result += "\n"
	}
	return result
}

type Node struct {
	pos        int
	parent     *Node
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

func (n nodeList) removeNode(pos) {
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

func Heuristic(tile, stop *Node) (h int) {
	h_diag := min(abs(tile.X-stop.X), abs(tile.Y-stop.Y))
	h_stra := abs(tile.X-stop.X) + abs(tile.Y-stop.Y)
	h = COST_DIAGONAL*h_diag + COST_STRAIGHT*(h_stra-2*h_diag)
	return
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

func (w World) setPass(id int, b bool) {
	w.rwLock.Lock()
	w.pass[id] = b
	w.rwLock.Unlock()
}

func (w World) Astar(start, stop int, isChase bool) []*Node {
	closedSet := newNodeList()
	openSet := newNodeList()
	pq := newPriorityQueue()

	startNode := NewNode(start)
	openSet.addNode(startNode)
	pq.PushNode(startNode)

	for openSet.len() != 0 {
		current := pq.PopNode()
		openSet.removeNode(current.pos)
		closedSet.addNode(current)

		if current.pos == stop {
			return retracePath(current)
		} else {
			for _, adir := range adjecentDirs {
				x, y := (current.X + adir[0]), (current.Y + adir[1])

				if (x < 0) || (x >= rows) || (y < 0) || (y >= cols) {
					continue
				}

				neighbor := graph.Node(x, y)
				if neighbor == nil || closedSet.hasNode(neighbor) {
					continue
				}

				g_score := current.g + adir[2]

				if !openSet.hasNode(neighbor) {
					neighbor.parent = current
					neighbor.g = g_score
					neighbor.f = neighbor.g + Heuristic(neighbor, stop)
					openSet.addNode(neighbor)
					pq.PushNode(neighbor)
				} else if g_score < neighbor.g {
					pq.RemoveNode(neighbor)
					neighbor.parent = current
					neighbor.g = g_score
					neighbor.f = neighbor.g + Heuristic(neighbor, stop)
					pq.PushNode(neighbor)
				}

			}
		}
	}
	return nil
}

func (w World) searchJumpPoint(node *Node, dir int, stop int) (*Node, bool) {
	for next := node.pos + dir; ; next += dir {
		if next == stop || w.isJumpPoint(next, dir) {
			jump := NewNode(next)
			jump.parent = node
			jump.f = gradeDistance(next, node.pos, dir) + heuristicDistance(next, stop)
			return jump, true
		}
	}
	return nil, false
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

func (w World) isJumpPoint(pos int, dir int) bool {
	switch dir {
	case DirUp:

	case DirDown:
	case DirLeft:
	case DirRight:
	case DirLeftUp:
	case DirLeftDown:
	case DirRightUp:
	case DirRightDown:
	}
}
