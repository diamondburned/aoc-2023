package aocutil

import (
	"reflect"
	"testing"
)

func TestWindow(t *testing.T) {
	var num = []int{0, 1, 2, 3, 4}
	var res = [][]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}}
	var turn int

	Window(num, 2, func(w []int) {
		if !reflect.DeepEqual(res[turn], w) {
			t.Errorf("mismatch at %d: %v", turn, w)
		}
		turn++
	})
}

func TestChunk(t *testing.T) {
	var num = []int{0, 1, 2, 3, 4}
	var res = [][]int{{0, 1}, {2, 3}, {4}}
	var turn int

	Chunk(num, 2, func(w []int) {
		if !reflect.DeepEqual(res[turn], w) {
			t.Errorf("mismatch at %d: %v", turn, w)
		}
		turn++
	})
}
