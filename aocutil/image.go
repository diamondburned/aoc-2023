package aocutil

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/exp/constraints"
)

// SaveImage saves the given image to a PNG image.
func SaveImage(img image.Image, dst string) {
	if ext := filepath.Ext(dst); ext != ".png" {
		panic(fmt.Errorf("invalid extension %q, only .png is supported", ext))
	}
	f := E2(os.Create(dst))
	defer f.Close()
	E1(png.Encode(f, img))
}

// OpenImage saves the given image to a temporary PNG image and opens it.
func OpenImage(img image.Image) {
	f := E2(os.CreateTemp("", "aocutil-*.png"))
	defer f.Close()
	E1(png.Encode(f, img))
	log.Printf("opening PNG at %q", f.Name())
	cmd := exec.Command("xdg-open", f.Name())
	cmd.Stderr = os.Stderr
	E1(cmd.Start())
}

// ScalarType is a type that can be used as a scalar.
type ScalarType interface {
	constraints.Float | constraints.Signed
}

// Point describes a point in 2D space.
type Point[T ScalarType] struct {
	X, Y T
}

// Pt returns a point in 2D space with the given coordinates.
func Pt[T ScalarType](x, y T) Point[T] {
	return Point[T]{x, y}
}

// String returns a string representation of the point.
func (p Point[T]) String() string {
	return fmt.Sprintf("(%v,%v)", p.X, p.Y)
}

// Add returns the sum of the two points.
func (p Point[T]) Add(q Point[T]) Point[T] {
	return Point[T]{p.X + q.X, p.Y + q.Y}
}

// Sub returns the difference of the two points.
func (p Point[T]) Sub(q Point[T]) Point[T] {
	return Point[T]{p.X - q.X, p.Y - q.Y}
}

// Mul returns the product of the point and the scalar.
func (p Point[T]) Mul(s T) Point[T] {
	return Point[T]{p.X * s, p.Y * s}
}

// Div returns the quotient of the point and the scalar.
func (p Point[T]) Div(s T) Point[T] {
	return Point[T]{p.X / s, p.Y / s}
}

// Manhattan returns the Manhattan distance between the two points.
func (p Point[T]) Manhattan(q Point[T]) T {
	return Abs(p.X-q.X) + Abs(p.Y-q.Y)
}

// Det returns the determinant of the two points.
func (p Point[T]) Det(q Point[T]) T {
	return p.X*q.Y - p.Y*q.X
}

// Dot returns the dot product of the two points.
func (p Point[T]) Dot(q Point[T]) T {
	return p.X*q.X + p.Y*q.Y
}

// Norm2 returns the squared norm of the point.
func (p Point[T]) Norm2() T {
	return p.Dot(p)
}

// Norm returns the norm of the point.
func (p Point[T]) Norm() float64 {
	return math.Sqrt(float64(p.Norm2()))
}

// Cross returns the cross product of the two points.
func (p Point[T]) Cross(q Point[T]) T {
	return p.X*q.Y - p.Y*q.X
}

// Cross2 returns the cross product of the two points, with p as the origin.
func (p Point[T]) Cross2(q, r Point[T]) T {
	return q.Sub(p).Cross(r.Sub(p))
}

// Rotate90 returns the point rotated 90 degrees clockwise.
func (p Point[T]) Rotate90() Point[T] {
	return Point[T]{-p.Y, p.X}
}

// Rotate180 returns the point rotated 180 degrees clockwise.
func (p Point[T]) Rotate180() Point[T] {
	return Point[T]{-p.X, -p.Y}
}

// Rotate270 returns the point rotated 270 degrees clockwise.
func (p Point[T]) Rotate270() Point[T] {
	return Point[T]{p.Y, -p.X}
}

// In returns whether the point is in the given rectangle.
func (p Point[T]) In(r image.Rectangle) bool {
	return true &&
		T(r.Min.X) <= p.X && p.X < T(r.Max.X) &&
		T(r.Min.Y) <= p.Y && p.Y < T(r.Max.Y)
}

// Within returns whether the point is within the given rectangle that is
// defined by its minimum and maximum points. The minimum point is inclusive,
// and the maximum point is exclusive.
func (p Point[T]) Within(min, max Point[T]) bool {
	return true &&
		min.X <= p.X && p.X < max.X &&
		min.Y <= p.Y && p.Y < max.Y
}

// WithinInclusive returns whether the point is within the given rectangle that
// is defined by its minimum and maximum points. Both the minimum and maximum
// points are inclusive. This function is useful if T is a floating point type.
func (p Point[T]) WithinInclusive(min, max Point[T]) bool {
	return true &&
		min.X <= p.X && p.X <= max.X &&
		min.Y <= p.Y && p.Y <= max.Y
}

// Line describes a line segment in 2D space.
type Line[T ScalarType] struct {
	Start, End Point[T]
}

// String returns a string representation of the line.
func (l Line[T]) String() string {
	return fmt.Sprintf("%v-%v", l.Start, l.End)
}

// Intersection describes the kind of intersection between two lines.
type Intersection uint8

const (
	NoIntersection Intersection = iota
	// Intersect implies a single intersection point.
	Intersect
	// Collinear implies an infinite number of intersection points.
	Collinear
)

// Intersection returns the intersection point of the two lines, if any.
func (l Line[T]) Intersection(m Line[T]) (Point[T], Intersection) {
	// https://github.com/kth-competitive-programming/kactl/blob/main/content/geometry/lineIntersection.h
	v1 := l.End.Sub(l.Start)
	v2 := m.End.Sub(m.Start)

	d := v1.Cross(v2)
	if d == 0 {
		if l.Start.Cross2(l.End, m.Start) == 0 {
			return Point[T]{}, Collinear
		}
		return Point[T]{}, NoIntersection
	}

	p := m.Start.Cross2(l.End, m.End)
	q := m.Start.Cross2(m.End, l.Start)
	i := l.Start.Mul(p).Add(l.End.Mul(q)).Div(d)
	return i, Intersect
}

// RayIntersection returns the intersection point of the two rays, if any.
// It is similar to Intersection, but the intersection point is required to be
// in the ray defined by the line.
func (l Line[T]) RayIntersection(m Line[T]) (Point[T], Intersection) {
	v1 := l.End.Sub(l.Start)
	v2 := m.End.Sub(m.Start)

	d := v1.Cross(v2)
	if d == 0 {
		if l.Start.Cross2(l.End, m.Start) == 0 {
			return Point[T]{}, Collinear
		}
		return Point[T]{}, NoIntersection
	}

	t2 := v1.Cross(l.Start.Sub(m.Start)) / d
	t1 := (m.Start.X - l.Start.X + v2.X*t2) / v1.X
	if t1 < 0 || t2 < 0 {
		// Intersection point is outside the line segments
		return Point[T]{}, NoIntersection
	}

	p := m.Start.Cross2(l.End, m.End)
	q := m.Start.Cross2(m.End, l.Start)
	i := l.Start.Mul(p).Add(l.End.Mul(q)).Div(d)
	return i, Intersect
}
