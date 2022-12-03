package aocutil

import (
	"reflect"
	"testing"
)

func TestSlidingWindow(t *testing.T) {
	var num = []int{0, 1, 2, 3, 4}
	var res = [][]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}}
	var turn int

	SlidingWindow(num, 2, func(w []int) {
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

func TestHeap(t *testing.T) {
	type test struct {
		name string
		opts HeapOpts[int]
		in   []int
		out  []int
	}

	tests := []test{
		{
			name: "min_preorder",
			opts: HeapOpts[int]{},
			in:   []int{1, 2, 3, 4, 5},
			out:  []int{1, 2, 3, 4, 5},
		},
		{
			name: "min_reverse",
			opts: HeapOpts[int]{},
			in:   []int{5, 4, 3, 2, 1},
			out:  []int{1, 2, 3, 4, 5},
		},
		{
			name: "min_random",
			opts: HeapOpts[int]{},
			in:   []int{3, 5, 1, 4, 2},
			out:  []int{1, 2, 3, 4, 5},
		},
		{
			name: "max_preorder",
			opts: HeapOpts[int]{Max: true},
			in:   []int{5, 4, 3, 2, 1},
			out:  []int{5, 4, 3, 2, 1},
		},
		{
			name: "max_reverse",
			opts: HeapOpts[int]{Max: true},
			in:   []int{1, 2, 3, 4, 5},
			out:  []int{5, 4, 3, 2, 1},
		},
		{
			name: "max_random",
			opts: HeapOpts[int]{Max: true},
			in:   []int{3, 5, 1, 4, 2},
			out:  []int{5, 4, 3, 2, 1},
		},
		{
			name: "min_limit",
			opts: HeapOpts[int]{Limit: 3},
			in:   []int{3, 5, 1, 4, 2},
			out:  []int{1, 2, 3},
		},
		{
			name: "max_limit",
			opts: HeapOpts[int]{Max: true, Limit: 3},
			in:   []int{3, 5, 1, 4, 2},
			out:  []int{5, 4, 3},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := NewHeapOpts(test.opts)
			for _, v := range test.in {
				h.Push(v)
			}

			t.Logf("heap: %v", h.ToSlice())

			for _, v := range test.out {
				if h.Pop() != v {
					t.Errorf("mismatch: %v", h)
				}
			}
		})
	}
}
