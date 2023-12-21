package main

import (
	"image"

	"libdb.so/aoc-2023/aocutil"
)

func main() {
	aocutil.Run(part1, part2)
}

const (
	RoundedRock = 'O'
	CubedRock   = '#'
	EmptySpace  = '.'
)

func parseInput(input string) aocutil.Map2D {
	return aocutil.NewMap2D(input)
}

type Direction = image.Point

var (
	North = image.Pt(0, -1)
	South = image.Pt(0, +1)
	East  = image.Pt(+1, 0)
	West  = image.Pt(-1, 0)
)

// slideRock returns the destination of the rock so that it slides in the
// direction d. If the rock cannot slide, then the destination is the same as
// the source.
func slideRock(m aocutil.Map2D, pt image.Point, d Direction) image.Point {
	dst := pt.Add(d)
	for dst.In(m.Bounds) && m.At(dst) == EmptySpace {
		dst = dst.Add(d)
	}
	return dst.Sub(d)
}

func tiltMap(m aocutil.Map2D, direction Direction) {
	var d image.Point
	var r image.Rectangle

	switch direction {
	case North, West:
		// Scan from the top-left corner.
		d = image.Pt(1, 1)
		r = m.Bounds
	case South, East:
		// Scan from the bottom-right corner.
		d = image.Pt(-1, -1)
		r = image.Rectangle{
			Min: m.Bounds.Max.Sub(image.Pt(1, 1)),
			Max: m.Bounds.Min.Sub(image.Pt(1, 1)),
		}
	}

	for pt := range aocutil.PointsIterateDelta(r, d) {
		at := m.At(pt)
		if at == RoundedRock {
			dst := slideRock(m, pt, direction)
			m.Set(pt, EmptySpace)
			m.Set(dst, RoundedRock)
		}
	}
}

func rockLoad(m aocutil.Map2D, pt image.Point) int {
	return m.Bounds.Dy() - pt.Y
}

func calculateTotalLoad(m aocutil.Map2D) int {
	var total int
	for p, at := range m.All() {
		if at == RoundedRock {
			total += rockLoad(m, p)
		}
	}
	return total
}

func part1(input string) int {
	m := parseInput(input)
	tiltMap(m, North)
	return calculateTotalLoad(m)
}

func part2(input string) int {
	m := parseInput(input)

	type cachedMap struct {
		// Map is the output of the tilts.
		Map aocutil.Map2D
		// Seen is a list of cycles where this map was seen.
		// A maximum of 2 cycles are stored.
		Seen []int
	}

	// cache stores the input aocutil.Map2D and the outcome of all 4 tilts.
	cache := make(map[string]cachedMap)

	const repeat = 1_000_000_000
	for i := 0; i < repeat; i++ {
		tiltMap(m, North)
		tiltMap(m, West)
		tiltMap(m, South)
		tiltMap(m, East)

		mstr := m.String()

		cached, ok := cache[mstr]
		if !ok {
			cached = cachedMap{Map: m.Clone()}
		}

		if len(cached.Seen) < 2 {
			cached.Seen = append(cached.Seen, i)
		} else {
			period := cached.Seen[1] - cached.Seen[0]
			i += (repeat - i) / period * period
		}

		cache[mstr] = cached
	}

	return calculateTotalLoad(m)
}
