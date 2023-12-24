package aocutil

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/mat"
)

// Point3D describes a point in 3D space.
type Point3D[T ScalarType] struct {
	X, Y, Z T
}

// Pt3 returns a point in 3D space with the given coordinates.
func Pt3[T ScalarType](x, y, z T) Point3D[T] {
	return Point3D[T]{x, y, z}
}

// String returns a string representation of the point.
func (p Point3D[T]) String() string {
	return fmt.Sprintf("(%v,%v,%v)", p.X, p.Y, p.Z)
}

// RemoveZ returns a 2D point with the Z coordinate removed.
func (p Point3D[T]) RemoveZ() Point[T] {
	return Point[T]{p.X, p.Y}
}

// Add returns the sum of the two points.
func (p Point3D[T]) Add(q Point3D[T]) Point3D[T] {
	return Point3D[T]{p.X + q.X, p.Y + q.Y, p.Z + q.Z}
}

// Sub returns the difference of the two points.
func (p Point3D[T]) Sub(q Point3D[T]) Point3D[T] {
	return Point3D[T]{p.X - q.X, p.Y - q.Y, p.Z - q.Z}
}

// Mul returns the product of the point and the scalar.
func (p Point3D[T]) Mul(s T) Point3D[T] {
	return Point3D[T]{p.X * s, p.Y * s, p.Z * s}
}

// Div returns the quotient of the point and the scalar.
func (p Point3D[T]) Div(s T) Point3D[T] {
	return Point3D[T]{p.X / s, p.Y / s, p.Z / s}
}

// Manhattan returns the Manhattan distance between the two points.
func (p Point3D[T]) Manhattan(q Point3D[T]) T {
	return Abs(p.X-q.X) + Abs(p.Y-q.Y) + Abs(p.Z-q.Z)
}

// Dot returns the dot product of the two points.
func (p Point3D[T]) Dot(q Point3D[T]) T {
	return p.X*q.X + p.Y*q.Y + p.Z*q.Z
}

// Cross returns the cross product of the two points.
func (p Point3D[T]) Cross(q Point3D[T]) Point3D[T] {
	return Point3D[T]{
		p.Y*q.Z - p.Z*q.Y,
		p.Z*q.X - p.X*q.Z,
		p.X*q.Y - p.Y*q.X,
	}
}

// Norm2 returns the squared norm of the point.
func (p Point3D[T]) Norm2() T {
	return p.Dot(p)
}

// Norm returns the norm of the point.
func (p Point3D[T]) Norm() float64 {
	return math.Sqrt(float64(p.Norm2()))
}

// Unit returns the unit vector of the point.
func (p Point3D[T]) Unit() Point3D[T] {
	norm := p.Norm()
	return Point3D[T]{
		T(float64(p.X) / norm),
		T(float64(p.Y) / norm),
		T(float64(p.Z) / norm),
	}
}

// Line3D describes a line in 3D space.
type Line3D[T ScalarType] struct {
	Start Point3D[T]
	End   Point3D[T]
}

// String returns a string representation of the line.
func (l Line3D[T]) String() string {
	return fmt.Sprintf("%s->%s", l.Start, l.End)
}

// RemoveZ returns a 2D line with the Z coordinate removed.
func (l Line3D[T]) RemoveZ() Line[T] {
	return Line[T]{l.Start.RemoveZ(), l.End.RemoveZ()}
}

// IsCoplanarWith returns true if the two lines are coplanar.
func (l Line3D[T]) IsCoplanarWith(m Line3D[T]) bool {
	v1 := l.End.Sub(l.Start)
	v2 := m.End.Sub(m.Start)
	det := mat.Det(mat.NewDense(3, 3, []float64{
		float64(m.Start.X - l.Start.X), float64(m.Start.Y - l.Start.Y), float64(m.Start.Z - l.Start.Z),
		float64(v1.X), float64(v1.Y), float64(v1.Z),
		float64(v2.X), float64(v2.Y), float64(v2.Z),
	}))
	return det == 0
}

// Cuboid3D describes a cuboid in 3D space.
type Cuboid3D[T ScalarType] struct {
	Min Point3D[T] // inclusive
	Max Point3D[T] // exclusive
}

// String returns a string representation of the rectangle.
func (r Cuboid3D[T]) String() string {
	return fmt.Sprintf("%s~%s", r.Min, r.Max)
}

// Canon returns a rectangle with the minimum point in each dimension first.
func (r Cuboid3D[T]) Canon() Cuboid3D[T] {
	if r.Min.X > r.Max.X {
		r.Min.X, r.Max.X = r.Max.X, r.Min.X
	}
	if r.Min.Y > r.Max.Y {
		r.Min.Y, r.Max.Y = r.Max.Y, r.Min.Y
	}
	if r.Min.Z > r.Max.Z {
		r.Min.Z, r.Max.Z = r.Max.Z, r.Min.Z
	}
	return r
}

// ContainsPt returns true if the point is contained in the rectangle.
func (r Cuboid3D[T]) ContainsPt(p Point3D[T]) bool {
	return true &&
		r.Min.X <= p.X && p.X < r.Max.X &&
		r.Min.Y <= p.Y && p.Y < r.Max.Y &&
		r.Min.Z <= p.Z && p.Z < r.Max.Z
}

// Contains returns true if the other rectangle is contained in the rectangle.
func (r Cuboid3D[T]) Contains(other Cuboid3D[T]) bool {
	return true &&
		r.Min.X <= other.Min.X && other.Max.X <= r.Max.X &&
		r.Min.Y <= other.Min.Y && other.Max.Y <= r.Max.Y &&
		r.Min.Z <= other.Min.Z && other.Max.Z <= r.Max.Z
}

// Overlaps returns true if the other rectangle overlaps the rectangle.
func (r Cuboid3D[T]) Overlaps(other Cuboid3D[T]) bool {
	return true &&
		other.Min.X < r.Max.X && r.Min.X < other.Max.X &&
		other.Min.Y < r.Max.Y && r.Min.Y < other.Max.Y &&
		other.Min.Z < r.Max.Z && r.Min.Z < other.Max.Z
}

// Intersect returns the intersection of the two rectangles.
func (r Cuboid3D[T]) Intersect(other Cuboid3D[T]) Cuboid3D[T] {
	return Cuboid3D[T]{
		Min: Point3D[T]{
			max(r.Min.X, other.Min.X),
			max(r.Min.Y, other.Min.Y),
			max(r.Min.Z, other.Min.Z),
		},
		Max: Point3D[T]{
			min(r.Max.X, other.Max.X),
			min(r.Max.Y, other.Max.Y),
			min(r.Max.Z, other.Max.Z),
		},
	}
}

// Union returns the union of the two rectangles.
func (r Cuboid3D[T]) Union(other Cuboid3D[T]) Cuboid3D[T] {
	return Cuboid3D[T]{
		Min: Point3D[T]{
			min(r.Min.X, other.Min.X),
			min(r.Min.Y, other.Min.Y),
			min(r.Min.Z, other.Min.Z),
		},
		Max: Point3D[T]{
			max(r.Max.X, other.Max.X),
			max(r.Max.Y, other.Max.Y),
			max(r.Max.Z, other.Max.Z),
		},
	}
}

// Size returns the size of the rectangle.
func (r Cuboid3D[T]) Size() Point3D[T] {
	return Point3D[T]{
		r.Max.X - r.Min.X,
		r.Max.Y - r.Min.Y,
		r.Max.Z - r.Min.Z,
	}
}

// Volume returns the volume of the rectangle.
func (r Cuboid3D[T]) Volume() T {
	s := r.Size()
	return s.X * s.Y * s.Z
}

// SurfaceArea returns the surface area of the rectangle.
func (r Cuboid3D[T]) SurfaceArea() T {
	s := r.Size()
	return 2 * (s.X*s.Y + s.X*s.Z + s.Y*s.Z)
}

// Translate returns a rectangle translated by the given point.
func (r Cuboid3D[T]) Translate(p Point3D[T]) Cuboid3D[T] {
	return Cuboid3D[T]{
		Min: r.Min.Add(p),
		Max: r.Max.Add(p),
	}
}
