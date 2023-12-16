package main

import (
	"fmt"
	"image"
	"log"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	aocutil.Run(part1, part2)
}

const (
	EmptySpace      = '.'
	Up90Mirror      = '/'
	Down90Mirror    = '\\'
	SplitVertical   = '|'
	SplitHorizontal = '-'
)

type MirrorMap struct {
	aocutil.Map2D
}

func parseInput(input string) MirrorMap {
	return MirrorMap{aocutil.NewMap2D(input)}
}

type LightBeam struct {
	Position  image.Point
	Direction image.Point
}

func (b LightBeam) String() string {
	return fmt.Sprintf("%v+%v", b.Position, b.Direction)
}

// NextBeam keeps the beam going one step forward. If the beam hits a mirror, it
// will change direction. If the beam hits a
func (b LightBeam) NextBeam(m MirrorMap) []LightBeam {
	nextPos := b.Position.Add(b.Direction)
	next := m.At(nextPos)
	if next == 0 {
		return nil
	}

	switch next {
	case Up90Mirror:
		switch b.Direction {
		case aocutil.VecUp:
			return []LightBeam{{nextPos, aocutil.VecRight}}
		case aocutil.VecDown:
			return []LightBeam{{nextPos, aocutil.VecLeft}}
		case aocutil.VecLeft:
			return []LightBeam{{nextPos, aocutil.VecDown}}
		case aocutil.VecRight:
			return []LightBeam{{nextPos, aocutil.VecUp}}
		}
	case Down90Mirror:
		switch b.Direction {
		case aocutil.VecUp:
			return []LightBeam{{nextPos, aocutil.VecLeft}}
		case aocutil.VecDown:
			return []LightBeam{{nextPos, aocutil.VecRight}}
		case aocutil.VecLeft:
			return []LightBeam{{nextPos, aocutil.VecUp}}
		case aocutil.VecRight:
			return []LightBeam{{nextPos, aocutil.VecDown}}
		}
	case SplitVertical:
		switch b.Direction {
		case aocutil.VecUp, aocutil.VecDown:
			// passes right through
			return []LightBeam{{nextPos, b.Direction}}
		case aocutil.VecLeft, aocutil.VecRight:
			// splits into two
			return []LightBeam{{nextPos, aocutil.VecUp}, {nextPos, aocutil.VecDown}}
		}
	case SplitHorizontal:
		switch b.Direction {
		case aocutil.VecUp, aocutil.VecDown:
			// splits into two
			return []LightBeam{{nextPos, aocutil.VecLeft}, {nextPos, aocutil.VecRight}}
		case aocutil.VecLeft, aocutil.VecRight:
			// passes right through
			return []LightBeam{{nextPos, b.Direction}}
		}
	case EmptySpace:
		return []LightBeam{{nextPos, b.Direction}}
	}
	log.Panicf("invalid beam at %v shooting %v", b.Position, b.Direction)
	return nil
}

func countEnergizedTiles(m MirrorMap, startingBeam LightBeam) int {
	traversed := aocutil.NewSet[LightBeam](0)
	energized := aocutil.NewSet[image.Point](0)
	aocutil.All(aocutil.BFS(startingBeam, func(b LightBeam) []LightBeam {
		if traversed.Has(b) {
			return nil
		}
		traversed.Add(b)
		if b.Position.In(m.Bounds) {
			energized.Add(b.Position)
		}
		return b.NextBeam(m)
	}))
	return len(energized)
}

func part1(input string) int {
	m := parseInput(input)
	return countEnergizedTiles(m, LightBeam{
		image.Point{-1, 0},
		aocutil.VecRight,
	})
}

func part2(input string) int {
	m := parseInput(input)

	var maxEnergized int
	test := func(start LightBeam) {
		// log.Printf("checking %v", start)
		energized := countEnergizedTiles(m, start)
		maxEnergized = max(maxEnergized, energized)
		// log.Printf("beam %v has %d energized tiles", start, energized)
	}

	for y := 0; y < m.Bounds.Max.Y; y++ {
		test(LightBeam{image.Pt(-1, y), aocutil.VecRight})
		test(LightBeam{image.Pt(m.Bounds.Max.X, y), aocutil.VecLeft})
	}
	for x := 0; x < m.Bounds.Max.X; x++ {
		test(LightBeam{image.Pt(x, -1), aocutil.VecDown})
		test(LightBeam{image.Pt(x, m.Bounds.Max.Y), aocutil.VecUp})
	}

	return maxEnergized
}
