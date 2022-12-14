package main

import (
	"fmt"
	"image"
	"io"
	"log"
	"math"
	"os"
	"strings"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	log.SetOutput(io.Discard)

	var drawingLines []Line

	input := aocutil.InputString()
	lines := aocutil.SplitLines(input)
	for _, line := range lines {
		drawingLines = append(drawingLines, ParseLines(line))
	}

	m := DrawMap(drawingLines)
	part1(aocutil.Clone(m))
	part2(aocutil.Clone(m))
}

var (
	dropSandAt = Pt{500, 0}
	mapOpts    = MapOpts{
		Overlays: map[Pt]byte{dropSandAt: '+'},
		Start:    Pt{480, 0},
		End:      Pt{},
	}
)

func part1(m Map) {
	sim := NewSimulator(m, SimulatorOpts{
		MapOpts: mapOpts,
	})

	var i int
	for sim.DropSandAt(dropSandAt) != (Pt{-1, -1}) {
		i++
	}
	fmt.Println("part 1:", i)
	sim.PrintMap(os.Stdout)
}

func part2(m Map) {
	sim := NewSimulator(m, SimulatorOpts{
		MapOpts:  mapOpts,
		FloorY:   m.Height() + 1,
		Infinite: true,
	})

	var i int
	for sim.DropSandAt(dropSandAt) != (Pt{-1, -1}) {
		i++
	}
	fmt.Println("part 2:", i)
	sim.PrintMap(os.Stdout)
}

// Pt is a point in 2D space.
type Pt = image.Point

var ZP = Pt{}

func addPt(pt Pt, dx, dy int) Pt { return pt.Add(Pt{dx, dy}) }

// Line is a draw instruction for a continuous line.
type Line []Pt

// ParseLines parses the given string for line instructions.
func ParseLines(str string) Line {
	var lines Line

	parts := strings.Split(str, " -> ")
	for _, part := range parts {
		xy := aocutil.SplitN(part, ",", 2)
		x := aocutil.Atoi[int](xy[0])
		y := aocutil.Atoi[int](xy[1])
		lines = append(lines, Pt{x, y})
	}

	return lines
}

type Block = byte

const (
	Void        Block = ' '
	Air         Block = '.'
	Rock        Block = '#'
	Sand        Block = '*'
	RestingSand Block = 'o'
)

// Map is a 2D map of blocks.
type Map [][]Block

// DrawMap draws the map using the given draw instructions, which is a list of
// lines.
func DrawMap(lines []Line) Map {
	var maxX, maxY int
	for _, line := range lines {
		for _, pt := range line {
			if pt.X > maxX {
				maxX = pt.X
			}
			if pt.Y > maxY {
				maxY = pt.Y
			}
		}
	}

	var m Map
	m.Grow(Pt{maxX + 1, maxY + 1})

	for _, line := range lines {
		prev := line[0]
		for _, curr := range line[1:] {
			if prev.X == curr.X {
				minY := aocutil.Min2(prev.Y, curr.Y)
				maxY := aocutil.Max2(prev.Y, curr.Y)
				for y := minY; y <= maxY; y++ {
					m[y][curr.X] = Rock
				}
			}
			if prev.Y == curr.Y {
				minX := aocutil.Min2(prev.X, curr.X)
				maxX := aocutil.Max2(prev.X, curr.X)
				for x := minX; x <= maxX; x++ {
					m[curr.Y][x] = Rock
				}
			}
			prev = curr
		}
	}

	return m
}

func (m Map) At(pt Pt) Block {
	if pt.Y < 0 || pt.Y >= len(m) || pt.X < 0 || pt.X >= len(m[pt.Y]) {
		return Void
	}
	return m[pt.Y][pt.X]
}

func (m Map) Width() int  { return len(m[0]) }
func (m Map) Height() int { return len(m) }

// GrowX grows the map by the given amount. If delta.X is negative, then the map
// is grown to the left. If delta.X is positive, then the map is grown to the
// right. delta.Y cannot be negative.
func (m *Map) Grow(delta Pt) {
	m.growY(delta.Y)
	m.growX(delta.X)
}

func (m *Map) growX(delta int) {
	amount := delta
	if amount < 0 {
		amount = -amount
	}

	mm := *m
	for y := range mm {
		mm[y] = append(mm[y], make([]Block, amount)...)
		// Are we growing to the left?
		if delta < 0 {
			// Yes. We'll need to shift the map to the right.
			copy(mm[y][amount:], mm[y])
			// Fill the left space with air.
			for x := 0; x < amount; x++ {
				mm[y][x] = Air
			}
		} else {
			// Fill the right space with air.
			for x := len(mm[y]) - amount; x < len(mm[y]); x++ {
				mm[y][x] = Air
			}
		}
	}

	*m = mm
}

func (m *Map) growY(delta int) {
	aocutil.Assert(delta >= 0, "delta must be positive")

	var w int
	if len(*m) > 0 {
		w = len((*m)[0])
	}

	mm := *m
	for i := 0; i < delta; i++ {
		row := make([]Block, w)
		// Fill it with air.
		for x := range row {
			row[x] = Air
		}

		mm = append(mm, row)
	}

	*m = mm
}

type Simulator struct {
	// m is the map to simulate on.
	m Map
	o SimulatorOpts
}

type SimulatorOpts struct {
	// FloorY is the Y-value of the floor. If this is 0, then there is no floor.
	// If FloorY is overbound, then it is useless unless Infinite is true.
	FloorY int
	// Infinite is whether the map is infinite. If this is true, then the map
	// will grow as needed.
	Infinite bool
	// MapOpts are the options to use when printing the map.
	MapOpts MapOpts
}

// MapOpts are the options to use when printing the map.
type MapOpts struct {
	Overlays map[Pt]byte
	Start    Pt
	End      Pt
}

// NewSimulator creates a new simulator for the given map.
func NewSimulator(m Map, o SimulatorOpts) *Simulator {
	s := Simulator{
		m: m,
		o: o,
	}

	if o.Infinite && o.FloorY >= m.Height() {
		s.m.Grow(Pt{0, o.FloorY - m.Height() + 1})
		s.fillFloor(0, m.Width())
	}

	return &s
}

func (s *Simulator) DropSandAt(pos Pt) (placed Pt) {
	// Check that the source is not blocked and is just air.
	if s.m.At(pos) != Air {
		return Pt{-1, -1}
	}

	fall := func(pt Pt, b Block) {
		if pos != (Pt{-1, -1}) {
			s.m[pos.Y][pos.X] = Air
		}
		if pt != (Pt{-1, -1}) {
			s.m[pt.Y][pt.X] = b
		}
		pos = pt
	}

fallLoop:
	for {
		switch s.m.At(addPt(pos, 0, +1)) {
		case Void:
			break fallLoop
		case Air:
			// We're still falling. Drop down.
			fall(addPt(pos, 0, +1), Sand)
			continue
		}

		// We're not falling anymore. We're either resting on a block, or we're
		// resting on a wall.

		// See if we can rest on the left.
		if s.offsetIsAir(pos, Pt{-1, +1}) {
			// We can rest on the left. Move left.
			fall(addPt(pos, -1, +1), Sand)
			continue
		}

		// See if we can rest on the right.
		if s.offsetIsAir(pos, Pt{+1, +1}) {
			// We can rest on the right. Move right.
			fall(addPt(pos, +1, +1), Sand)
			continue
		}

		// We can't move left or right. We'll just sit here.
		fall(pos, RestingSand)
		return pos
	}

	// Check if sand is hitting the floor if we even have one.
	if s.o.FloorY > 0 && pos.Y >= s.o.FloorY {
		// If we have a floor, then we'll just sit here.
		fall(pos, RestingSand)
	} else {
		// If we don't have a floor, then we'll just delete the sand.
		fall(Pt{-1, -1}, Void)
	}

	return pos
}

func (s *Simulator) offsetIsAir(pt Pt, offsets ...Pt) bool {
	for _, offset := range offsets {
		switch pt := pt.Add(offset); s.m.At(pt) {
		case Air:
			continue
		case Void:
			// Are we simulating an infinite map?
			if s.o.Infinite {
				// Yes. Grow the map.
				s.growX(offset.X)
				// If we're not on floorY, then we're in the air, so we can
				// continue.
				if pt.Y != s.o.FloorY {
					continue
				}
			}
			fallthrough
		default:
			return false
		}
	}
	return true
}

func (s *Simulator) growX(delta int) {
	s.m.Grow(Pt{delta, 0})
	if s.o.FloorY > 0 {
		s.fillFloor(len(s.m[0])-delta, len(s.m[0]))
	}
}

func (s *Simulator) fillFloor(x1, x2 int) {
	for x := x1; x < x2 && x < len(s.m[s.o.FloorY]); x++ {
		s.m[s.o.FloorY][x] = Rock
	}
}

func (s *Simulator) PrintMap(w io.Writer) {
	if w == io.Discard {
		return
	}

	o := s.o.MapOpts
	if o.End == ZP {
		o.End = Pt{math.MaxInt, math.MaxInt}
	}

	for y := o.Start.Y; y <= aocutil.Min2(len(s.m)-1, o.End.Y); y++ {
		for x := o.Start.X; x <= aocutil.Min2(len(s.m[y])-1, o.End.X); x++ {
			switch {
			case s.m[y][x] != Air:
				w.Write([]byte{s.m[y][x]})
			case o.Overlays[Pt{x, y}] != 0:
				w.Write([]byte{o.Overlays[Pt{x, y}]})
			case s.o.FloorY != 0 && y == s.o.FloorY:
				w.Write([]byte{'#'})
			default:
				w.Write([]byte{Air})
			}
		}
		w.Write([]byte{'\n'})
	}
}
