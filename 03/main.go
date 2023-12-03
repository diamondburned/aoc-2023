package main

import (
	"bytes"
	"image"
	"strconv"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	aocutil.Run(part1, part2)
}

type Schematic [][]byte

func parseSchematic(input string) Schematic {
	lines := aocutil.SplitLines(input)
	schematic := make(Schematic, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			continue
		}
		schematic = append(schematic, []byte(line))
	}
	return schematic
}

const Gear = '*'

func isDigit[T byte | rune](b T) bool {
	return b >= '0' && b <= '9'
}

func isSymbol[T byte | rune](b T) bool {
	return b != 0 && !isDigit(b) && b != '.'
}

// Bounds returns the bounds of the schematic.
func (s Schematic) Bounds() image.Rectangle {
	return image.Rect(0, 0, len(s[0]), len(s))
}

// In returns true if the given position is within the schematic.
func (s Schematic) In(pt image.Point) bool {
	return pt.X >= 0 && pt.Y >= 0 && pt.Y < len(s) && pt.X < len(s[pt.Y])
}

// At returns the byte at the given position, or 0 if out of bounds.
func (s Schematic) At(pt image.Point) byte {
	if !s.In(pt) {
		return 0
	}
	return s[pt.Y][pt.X]
}

// SearchNumber searches for a number at the given position, returning the
// number of digits and the end position. If no number is found, it returns 0
// and the start position.
func (s Schematic) SearchNumber(pt image.Point) (digits int, r image.Rectangle) {
	if !s.In(pt) {
		return 0, image.Rectangle{}
	}
	row := s[pt.Y]
	end := bytes.IndexFunc(row[pt.X:], func(r rune) bool { return !isDigit(r) })
	if end == -1 {
		end = len(row) - pt.X
	}

	n, err := strconv.Atoi(string(row[pt.X : pt.X+end]))
	if err != nil {
		return 0, image.Rectangle{}
	}

	return n, image.Rect(pt.X, pt.Y, pt.X+end, pt.Y+1)
}

// EachAdjacent calls match for each adjacent cell to the given rectangle
// defined by p1 and p2. The match function is called with the byte at the
// position, and if it returns true, the position is returned.
func (s Schematic) EachAdjacent(rect image.Rectangle, match func(pt image.Point, b byte) bool) {
	rect = rect.Canon()
	rect = rect.Intersect(s.Bounds())

	border := rect
	border = image.Rect(border.Min.X-1, border.Min.Y-1, border.Max.X+1, border.Max.Y+1)
	border = border.Intersect(s.Bounds())

	for y := border.Min.Y; y < border.Max.Y; y++ {
		for x := border.Min.X; x < border.Max.X; x++ {
			pt := image.Pt(x, y)
			if !pt.In(rect) && match(pt, s.At(pt)) {
				return
			}
		}
	}
}

func part1(input string) int {
	schematic := parseSchematic(input)

	var sum int
	for y, row := range schematic {
		for x := 0; x < len(row); x++ {
			pt := image.Pt(x, y)

			n, rect := schematic.SearchNumber(pt)
			if rect.Empty() {
				continue
			}

			var nearSymbol bool
			schematic.EachAdjacent(rect, func(pt image.Point, b byte) bool {
				nearSymbol = isSymbol(b)
				return nearSymbol
			})

			if nearSymbol {
				sum += n
			}

			x = rect.Max.X - 1
		}
	}

	return sum
}

func part2(input string) int {
	schematic := parseSchematic(input)

	// gears tracks the gears at each point and its adjacent numbers.
	gears := make(map[image.Point][]int)

	for y, row := range schematic {
		for x := 0; x < len(row); x++ {
			pt := image.Pt(x, y)

			n, rect := schematic.SearchNumber(pt)
			if rect.Empty() {
				continue
			}

			schematic.EachAdjacent(rect, func(pt image.Point, b byte) bool {
				if b == Gear {
					gears[pt] = append(gears[pt], n)
				}
				return false
			})

			x = rect.Max.X - 1
		}
	}

	var sum int
	for _, nums := range gears {
		if len(nums) == 2 {
			sum += aocutil.Mul(nums)
		}
	}

	return sum
}
