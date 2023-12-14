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
	RoundedRock = 'O'
	CubedRock   = '#'
	EmptySpace  = '.'
)

type Map struct {
	Data   [][]byte
	Bounds image.Rectangle
}

func parseInput(input string) Map {
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
	sb.Grow(len(m.Data) * (len(m.Data[0]) + 1))
	for _, line := range m.Data {
		sb.Write(line)
		sb.WriteByte('\n')
	}
	return sb.String()
}

// Load returns the amount of load caused by ar ock.
func (m Map) Load(pt image.Point) int {
	return m.Bounds.Dy() - pt.Y
}

type Direction = image.Point

var (
	North = image.Pt(0, -1)
	South = image.Pt(0, +1)
	East  = image.Pt(+1, 0)
	West  = image.Pt(-1, 0)
)

// canSlide returns true and the point where the rock will slide to if the
// rock at the given point can slide in the given direction.
func canSlide(m Map, pt image.Point, d Direction) (image.Point, bool) {
	dst := pt.Add(d)
	for {
		if !dst.In(m.Bounds) || m.At(dst) != EmptySpace {
			dst = dst.Sub(d)
			break
		}
		dst = dst.Add(d)
	}

	aocutil.Assertf(
		dst == pt || m.At(dst) == EmptySpace,
		"rock at %v cannot slide over non-empty space at %v", pt, dst,
	)

	return dst, dst != pt
}

func tiltMap(m Map, direction Direction) {
	var delta image.Point
	var start, end image.Point // [start, end)

	switch direction {
	case North, West:
		// Scan from the top-left corner.
		delta = image.Pt(1, 1)
		start = m.Bounds.Min
		end = m.Bounds.Max
	case South, East:
		// Scan from the bottom-right corner.
		delta = image.Pt(-1, -1)
		start = m.Bounds.Max.Sub(image.Pt(1, 1))
		end = m.Bounds.Min.Sub(image.Pt(1, 1))
	}

	for y := start.Y; y != end.Y; y += delta.Y {
		for x := start.X; x != end.X; x += delta.X {
			pt := image.Pt(x, y)
			at := m.At(pt)
			if at != RoundedRock {
				continue
			}

			dst, canSlide := canSlide(m, pt, direction)
			if !canSlide {
				continue
			}

			m.Set(image.Pt(x, y), EmptySpace)
			m.Set(dst, RoundedRock)
		}
	}
}

func calculateTotalLoad(m Map) int {
	var total int
	for y := m.Bounds.Min.Y; y < m.Bounds.Max.Y; y++ {
		for x := m.Bounds.Min.X; x < m.Bounds.Max.X; x++ {
			pt := image.Pt(x, y)
			if m.At(pt) != RoundedRock {
				continue
			}
			total += m.Load(pt)
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
	reps := make(map[string]int)

	type cachedMap struct {
		Map
		// Seen is a list of cycles where this map was seen.
		// A maximum of 2 cycles are stored.
		Seen []int
	}

	// cache stores the input Map and the outcome of all 4 tilts.
	cache := make(map[string]cachedMap)

	const repeat = 1_000_000_000
	for i := 0; i < repeat; i++ {
		log.Printf("cycle %d", i)
		mstr := m.String()

		tiltMap(m, North)
		tiltMap(m, West)
		tiltMap(m, South)
		tiltMap(m, East)
		log.Print("cycle ", i, ":\n", m)

		cached, ok := cache[mstr]
		if !ok {
			cached = cachedMap{Map: m.Clone(), Seen: []int{i}}
			cache[mstr] = cached
			continue
		}

		// We've seen this map before. Try to track how many cycles it
		// takes to repeat.
		if len(cached.Seen) < 5 {
			cached.Seen = append(cached.Seen, i)
			cache[mstr] = cached
			continue
		}

		m = cached.Map
		log.Print("cycle ", i, " (cached):\n", m)

		// We've seen this sequence twice. Jump ahead to the next time it
		// repeats.
		seen := cached.Seen
		period := seen[1] - seen[0]
		jump := (repeat - i) / period
		log.Printf("%v: jumping ahead %d cycles", seen, jump*period)
		i += jump * period
	}

	var cycles []Map
	for mstr, rep := range reps {
		if rep == 100 {
			cycles = append(cycles, parseInput(mstr))
			log.Print("map:\n", mstr)
		}
	}

	return calculateTotalLoad(m)
}
