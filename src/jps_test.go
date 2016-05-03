package jps

import (
	"fmt"
	"testing"
	"time"
)

const (
	k  = 1000
	us = k
	ms = us * k
	s  = ms * k
)

func TestPNGReader(t *testing.T) {
	map_data := GetMapFromImage("../map/map3kx3k.png")
	if map_data == nil {
		t.Errorf("Could not open map")
		return
	}

	world := newWorld()
	world.pass = map_data

	fmt.Println(len(map_data))
	t1 := time.Now().UnixNano()
	path := world.Jps(1, 8999999)
	t2 := time.Now().UnixNano()
	printTime("jps", t1, t2)
	fmt.Printf("%#v\n", path)

	t1 = time.Now().UnixNano()
	path, ok := world.Astar(1, 8999999, true)
	t2 = time.Now().UnixNano()
	printTime("astar", t1, t2)
	fmt.Printf("%#v  %v\n", path, ok)

	// fmt.Println(str_map(map_data, path))

	// for _, n := range path {
	// fmt.Printf("%v>>>", n.pos)
	// }

	// timeAstarStart := time.Now().UnixNano()

	// nodes_path := Astar(map_data, 0, 0, 799, 599, true)
	// Astar(map_data, 0, 0, 2999, 2999, true)

	// timeAstarEnd := time.Now().UnixNano()
	// printTime("Astar", timeAstarStart, timeAstarEnd)
	// fmt.Println(str_map(map_data, nodes_path))

	// nodes_path = Astar(map_data, 35, 5, 5, 35, true)
	// fmt.Println(str_map(map_data, nodes_path))

	// nodes_path = Astar(map_data, 35, 5, 5, 5, true)
	// fmt.Println(str_map(map_data, nodes_path))
}

func printTime(str string, start int64, end int64) {
	time := end - start
	if time > s {
		fmt.Printf("%s  %v s\n", str, float64(time)/s)
	} else {
		if time > ms {
			fmt.Printf("%s  %v ms\n", str, float64(time)/ms)
		} else {
			if time > us {
				fmt.Printf("%s  %v us\n", str, float64(time)/us)
			} else {
				fmt.Printf("%s  %v ns\n", str, time)
			}
		}
	}
}

func dirToStr(dir int) string {
	switch dir {
	case DirStart:
		return "DirStart"
	case DirUp:
		return "DirUp"
	case DirDown:
		return "DirDown"
	case DirLeft:
		return "DirLeft"
	case DirRight:
		return "DirRight"
	case DirLeftUp:
		return "DirLeftUp"
	case DirLeftDown:
		return "DirLeftDown"
	case DirRightUp:
		return "DirRightUp"
	case DirRightDown:
		return "DirRightDown"
	default:
		return ""
	}
}

func (this World) Astar(startPos int, stopPos int, isChase bool) ([]int, bool) {
	if startPos == stopPos {
		return []int{}, true
	} else {
		openSet := newNodeList()
		closedSet := newNodeList()
		pq := newPriorityQueue()
		start := NewNode(startPos)
		openSet.addNode(start)
		pq.PushNode(start)
		for openSet.len() != 0 {
			// if closedSet.len() > 100000 {
			// 	return nil, false
			// } else {
			current := pq.PopNode()
			openSet.removeNode(current.pos)
			closedSet.addNode(current)
			if current.pos == stopPos {
				return retracePath(current), true
			} else {
				for dir, adir := range adjecentDirs {
					next := current.newNext(dir)
					isPass := false
					if isChase && (next.pos == stopPos) {
						isPass = true
					} else {
						isPass = this.isPassRowCol(next.row, next.col)
					}
					if isPass {
						ok := closedSet.hasNode(next.pos)
						if !ok {
							g_score := current.g + adir[2]
							ok = openSet.hasNode(next.pos)
							if ok {
								neighbor := openSet.getNode(next.pos)
								if g_score < neighbor.g {
									pq.RemoveNode(neighbor)
									neighbor.parent = current
									neighbor.g = g_score
									neighbor.f = neighbor.g + heuristicDistance(neighbor.pos, stopPos)
									pq.PushNode(neighbor)
								}
							} else {
								next.f = next.g + heuristicDistance(next.pos, stopPos)
								openSet.addNode(next)
								pq.PushNode(next)
							}
						}
					}
				}
				// }
			}
		}
		return []int{}, false
	}
}
