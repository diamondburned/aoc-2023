package main

import (
	"fmt"
	"io"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	var movs []Movement

	input := aocutil.InputString()
	lines := aocutil.SplitLines(input)

	for _, line := range lines {
		var direction string
		var amount int
		_, err := fmt.Sscanf(line, "%s %d", &direction, &amount)
		aocutil.E1(err)

		movs = append(movs, Movement{
			Delta:  directions[direction],
			Amount: amount,
		})
	}

	// {
	// 	rope := NewRope(0)
	// 	tailPos := make(map[Pt]struct{})
	// 	for _, mov := range movs {
	// 		for i := 0; i < mov.Amount; i++ {
	// 			rope.MoveHead(mov.Delta)
	// 			tailPos[rope[len(rope)-1]] = struct{}{}
	// 		}
	// 	}

	// 	fmt.Println(len(tailPos))
	// }

	{
		rope := NewRope(10)
		tailPos := make(map[Pt]struct{})
		for _, mov := range movs {
			for i := 0; i < mov.Amount; i++ {
				rope.MoveHead(mov.Delta)
				tailPos[rope[len(rope)-1]] = struct{}{}
			}
		}

		fmt.Println(len(tailPos))
	}
}

type Movement struct {
	Delta  Pt
	Amount int
}

type Pt struct{ X, Y int }

func (pt *Pt) Add(other Pt) {
	pt.X += other.X
	pt.Y += other.Y
}

var directions = map[string]Pt{
	"R": {1, 0},
	"L": {-1, 0},
	"U": {0, 1},
	"D": {0, -1},
}

type Rope []Pt

func NewRope(knots int) Rope {
	if knots < 0 {
		knots = 1 // need head
	}
	return make([]Pt, knots)
}

func (r *Rope) PrintMap(w io.Writer, minX, minY, maxX, maxY int) {
	for _, pt := range *r {
		if pt.X > maxX {
			maxX = pt.X
		}
		if pt.X < minX {
			minX = pt.X
		}
		if pt.Y > maxY {
			maxY = pt.Y
		}
		if pt.Y < minY {
			minY = pt.Y
		}
	}

	for y := maxY; y >= minY; y-- {
		for x := minX; x <= maxX; x++ {
			for i, pt := range *r {
				if pt.X == x && pt.Y == y {
					switch i {
					case 0:
						io.WriteString(w, "H")
					default:
						fmt.Fprintf(w, "%d", i)
					}
					goto next
				}
			}
			io.WriteString(w, ".")
		next:
		}
		w.Write([]byte("\n"))
	}
}

func (r *Rope) MoveHead(delta Pt) {
	prev := &(*r)[0]
	prev.Add(delta)

	for i := 1; i < len(*r); i++ {
		knot := &(*r)[i]
		if !isTouching(*prev, *knot) {
			delta := moveDelta(*prev, *knot)
			knot.Add(delta)
		}
		prev = knot
	}
}

func moveDelta(h, t Pt) Pt {
	var delta Pt
	if h.X > t.X {
		delta.X = 1
	}
	if h.X < t.X {
		delta.X = -1
	}
	if h.Y > t.Y {
		delta.Y = 1
	}
	if h.Y < t.Y {
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
