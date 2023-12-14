package aocutil

// Iter copies x/exp/xiter.
type Iter[T any] func(yield func(T) bool)

// Iter2 copies x/exp/xiter.
type Iter2[T1, T2 any] func(yield func(T1, T2) bool)

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
	var vs []T
	for v := range iter {
		vs = append(vs, v)
	}
	return vs
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
