package aocutil

import (
	"container/heap"
	"sort"

	"golang.org/x/exp/constraints"
)

// Heap is a heap of ordered values.
type Heap[T any] internalHeap[T]

type HeapType int

const (
	MinHeap HeapType = iota
	MaxHeap
)

// NewMinHeap returns a new min heap.
func NewMinHeap[T constraints.Ordered](limit int) *Heap[T] {
	return NewHeap(MinHeap, limit, CompareOrdered[T])
}

// NewMaxHeap returns a new max heap.
func NewMaxHeap[T constraints.Ordered](limit int) *Heap[T] {
	return NewHeap(MaxHeap, limit, CompareOrdered[T])
}

// NewHeap creates a new heap with options.
func NewHeap[T any](t HeapType, limit int, compare func(a, b T) int) *Heap[T] {
	if limit > 0 {
		// If we have a limit, then flip the heap type. Not entirely sure why we
		// need this, though!
		t = HeapType((int(t) + 1) % 2)
	}
	if t == MaxHeap {
		compare = ReverseCompare(compare)
	}

	slice := make([]T, 0, 12)
	h := internalHeap[T]{
		cmp:   compare,
		heap:  slice,
		htyp:  t,
		limit: limit,
	}
	heap.Init((*internalHeap[T])(&h))
	return (*Heap[T])(&h)
}

// Push pushes a value onto the heap.
func (h *Heap[T]) Push(v T) {
	if h.limit == 0 || len(h.heap) < h.limit {
		h.push(v)
		return
	}

	old := h.heap[0]
	if h.cmp(v, old) > 0 {
		h.heap[0] = v
		h.Fix(0)
	}
}

func (h *Heap[T]) push(v T) {
	heap.Push((*internalHeap[T])(h), v)
}

// Pop pops a value from the heap.
func (h *Heap[T]) Pop() T {
	return heap.Pop((*internalHeap[T])(h)).(T)
}

// Fix fixes the heap after a value has been modified at the given index.
func (h *Heap[T]) Fix(i int) {
	heap.Fix((*internalHeap[T])(h), i)
}

// Len returns the number of items in the heap.
func (h *Heap[T]) Len() int {
	return len(h.heap)
}

// Sort sorts the internal heap slice
func (h *Heap[T]) Sort() {
	sort.Sort((*internalHeap[T])(h))
}

// ToSlice returns the underlying heap slice. The returned slice is not sorted.
func (h *Heap[T]) ToSlice() []T {
	return h.heap
}

type internalHeap[T any] struct {
	cmp   func(a, b T) int
	heap  []T
	htyp  HeapType
	limit int
}

var _ heap.Interface = (*internalHeap[int])(nil)
var _ sort.Interface = (*internalHeap[int])(nil)

func (h internalHeap[T]) Len() int {
	return len(h.heap)
}

func (h internalHeap[T]) Less(i, j int) bool {
	return h.cmp(h.heap[i], h.heap[j]) < 0
}

func (h internalHeap[T]) Swap(i, j int) {
	h.heap[i], h.heap[j] = h.heap[j], h.heap[i]
}

// Push pushes the element x onto the heap.
func (h *internalHeap[T]) Push(x interface{}) {
	h.heap = append(h.heap, x.(T))
}

// Pop removes the last element from the heap and returns it.
func (h *internalHeap[T]) Pop() interface{} {
	x := h.heap[len(h.heap)-1]
	h.heap = h.heap[:len(h.heap)-1]
	return x
}
