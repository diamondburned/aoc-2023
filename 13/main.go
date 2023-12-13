package main

import (
	"image"
	"log"
	"slices"
	"strings"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	aocutil.Run(part1, part2)
}

const (
	Ash  = '.'
	Rock = '#'
)

func flipBlock(block byte) byte {
	switch block {
	case Ash:
		return Rock
	case Rock:
		return Ash
	}
	panic("unreachable")
}

type Map struct {
	Data   [][]byte
	Bounds image.Rectangle
}

func parseInput(input string) []Map {
	blocks := strings.Split(input, "\n\n")
	blocks = aocutil.FilterEmptyStrings(blocks)
	maps := make([]Map, len(blocks))
	for i, block := range blocks {
		maps[i] = parseMap(block)
	}
	return maps
}

func parseMap(input string) Map {
	var m Map
	lines := aocutil.SplitLines(input)
	lines = aocutil.FilterEmptyStrings(lines)
	m.Data = make([][]byte, len(lines))
	for i, line := range lines {
		m.Data[i] = []byte(line)
	}
	m.Bounds = image.Rect(0, 0, len(m.Data[0]), len(m.Data))
	return m
}

// At returns the byte at the given point. If the point is out of bounds, then
// 0 is returned.
func (m Map) At(p image.Point) byte {
	if !p.In(m.Bounds) {
		return 0
	}
	return m.Data[p.Y][p.X]
}

// Set sets the byte at the given point.
func (m Map) Set(p image.Point, b byte) {
	m.Data[p.Y][p.X] = b
}

// Clone makes a copy of the map.
func (m Map) Clone() Map {
	var n Map
	n.Data = make([][]byte, len(m.Data))
	for i, line := range m.Data {
		n.Data[i] = slices.Clone(line)
	}
	n.Bounds = m.Bounds
	return n
}

func (m Map) String() string {
	var sb strings.Builder
	for _, line := range m.Data {
		sb.Write(line)
		sb.WriteByte('\n')
	}
	return sb.String()
}

// findReflections finds all reflections in the map. A reflection is a line
// where the map is mirrored on both sides. If the line is horizontal, then {0,
// Y} is returned. If the line is vertical, then {X, 0} is returned.
func findReflections(m Map) aocutil.Iter[image.Point] {
	return func(yield func(pt image.Point) bool) {
		// Search vertically first. Our strategy is to compare the first line and
		// see if any part of it is mirrored on the other side.
	mirrorX:
		for x := 0; x < m.Bounds.Max.X; x++ {
			// Continue checking for Y > 0.
			for y := 0; y < m.Bounds.Max.Y; y++ {
				if !isMirroredAt(m, image.Pt(x, y), Vertical) {
					continue mirrorX
				}
			}
			if !yield(image.Pt(x, 0)) {
				return
			}
		}

	mirrorY:
		for y := 0; y < m.Bounds.Max.Y; y++ {
			for x := 0; x < m.Bounds.Max.X; x++ {
				if !isMirroredAt(m, image.Pt(x, y), Horizontal) {
					continue mirrorY
				}
			}
			if !yield(image.Pt(0, y)) {
				return
			}
		}
	}
}

type Axis int

const (
	Horizontal Axis = iota
	Vertical
)

func (a Axis) String() string {
	switch a {
	case Horizontal:
		return "Y"
	case Vertical:
		return "X"
	}
	panic("unreachable")
}

func (a Axis) Other() Axis {
	switch a {
	case Horizontal:
		return Vertical
	case Vertical:
		return Horizontal
	}
	panic("unreachable")
}

// isMirroredAt returns true if the line is mirrored at the given point.
// atFunc is a function that returns the byte at the given index. If the index
// is out of bounds, then it should return 0.
func isMirroredAt(m Map, pt image.Point, axis Axis) bool {
	switch pt {
	case
		m.Bounds.Min,
		m.Bounds.Max,
		image.Pt(m.Bounds.Min.X, m.Bounds.Max.Y),
		image.Pt(m.Bounds.Max.X, m.Bounds.Min.Y):

		// We don't check the corners.
		return false
	}

	// log.Printf("checking at %s for axis %s", pt, axis)

	var at int
	switch axis {
	case Horizontal:
		at = pt.Y
	case Vertical:
		at = pt.X
	}

	// Search left and right simultaneously.
	for i1, i2 := at-1, at; true; i1, i2 = i1-1, i2+1 {
		var p1, p2 image.Point
		switch axis {
		case Horizontal:
			p1 = image.Pt(pt.X, i1)
			p2 = image.Pt(pt.X, i2)
		case Vertical:
			p1 = image.Pt(i1, pt.Y)
			p2 = image.Pt(i2, pt.Y)
		}
		b1 := m.At(p1)
		b2 := m.At(p2)
		// log.Printf("  %s='%c', %s='%c'", p1, b1, p2, b2)
		if b1 == 0 || b2 == 0 {
			return true
		}
		if b1 != b2 {
			// We're not mirrored.
			return false
		}
	}
	panic("unreachable")
}

func encodePt(pt image.Point) int {
	return pt.X + pt.Y*100
}

func part1(input string) int {
	maps := parseInput(input)
	var summary int
	for _, m := range maps {
		reflectPt, ok := aocutil.Once(findReflections(m))
		if !ok {
			panic("no reflection found")
		}
		summary += encodePt(reflectPt)
	}
	return summary
}

type UnsmudgedReflection struct {
	// Map is the unsmudged map.
	// The new reflection can be found in this map.
	Map Map
	// Reflection is the reflection point in the new map.
	Reflection image.Point
}

func findNewReflection(m Map) (UnsmudgedReflection, bool) {
	log.Print("\n", m)
	midpoint := image.Pt(m.Bounds.Max.X/2, m.Bounds.Max.Y/2)

	reflectPt, ok := aocutil.Once(findReflections(m))
	if !ok {
		return UnsmudgedReflection{}, false
	}

	log.Println("midpoint:", midpoint)
	log.Println("reflection:", reflectPt)

	for y := m.Bounds.Min.Y; y < m.Bounds.Max.Y; y++ {
		for x := m.Bounds.Min.X; x < m.Bounds.Max.X; x++ {
			pt := image.Pt(x, y)
			old := m.At(pt)

			// Flip the block and find a new reflection.
			m.Set(pt, flipBlock(old))

			for newReflectPt := range findReflections(m) {
				switch newReflectPt {
				case (image.Point{}), reflectPt:
					continue
				}
				return UnsmudgedReflection{
					Map:        m,
					Reflection: newReflectPt,
				}, true
			}

			// Flip it back.
			m.Set(pt, old)
		}
	}

	return UnsmudgedReflection{}, false
}

func part2(input string) int {
	maps := parseInput(input)
	var summary int
	for _, m := range maps {
		m2, ok := findNewReflection(m)
		if !ok {
			panic("no new reflection found")
		}
		reflectPt := m2.Reflection
		log.Printf("new reflection: %s", reflectPt)
		summary += encodePt(reflectPt)
	}
	return summary
}
