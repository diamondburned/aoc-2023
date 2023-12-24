package aocutil

import "fmt"

// Interval describes an interval.
type Interval[T ScalarType] struct {
	Start T // inclusive
	End   T // inclusive
}

// String returns a string representation of the range.
func (r Interval[T]) String() string {
	return fmt.Sprintf("[%v, %v]", r.Start, r.End)
}

// Length returns the length of the range.
func (r Interval[T]) Length() T {
	return r.End - r.Start + 1
}

// Intersect returns the intersection of the two ranges.
func (r Interval[T]) Intersect(other Interval[T]) Interval[T] {
	i := Interval[T]{
		Start: max(r.Start, other.Start),
		End:   min(r.End, other.End),
	}
	i.End = max(i.Start, i.End)
	return i
}

// Except returns the ranges that are not in the other range.
// It may return a maximum of two ranges, one on the left and one on the right.
func (r Interval[T]) Except(other Interval[T]) []Interval[T] {
	var ranges []Interval[T]
	if other.Start > r.Start {
		ranges = append(ranges, Interval[T]{
			Start: r.Start,
			End:   other.Start - 1,
		})
	}
	if other.End < r.End {
		ranges = append(ranges, Interval[T]{
			Start: other.End,
			End:   r.End,
		})
	}
	return ranges
}

// Canon returns a range with the minimum value first.
func (r Interval[T]) Canon() Interval[T] {
	if r.Start > r.End {
		r.Start, r.End = r.End, r.Start
	}
	return r
}

// ContainsInterval returns true if the other range is contained in the range.
func (r Interval[T]) ContainsInterval(other Interval[T]) bool {
	return r.Start <= other.Start && other.End <= r.End
}

// Contains returns true if the value is contained in the range.
func (r Interval[T]) Contains(v T) bool {
	return r.Start <= v && v <= r.End
}
