package aocutil

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var (
	// VecUp is a vector pointing up.
	VecUp = image.Pt(0, -1)
	// VecDown is a vector pointing down.
	VecDown = image.Pt(0, +1)
	// VecLeft is a vector pointing left.
	VecLeft = image.Pt(-1, 0)
	// VecRight is a vector pointing right.
	VecRight = image.Pt(+1, 0)
)

// CardinalDirections is a list of cardinal directions.
var CardinalDirections = []image.Point{
	VecUp,
	VecDown,
	VecLeft,
	VecRight,
}

// Map2D is a 2D map of bytes.
type Map2D struct {
	Data   [][]byte
	Bounds image.Rectangle
}

// NewMap2D creates a new Map2D from the given input.
func NewMap2D(input string) Map2D {
	lines := SplitLines(input)
	data := make([][]byte, len(lines))
	for i, line := range lines {
		data[i] = []byte(line)
	}
	return NewMap2DFromData(data)
}

func NewEmptyMap2D(bounds image.Rectangle) Map2D {
	data := make([][]byte, bounds.Dy())
	for i := range data {
		data[i] = make([]byte, bounds.Dx())
	}
	m := NewMap2DFromData(data)
	m.Bounds = bounds
	return m
}

// NewMap2DFromData creates a new Map2D from the given data.
func NewMap2DFromData(data [][]byte) Map2D {
	m := Map2D{Data: data}
	m.Bounds = image.Rect(0, 0, len(m.Data[0]), len(m.Data))
	return m
}

func (m Map2D) pt(pt image.Point) image.Point {
	return pt.Sub(m.Bounds.Min)
}

// At returns the byte at the given point. If the point is out of bounds, then
// 0 is returned.
func (m Map2D) At(p image.Point) byte {
	if !p.In(m.Bounds) {
		return 0
	}
	p = m.pt(p)
	return m.Data[p.Y][p.X]
}

// Set sets the byte at the given point.
func (m Map2D) Set(p image.Point, b byte) {
	if !p.In(m.Bounds) {
		return
	}
	p = m.pt(p)
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

// Draw draws the map to a paletted image.
func (m Map2D) Draw(colorMap map[byte]color.RGBA) *image.Paletted {
	palette := make(color.Palette, 1, 1+len(colorMap))
	palette[0] = color.RGBA{0, 0, 0, 255}

	colorIx := make(map[byte]uint8)
	for k, v := range colorMap {
		palette = append(palette, v)
		colorIx[k] = uint8(len(palette) - 1)
	}

	img := image.NewPaletted(m.Bounds, palette)
	for pt, v := range m.All() {
		ix := colorIx[v]
		img.SetColorIndex(pt.X, pt.Y, ix)
	}

	return img
}

// Equal returns true if the maps are equal.
func (m Map2D) Equal(other Map2D) bool {
	if m.Bounds != other.Bounds {
		return false
	}
	for i, line := range m.Data {
		if !bytes.Equal(line, other.Data[i]) {
			return false
		}
	}
	return true
}

// Transpose returns a transposed copy of the map.
func (m Map2D) Transpose() Map2D {
	data := make([][]byte, len(m.Data[0]))
	for i := range data {
		data[i] = make([]byte, len(m.Data))
		for j := range data[i] {
			data[i][j] = m.Data[j][i]
		}
	}
	return NewMap2DFromData(data)
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
				pt := image.Pt(x, y)
				if !yield(pt, m.At(pt)) {
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

// SaveImage saves the given image to a PNG image.
func SaveImage(img image.Image, dst string) {
	if ext := filepath.Ext(dst); ext != ".png" {
		panic(fmt.Errorf("invalid extension %q, only .png is supported", ext))
	}
	f := E2(os.Create(dst))
	defer f.Close()
	E1(png.Encode(f, img))
}

// RectangleContainingPoints returns the smallest rectangle that contains all
// the given points. It takes the union of all the points.
func RectangleContainingPoints(pts []image.Point) image.Rectangle {
	var r image.Rectangle
	for _, pt := range pts {
		r = r.Union(image.Rectangle{Min: pt, Max: pt})
	}
	return r
}
