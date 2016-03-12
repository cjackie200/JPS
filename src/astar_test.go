package astar

import (
	"fmt"
	"testing"
	"time"
)

func TestPNGReader(t *testing.T) {
	map_data := GetMapFromImage("../map/map800x600.png")
	if map_data == nil {
		t.Errorf("Could not open test.png")
		return
	}
	timeAstarStart := time.Now().UnixNano()

	nodes_path := Astar(map_data, 0, 0, 799, 599, true)
	timeAstarEnd := time.Now().UnixNano()

	fmt.Printf("Astar  %v ns\n", timeAstarEnd-timeAstarStart)

	fmt.Println(str_map(map_data, nodes_path))

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
