package aocutil

import (
	"slices"
	"testing"
)

func TestHeap(t *testing.T) {
	type test struct {
		name  string
		heap  HeapType
		limit int
		in    []int
		out   []int
	}

	tests := []test{
		{
			name: "min_preorder",
			in:   []int{1, 2, 3, 4, 5},
			out:  []int{1, 2, 3, 4, 5},
		},
		{
			name: "min_reverse",
			in:   []int{5, 4, 3, 2, 1},
			out:  []int{1, 2, 3, 4, 5},
		},
		{
			name: "min_random",
			in:   []int{3, 5, 1, 4, 2},
			out:  []int{1, 2, 3, 4, 5},
		},
		{
			name: "max_preorder",
			heap: MaxHeap,
			in:   []int{5, 4, 3, 2, 1},
			out:  []int{5, 4, 3, 2, 1},
		},
		{
			name: "max_reverse",
			heap: MaxHeap,
			in:   []int{1, 2, 3, 4, 5},
			out:  []int{5, 4, 3, 2, 1},
		},
		{
			name: "max_random",
			heap: MaxHeap,
			in:   []int{3, 5, 1, 4, 2},
			out:  []int{5, 4, 3, 2, 1},
		},
		{
			name:  "min_limit",
			limit: 3,
			in:    []int{3, 5, 1, 4, 2},
			out:   []int{3, 2, 1},
		},
		{
			name:  "max_limit",
			heap:  MaxHeap,
			limit: 3,
			in:    []int{3, 5, 1, 4, 2},
			out:   []int{3, 4, 5},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := NewHeap(test.heap, test.limit, CompareOrdered[int])
			for _, v := range test.in {
				h.Push(v)
			}

			h.Sort()
			t.Logf("heap: %v", h.ToSlice())

			got := h.ToSlice()
			expect := test.out

			if !slices.Equal(got, expect) {
				t.Errorf("heap:\n"+
					"got:    %v\n"+
					"expect: %v",
					got, expect)
			}
		})
	}
}
