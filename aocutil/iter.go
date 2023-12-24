package aocutil

import "golang.org/x/exp/constraints"

// Iter copies x/exp/xiter.
type Iter[T any] func(yield func(T) bool)

// All returns all items from the iterator.
func (i Iter[T]) All() []T {
	var vs []T
	for v := range i {
		vs = append(vs, v)
	}
	return vs
}

// ContainsFunc returns true if the iterator contains a value that satisfies the
// given function.
func (i Iter[T]) ContainsFunc(eq func(T) bool) bool {
	for x := range i {
		if eq(x) {
			return true
		}
	}
	return false
}

func (i Iter[T]) Filter(f func(T) bool) Iter[T] {
	return func(yield func(T) bool) {
		for x := range i {
			if f(x) && !yield(x) {
				break
			}
		}
	}
}

// Any returns true if the iterator contains any value.
func (i Iter[T]) Any() bool {
	for _ = range i {
		return true
	}
	return false
}

// Count returns the number of items in the iterator.
func (i Iter[T]) Count() int {
	var n int
	for _ = range i {
		n++
	}
	return n
}

// Iter2 copies x/exp/xiter.
type Iter2[T1, T2 any] func(yield func(T1, T2) bool)

// All returns all items from the iterator.
func (i Iter2[T1, T2]) All() ([]T1, []T2) {
	var vs1 []T1
	var vs2 []T2
	for v1, v2 := range i {
		vs1 = append(vs1, v1)
		vs2 = append(vs2, v2)
	}
	return vs1, vs2
}

// One returns the first item from the iterator, or false if none.
func One[T any](iter Iter[T]) (T, bool) {
	for v := range iter {
		return v, true
	}
	var v T
	return v, false
}

// One2 returns the first item from the iterator, or false if none.
func One2[T1, T2 any](iter Iter2[T1, T2]) (T1, T2, bool) {
	for v1, v2 := range iter {
		return v1, v2, true
	}
	var va1 T1
	var va2 T2
	return va1, va2, false
}

// All returns all items from the iterator.
func All[T any](iter Iter[T]) []T {
	return iter.All()
}

// SliceIter returns an iterator that yields the items in the slice.
func SliceIter[T any](vs []T) Iter[T] {
	return func(yield func(T) bool) {
		for _, v := range vs {
			if !yield(v) {
				break
			}
		}
	}
}

// Range returns an iterator that yields the items in the range.
func Range[T constraints.Integer | constraints.Float](start, end T) Iter[T] {
	return func(yield func(T) bool) {
		for i := start; i < end; i++ {
			if !yield(i) {
				break
			}
		}
	}
}

// IterContains returns true if the iterator contains the value.
func IterContains[T comparable](iter Iter[T], v T) bool {
	for x := range iter {
		if x == v {
			return true
		}
	}
	return false
}

// IterateReverse returns an iterator that yields the items in the slice in
// reverse order.
func IterateReverse[T any](slice []T) Iter[T] {
	return func(yield func(T) bool) {
		for i := len(slice) - 1; i >= 0; i-- {
			if !yield(slice[i]) {
				break
			}
		}
	}
}
