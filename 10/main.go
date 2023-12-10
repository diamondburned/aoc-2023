package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	aocutil.Run(part1, part2)
}

type Block = byte

const (
	VerticalPipe   Block = '|' // North-South
	HorizontalPipe       = '-' // East-West
	NEBend         Block = 'L' // 90 degree bend, North-East
	NWBend         Block = 'J' // 90 degree bend, North-West
	SWBend         Block = '7' // 90 degree bend, South-West
	SEBend         Block = 'F' // 90 degree bend, South-East
	Ground         Block = '.' // Ground
	Start          Block = 'S' // Start
)

type Map [][]Block

func parseInput(input string) Map {
	lines := strings.Split(input, "\n")
	m := make(Map, len(lines))

	for i, line := range lines {
		m[i] = make([]Block, len(line))
		for j, block := range line {
			m[i][j] = Block(block)
		}
	}

	return m
}

func (m Map) PtIn(pt image.Point) bool {
	return pt.Y >= 0 && pt.Y < len(m) && pt.X >= 0 && pt.X < len(m[pt.Y])
}

func (m Map) At(pt image.Point) Block {
	if !m.PtIn(pt) {
		return 0
	}
	return m[pt.Y][pt.X]
}

func (m Map) String() string {
	return string(bytes.Join([][]byte(m), []byte("\n")))
}

func (m Map) Copy() Map {
	c := make(Map, len(m))
	for i, line := range m {
		c[i] = make([]Block, len(line))
		copy(c[i], line)
	}
	return c
}

var (
	DeltaSouth = image.Point{0, 1}
	DeltaNorth = image.Point{0, -1}
	DeltaEast  = image.Point{1, 0}
	DeltaWest  = image.Point{-1, 0}
)

func (m Map) AdjacentDeltas(pt image.Point) []image.Point {
	if !m.PtIn(pt) {
		return nil
	}

	block := m[pt.Y][pt.X]
	switch block {
	case VerticalPipe:
		return []image.Point{{0, -1}, {0, 1}}
	case HorizontalPipe:
		return []image.Point{{-1, 0}, {1, 0}}
	case NEBend:
		return []image.Point{{0, -1}, {1, 0}}
	case NWBend:
		return []image.Point{{0, -1}, {-1, 0}}
	case SWBend:
		return []image.Point{{0, 1}, {-1, 0}}
	case SEBend:
		return []image.Point{{0, 1}, {1, 0}}
	case Start:
		var traversable []image.Point
		if slices.Contains(m.AdjacentDeltas(pt.Add(DeltaNorth)), DeltaSouth) {
			traversable = append(traversable, DeltaNorth)
		}
		if slices.Contains(m.AdjacentDeltas(pt.Add(DeltaSouth)), DeltaNorth) {
			traversable = append(traversable, DeltaSouth)
		}
		if slices.Contains(m.AdjacentDeltas(pt.Add(DeltaEast)), DeltaWest) {
			traversable = append(traversable, DeltaEast)
		}
		if slices.Contains(m.AdjacentDeltas(pt.Add(DeltaWest)), DeltaEast) {
			traversable = append(traversable, DeltaWest)
		}
		log.Printf("Start traversable: %v", traversable)
		return traversable
	default:
		return nil
	}
}

func (m Map) TraverseFrom(start image.Point, f func(pt image.Point, dist int) bool) bool {
	m = m.Copy()

	type traversedPt struct {
		pt   image.Point
		dist int
	}

	queue := []traversedPt{{start, 0}}

	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]

		if !m.PtIn(item.pt) {
			continue
		}

		if !f(item.pt, item.dist) {
			return false
		}

		for _, nextDelta := range m.AdjacentDeltas(item.pt) {
			nextPt := item.pt.Add(nextDelta)
			queue = append(queue, traversedPt{nextPt, item.dist + 1})
		}

		m[item.pt.Y][item.pt.X] = Ground
	}

	return true
}

func (m Map) TraverseMaxDist(start image.Point) int {
	var maxDist int
	m.TraverseFrom(start, func(pt image.Point, dist int) bool {
		maxDist = max(maxDist, dist)
		return true
	})
	return maxDist
}

func (m Map) FindStart() image.Point {
	for y, line := range m {
		xi := strings.IndexByte(string(line), byte(Start))
		if xi != -1 {
			return image.Point{xi, y}
		}
	}
	panic("no start found")
}

// PipeBitmap is a blown up map, where each block is represented by a 3x3 pixel
// square. This square accurately describes the shape of the block so that
// a path is a regular line.
type PipeBitmap = image.Paletted

// NewPipeBitmap creates a PipeBitmap from a Map.
func NewPipeBitmap(m Map) *PipeBitmap {
	img := image.NewPaletted(image.Rect(0, 0, len(m[0])*3, len(m)*3), []color.Color{
		color.White, // ground
		color.Black, // pipe
		color.RGBA{255, 0, 0, 255},
		color.RGBA{0, 255, 0, 255},
		color.RGBA{0, 0, 255, 255},
	})

	for y, line := range m {
		for x := range line {
			corner := image.Point{x * 3, y * 3}
			middle := corner.Add(image.Point{1, 1})

			deltas := m.AdjacentDeltas(image.Point{x, y})
			if len(deltas) > 0 {
				img.SetColorIndex(middle.X, middle.Y, 1)
				for _, delta := range deltas {
					switch delta {
					case DeltaNorth:
						img.SetColorIndex(middle.X, middle.Y-1, 1)
					case DeltaSouth:
						img.SetColorIndex(middle.X, middle.Y+1, 1)
					case DeltaEast:
						img.SetColorIndex(middle.X+1, middle.Y, 1)
					case DeltaWest:
						img.SetColorIndex(middle.X-1, middle.Y, 1)
					}
				}
			}
		}
	}

	return img
}

type FilledPipeBitmap struct {
	*PipeBitmap
	FilledColors []color.Color
}

// FloodFill fills the bitmap with the given color, starting at the given
// point. It returns the number of pixels filled.
func FloodFill(img *PipeBitmap, start image.Point, color uint8) int {
	if color == 0 {
		panic("color 0 is reserved")
	}

	if img.ColorIndexAt(start.X, start.Y) != 0 {
		// Start point is already filled.
		log.Printf(
			"start point %v is already filled with color %v",
			start, img.ColorIndexAt(start.X, start.Y))
		return 0
	}

	var filled int
	queue := []image.Point{start}
	for len(queue) > 0 {
		pt := queue[0]
		queue = queue[1:]

		if !pt.In(img.Bounds()) || img.ColorIndexAt(pt.X, pt.Y) != 0 {
			continue
		}
		filled++
		img.SetColorIndex(pt.X, pt.Y, color)

		queue = append(queue, pt.Add(DeltaNorth))
		queue = append(queue, pt.Add(DeltaSouth))
		queue = append(queue, pt.Add(DeltaEast))
		queue = append(queue, pt.Add(DeltaWest))
	}

	log.Println("filled", start, "with color", color, "for", filled, "pixels")
	return filled
}

func floodFill(img *PipeBitmap, pt image.Point, colorIx int) {

}

// getSinglePaletteColor returns the index of the color in the palette at the
// given rectangle. It returns 0 if the rectangle contains multiple colors.
func getSinglePaletteColor(img *image.Paletted, r image.Rectangle) uint8 {
	var color uint8
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			if color == 0 {
				color = img.ColorIndexAt(x, y)
			} else if color != img.ColorIndexAt(x, y) {
				return 0
			}
		}
	}
	return color
}

func part1(stdin string) int {
	m := parseInput(stdin)
	s := m.FindStart()
	return m.TraverseMaxDist(s)
}

func part2(stdin string) int {
	m := parseInput(stdin)
	s := m.FindStart()

	// "Any tile that isn't part of the main loop can count as being enclosed by
	// the loop", so any tile that is not the main loop is just junk and can be
	// erased completely.

	// We'll start off by tracking all the points that are part of the main
	// loop.
	loopPoints := map[image.Point]struct{}{}
	m.TraverseFrom(s, func(pt image.Point, dist int) bool {
		loopPoints[pt] = struct{}{}
		return true
	})

	// Then, we'll just eviscerate all the points that aren't part of the main
	// loop.
	for y, line := range m {
		for x := range line {
			pt := image.Point{x, y}
			if _, ok := loopPoints[pt]; !ok {
				m[y][x] = Ground
			}
		}
	}

	// Create a zoomed-in bitmap of the map. It basically translates each pipe
	// character (e.g. L, 7, etc.) to actual lines on an image.
	bitmap := NewPipeBitmap(m)
	bitmapS := s.Mul(3).Add(image.Point{1, 1})

	// We can then run our flood fill algorithm to figure out which color each
	// area is.
	corners := []image.Point{
		bitmapS.Add(image.Point{-1, -1}),
		bitmapS.Add(image.Point{+1, -1}),
		bitmapS.Add(image.Point{-1, +1}),
		bitmapS.Add(image.Point{+1, +1}),
	}
	// This really never goes beyond 2 colors, since our input is well-crafted
	// enough, and we already cleared out all the junk.
	fillColor := uint8(2)
	for _, corner := range corners {
		if FloodFill(bitmap, corner, fillColor) == 0 {
			continue
		}
		fillColor++
	}

	// Write this so I can show it off :)
	pngFile := aocutil.E2(os.Create("bitmap.png"))
	defer pngFile.Close()
	aocutil.E1(png.Encode(pngFile, (*image.Paletted)(bitmap)))

	// Start counting blocks of 3x3 pixels and count whole color blocks, aka
	// blocks that only contain one color.
	wholeBlocks := make(map[uint8]int, len(bitmap.Palette))
	for y := 0; y < bitmap.Bounds().Dy(); y += 3 {
		for x := 0; x < bitmap.Bounds().Dx(); x += 3 {
			r := image.Rect(x, y, x+3, y+3)
			c := getSinglePaletteColor(bitmap, r)
			if c == 0 {
				continue
			}
			wholeBlocks[c]++
		}
	}

	// Pick the color at (0, 0). This is the "outside" color.
	outsideColor := bitmap.ColorIndexAt(0, 0)
	log.Printf("outside color at (0, 0) is %v", outsideColor)

	// We want the inside color, so we pick the color that is not the outside
	// color.
	for color, count := range wholeBlocks {
		if color != outsideColor {
			return count
		}
	}

	log.Fatalln("could not find inside color :(")
	return 0
}
