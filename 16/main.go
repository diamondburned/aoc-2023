package main

import (
	"fmt"
	"image"

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

type LightBeam struct {
	Position  image.Point
	Direction image.Point
}

func (b LightBeam) String() string {
	return fmt.Sprintf("%v+%v", b.Position, b.Direction)
}

// NextBeam keeps the beam going one step forward. If the beam hits a mirror, it
// will change direction. If the beam hits a
func (b LightBeam) NextBeam(m aocutil.Map2D) []LightBeam {
	nextPos := b.Position.Add(b.Direction)
	switch m.At(nextPos) {
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
	return nil
}

func countEnergizedTiles(m aocutil.Map2D, startingBeam LightBeam) int {
	energized := aocutil.NewSet[image.Point](0)
	aocutil.All(aocutil.AcyclicBFS(startingBeam, func(b LightBeam) []LightBeam {
		if b.Position.In(m.Bounds) {
			energized.Add(b.Position)
		}
		return b.NextBeam(m)
	}))
	return len(energized)
}

func parseInput(input string) aocutil.Map2D {
	return aocutil.NewMap2D(input)
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
		energized := countEnergizedTiles(m, start)
		maxEnergized = max(maxEnergized, energized)
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
