package astar

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
	timeAstarStart := time.Now().UnixNano()

	// nodes_path := Astar(map_data, 0, 0, 799, 599, true)
	Astar(map_data, 0, 0, 2999, 2999, true)

	timeAstarEnd := time.Now().UnixNano()
	printTime("Astar", timeAstarStart, timeAstarEnd)
	// fmt.Println(str_map(map_data, nodes_path))

	// nodes_path = Astar(map_data, 35, 5, 5, 35, true)
	// fmt.Println(str_map(map_data, nodes_path))

	// nodes_path = Astar(map_data, 35, 5, 5, 5, true)
	// fmt.Println(str_map(map_data, nodes_path))
}

func BenchmarkAstar4Dirs100x100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		map_data := NewMapData(100, 100)
		Astar(map_data, 0, 0, 99, 99, false)
	}
}

func BenchmarkAstar8Dirs100x100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		map_data := NewMapData(100, 100)
		Astar(map_data, 0, 0, 99, 99, false)
	}
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
