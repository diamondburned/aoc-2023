package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	var movs []Movement

	input := aocutil.InputString()
	lines := aocutil.SplitLines(input)

	for _, line := range lines {
		var direction string
		var amount int
		aocutil.Sscanf(line, "%s %d", &direction, &amount)

		movs = append(movs, Movement{
			Delta:  directions[direction],
			Amount: amount,
		})
	}

	playMovements := func(f func(delta Pt)) {
		for _, mov := range movs {
			for i := 0; i < mov.Amount; i++ {
				f(mov.Delta)
			}
		}
	}

	mapMin := Pt{-5, -5}
	mapMax := Pt{+5, +5}

	{
		rope := NewRope(2)
		tailPos := aocutil.NewSet[Pt](0)
		playMovements(func(delta Pt) {
			rope.MoveHead(delta)
			tailPos.Add(rope.Tail())
		})

		fmt.Println("part 1: the tail has been on", len(tailPos), "positions")
		rope.PrintMap(os.Stdout, mapMin, mapMax)
	}

	fmt.Println("=======")

	{
		rope := NewRope(10)
		tailPos := aocutil.NewSet[Pt](0)
		playMovements(func(delta Pt) {
			rope.MoveHead(delta)
			tailPos.Add(rope.Tail())
		})

		fmt.Println("part 2: the tail has been on", len(tailPos), "positions")
		rope.PrintMap(os.Stdout, mapMin, mapMax)
	}
}

var directions = map[string]Pt{
	"R": {1, 0},
	"L": {-1, 0},
	"U": {0, 1},
	"D": {0, -1},
}

// Movement represents a movement of the head of the rope.
type Movement struct {
	Delta  Pt
	Amount int
}

// Pt is a point in 2D space.
type Pt struct{ X, Y int }

// Add adds the given pt to the point.
func (pt *Pt) Add(other Pt) {
	pt.X += other.X
	pt.Y += other.Y
}

// Rope is a rope of knots. Each knot is a point.
type Rope []Pt

// NewRope creates a new rope with the given length. The rope will be created
// with the head at the origin, and the tail at the given length.
func NewRope(knots int) Rope {
	if knots < 2 {
		knots = 2 // need head and tail
	}
	return make([]Pt, knots)
}

// PrintMap prints the map of the rope. The head is marked with an 'H', the tail
// and knots are marked with a 'T' or are numbered.
func (r Rope) PrintMap(w io.Writer, min, max Pt) {
	for _, pt := range r {
		max.X = aocutil.Max2(max.X, pt.X)
		max.Y = aocutil.Max2(max.Y, pt.Y)
		min.X = aocutil.Min2(min.X, pt.X)
		min.Y = aocutil.Min2(min.Y, pt.Y)
	}

	mapb := make([][]byte, max.Y-min.Y+1)
	for i := range mapb {
		mapb[i] = make([]byte, max.X-min.X+1)
		for j := range mapb[i] {
			x := j + min.X
			y := i + min.Y
			if x == 0 && y == 0 {
				mapb[i][j] = 's'
			} else {
				mapb[i][j] = '.'
			}
		}
	}

	for i, pt := range r {
		var char byte
		switch {
		case i == 0:
			char = 'H'
		case len(r) == 2:
			char = 'T'
		default:
			v := strconv.Itoa(i)
			char = v[0]
		}
		mapb[pt.Y-min.Y][pt.X-min.X] = char
	}

	for i := len(mapb) - 1; i >= 0; i-- {
		w.Write(mapb[i])
		w.Write([]byte{'\n'})
	}
}

// MoveHead moves the head of the rope by the given delta. The knots will be
// moved accordingly.
func (r Rope) MoveHead(delta Pt) {
	prev := &r[0]
	prev.Add(delta)

	for i := 1; i < len(r); i++ {
		knot := &r[i]
		if !isTouching(*prev, *knot) {
			delta := moveDelta(*prev, *knot)
			knot.Add(delta)
		}
		prev = knot
	}
}

// Head returns the head of the rope. The head is the first knot.
func (r Rope) Head() Pt {
	return r[0]
}

// Tail returns the tail of the rope. The tail is the last knot.
func (r Rope) Tail() Pt {
	return r[len(r)-1]
}

func moveDelta(h, t Pt) Pt {
	var delta Pt
	switch {
	case h.X > t.X:
		delta.X = 1
	case h.X < t.X:
		delta.X = -1
	}
	switch {
	case h.Y > t.Y:
		delta.Y = 1
	case h.Y < t.Y:
		delta.Y = -1
	}
	return delta
}

func isTouching(h, t Pt) bool {
	return false ||
		h == t ||
		// Top, bottom, left, right.
		h.X+1 == t.X && h.Y == t.Y ||
		h.X-1 == t.X && h.Y == t.Y ||
		h.X == t.X && h.Y+1 == t.Y ||
		h.X == t.X && h.Y-1 == t.Y ||
		// Four corners.
		h.X+1 == t.X && h.Y+1 == t.Y ||
		h.X+1 == t.X && h.Y-1 == t.Y ||
		h.X-1 == t.X && h.Y+1 == t.Y ||
		h.X-1 == t.X && h.Y-1 == t.Y
}
