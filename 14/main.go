package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"

	"github.com/diamondburned/aoc-2022/aocutil"
)

// var stderr = os.Stderr

var stderr = io.Discard

func main() {
	log.SetOutput(stderr)

	var drawingLines []Lines

	input := aocutil.InputString()
	lines := aocutil.SplitLines(input)
	for _, line := range lines {
		drawingLines = append(drawingLines, ParseLines(line))
	}

	m := DrawMap(drawingLines)
	part1(aocutil.Clone(m))
	// part2(aocutil.Clone(m))
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
		log.Println("Dropped", i, "sand")
	}
	fmt.Println("part 1:", i)
	sim.PrintMap(os.Stdout)
}

func part2(m Map) {
	sim := NewSimulator(m, SimulatorOpts{
		MapOpts:      mapOpts,
		Infinite:     true,
		FloorOffsetY: 2,
	})

	var i int
	for sim.DropSandAt(dropSandAt) != (Pt{-1, -1}) {
		i++
		log.Println("Dropped", i, "sand")
	}
	fmt.Println("part 2:", i)
	sim.PrintMap(os.Stdout)
}

// Lines is a draw instruction for multiple lines.
type Lines []Pt

func ParseLines(str string) Lines {
	var lines Lines

	parts := strings.Split(str, " -> ")
	for _, part := range parts {
		xy := aocutil.SplitN(part, ",", 2)
		x := aocutil.Atoi[int](xy[0])
		y := aocutil.Atoi[int](xy[1])
		lines = append(lines, Pt{x, y})
	}

	return lines
}

type Pt struct{ X, Y int }

func (pt Pt) Add(other Pt) Pt {
	return Pt{pt.X + other.X, pt.Y + other.Y}
}

func (pt Pt) Addn(x, y int) Pt {
	return Pt{pt.X + x, pt.Y + y}
}

// Map is a 2D map of blocks.
type Map [][]Block

// DrawMap draws the map using the given draw instructions, which is a list of
// lines.
func DrawMap(lines []Lines) Map {
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

type Block = byte

const (
	Void        Block = ' '
	Air         Block = '.'
	Rock        Block = '#'
	Sand        Block = '*'
	RestingSand Block = 'o'
)

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
	log.Println("Growing map horizontally by", delta)

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
	log.Println("Growing map vertically by", delta)

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

	floorY int
}

type SimulatorOpts struct {
	// FloorOffsetY is the offset from the bottom of the map to the floor.
	// If this is 0, then there is no floor.
	FloorOffsetY int
	// Infinite is whether the map is infinite. If this is true, then the map
	// will grow as needed.
	Infinite bool
	// MapOpts are the options to use when printing the map.
	MapOpts MapOpts
}

type MapOpts struct {
	Overlays map[Pt]byte
	Start    Pt
	End      Pt
}

// NewSimulator creates a new simulator for the given map.
func NewSimulator(m Map, o SimulatorOpts) *Simulator {
	floorY := 0
	if o.FloorOffsetY >= 0 {
		floorY = m.Height() - 1 + o.FloorOffsetY
		// Pre-grow our map.
		m.Grow(Pt{0, o.FloorOffsetY})
		// Fill the floor with rock.
		for x := range m[floorY] {
			m[floorY][x] = Rock
		}
	}

	return &Simulator{
		m: m,
		o: o,

		floorY: floorY,
	}
}

func (s *Simulator) DropSandAt(pt Pt) (placed Pt) {
	// Check that the source is not blocked and is just air.
	if s.m.At(pt) != Air {
		return Pt{-1, -1}
	}

	defer func() {
		log.Println("dropped sand:")
		s.PrintMap(log.Writer())
	}()

	fall := fallingState{s.m, pt, Pt{-1, -1}}

fallLoop:
	for {
		switch s.m.At(pt.Addn(0, +1)) {
		case Void:
			log.Println("falling off the map at", pt)
			break fallLoop
		case Air:
			// We're still falling. Drop down.
			pt = fall.set(pt.Addn(0, +1), Sand)
			continue
		}

		// We're not falling anymore. We're either resting on a block, or we're
		// resting on a wall.

		// See if we can rest on the left.
		if s.offsetIsAir(pt, Pt{-1, +1}) {
			// We can rest on the left. Move left.
			pt = fall.set(pt.Addn(-1, +1), Sand)
			continue
		}

		// See if we can rest on the right.
		if s.offsetIsAir(pt, Pt{+1, +1}) {
			// We can rest on the right. Move right.
			pt = fall.set(pt.Addn(+1, +1), Sand)
			continue
		}

		// We can't move left or right. We'll just sit here.
		return fall.set(pt, RestingSand)
	}

	// Check if sand is hitting the floor if we even have one.
	log.Println("seeing floor at", s.floorY)
	if s.floorY > 0 && pt.Y >= s.floorY {
		// If we have a floor, then we'll just sit here.
		return fall.set(pt, RestingSand)
	}

	// If we don't have a floor, then we'll just delete the sand.
	fall.delete()
	return Pt{-1, -1}
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
				if pt.Y != s.floorY {
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
	// Fill the new floor if we have to.
	if s.floorY > 0 {
		for x := len(s.m[s.floorY]) - delta; x < len(s.m[s.floorY]); x++ {
			s.m[s.floorY][x] = Rock
		}
	}
}

type fallingState struct {
	m     Map
	start Pt
	curr  Pt
}

func (s *fallingState) set(newPt Pt, newBlock Block) Pt {
	s.delete()
	s.m[newPt.Y][newPt.X] = newBlock
	s.curr = newPt
	return newPt
}

func (s *fallingState) delete() {
	if s.curr != (Pt{-1, -1}) {
		s.m[s.curr.Y][s.curr.X] = Air
	}
}

func (s *Simulator) PrintMap(w io.Writer) {
	o := s.o.MapOpts
	if w == io.Discard {
		return
	}

	if o.End == (Pt{}) {
		o.End = Pt{math.MaxInt, math.MaxInt}
	}

	for y := o.Start.Y; y <= aocutil.Min2(len(s.m)-1, o.End.Y); y++ {
		for x := o.Start.X; x <= aocutil.Min2(len(s.m[y])-1, o.End.X); x++ {
			switch {
			case s.m[y][x] != Air:
				w.Write([]byte{byte(s.m[y][x])})
			case o.Overlays[Pt{x, y}] != 0:
				w.Write([]byte{o.Overlays[Pt{x, y}]})
			case y == s.floorY:
				w.Write([]byte{'#'})
			default:
				w.Write([]byte{Air})
			}
		}
		w.Write([]byte{'\n'})
	}
}
