package aocutil

import (
	"image"
	"slices"
	"strings"
)

// Map2D is a 2D map of bytes.
type Map2D struct {
	Data   [][]byte
	Bounds image.Rectangle
}

// NewMap2D creates a new Map2D from the given input.
func NewMap2D(input string) Map2D {
	lines := SplitLines(input)
	m := Map2D{Data: make([][]byte, len(lines))}
	for i, line := range lines {
		m.Data[i] = []byte(line)
	}
	m.Bounds = image.Rect(0, 0, len(m.Data[0]), len(m.Data))
	return m
}

// At returns the byte at the given point. If the point is out of bounds, then
// 0 is returned.
func (m Map2D) At(p image.Point) byte {
	if !p.In(m.Bounds) {
		return 0
	}
	return m.Data[p.Y][p.X]
}

// Set sets the byte at the given point.
func (m Map2D) Set(p image.Point, b byte) {
	m.Data[p.Y][p.X] = b
}

// Clone makes a copy of the map.
func (m Map2D) Clone() Map2D {
	var n Map2D
	n.Data = make([][]byte, len(m.Data))
	for i, line := range m.Data {
		n.Data[i] = slices.Clone(line)
	}
	n.Bounds = m.Bounds
	return n
}

// String returns a string representation of the map.
func (m Map2D) String() string {
	var sb strings.Builder
	sb.Grow(len(m.Data) * (len(m.Data[0]) + 1))
	for _, line := range m.Data {
		sb.Write(line)
		sb.WriteByte('\n')
	}
	return sb.String()
}

// All returns an iterator that iterates over all points in the map.
func (m Map2D) All() Iter2[image.Point, byte] {
	return m.AllWithin(m.Bounds)
}

// AllWithin returns an iterator that iterates over all points within the given
// rectangle.
func (m Map2D) AllWithin(r image.Rectangle) Iter2[image.Point, byte] {
	r = r.Canon()
	r = r.Intersect(m.Bounds)

	return func(yield func(image.Point, byte) bool) {
		for y := r.Min.Y; y < r.Max.Y; y++ {
			for x := r.Min.X; x < r.Max.X; x++ {
				if !yield(image.Pt(x, y), m.Data[y][x]) {
					return
				}
			}
		}
	}
}

// PointsWithin returns an iterator that iterates over all points within the
// given rectangle.
func PointsWithin(r image.Rectangle) Iter[image.Point] {
	return func(yield func(image.Point) bool) {
		for y := r.Min.Y; y < r.Max.Y; y++ {
			for x := r.Min.X; x < r.Max.X; x++ {
				if !yield(image.Pt(x, y)) {
					return
				}
			}
		}
	}
}

// PointsIterateDelta returns an iterator that iterates over all points within
// the given rectangle, with the given delta. The min and max are strictly
// respected without being canonized, meaning that the min point can be greater
// than the max point.
func PointsIterateDelta(r image.Rectangle, delta image.Point) Iter[image.Point] {
	return func(yield func(image.Point) bool) {
		for y := r.Min.Y; y != r.Max.Y; y += delta.Y {
			for x := r.Min.X; x != r.Max.X; x += delta.X {
				if !yield(image.Pt(x, y)) {
					return
				}
			}
		}
	}
}
