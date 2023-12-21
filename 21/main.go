package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"slices"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	aocutil.Run(part1, part2)
}

const (
	Start      = 'S'
	GardenPlot = '.'
	Rock       = '#'
)

type Map struct {
	aocutil.Map2D
	Start    image.Point
	Infinite bool
}

func parseInput(input string) Map {
	m := aocutil.NewMap2D(input)
	var start image.Point
	for pt, b := range m.All() {
		if b == Start {
			start = pt
			break
		}
	}
	return Map{
		Map2D: m,
		Start: start,
	}
}

// At assumes the map is infinite in all directions.
func (m Map) At(pt image.Point) byte {
	if m.Infinite {
		pt = m.loopPointBack(pt)
	}
	return m.Map2D.At(pt)
}

func (m Map) loopPointBack(pt image.Point) image.Point {
	return image.Point{
		X: aocutil.PositiveMod(pt.X, m.Bounds.Dx()),
		Y: aocutil.PositiveMod(pt.Y, m.Bounds.Dy()),
	}
}

func (m Map) loopTileID(pt image.Point) image.Point {
	return image.Point{
		X: pt.X / m.Bounds.Dx(),
		Y: pt.Y / m.Bounds.Dy(),
	}
}

func (m Map) Clone() Map {
	m.Map2D = m.Map2D.Clone()
	return m
}

func (m Map) CountSteps(stepsNeeded int) []image.Point {
	type nodeType struct {
		pt    image.Point
		steps int
	}

	queue := []nodeType{{m.Start, 0}}
	var plots []image.Point

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		// log.Printf("node: %+v", node)

		if node.steps == stepsNeeded {
			plots = append(plots, node.pt)
			continue
		}

		pt := node.pt

		for _, dir := range aocutil.CardinalDirections {
			pt := pt.Add(dir)
			at := m.At(pt)
			if at != GardenPlot && at != Start {
				continue
			}
			next := nodeType{pt, node.steps + 1}
			if slices.Contains(queue, next) {
				continue
			}
			queue = append(queue, next)
		}
	}

	return plots
}

func (m Map) TraverseSteps() aocutil.Iter2[int, []image.Point] {
	type nodeType struct {
		pt    image.Point
		steps int
	}

	type state uint8

	const (
		stateEven state = iota
		stateOdd
	)

	calculateState := func(steps int) state {
		return state(steps % 2)
	}

	type tileKey struct {
		pt    image.Point
		state state
	}

	type visitedKey struct {
		pt    image.Point
		state state
	}

	return func(yield func(int, []image.Point) bool) {
		queue := []nodeType{{m.Start, 0}}

		var plotSteps int
		var plots []image.Point
		plotMap := aocutil.NewSet[image.Point](0)

		visitedStates := make(map[visitedKey]bool)
		// tileMap := make(map[tileKey]bool)

		for len(queue) > 0 {
			node := queue[0]
			queue = queue[1:]
			state := calculateState(node.steps)
			// log.Printf("node: %+v", node)

			if visitedStates[visitedKey{node.pt, state}] {
				continue
			}
			visitedStates[visitedKey{node.pt, state}] = true

			if node.steps != plotSteps {
				// Finish up plot steps.
				for k := range visitedStates {
					if k.state == calculateState(plotSteps) && !plotMap.Has(k.pt) {
						plotMap.Add(k.pt)
						plots = append(plots, k.pt)
					}
				}
				if !yield(plotSteps, plots) {
					return
				}
				plotSteps = node.steps
				plots = plots[:0]
				plotMap = aocutil.NewSet[image.Point](cap(plots))
			}

			if node.steps == plotSteps {
				plots = append(plots, node.pt)
				plotMap.Add(node.pt)
			}

			pt := node.pt

			for _, dir := range aocutil.CardinalDirections {
				pt := pt.Add(dir)
				at := m.At(pt)
				if at != GardenPlot && at != Start {
					continue
				}
				next := nodeType{pt, node.steps + 1}
				queue = append(queue, next)
			}
		}
	}
}

func part1(input string) int {
	m := parseInput(input)
	m.Infinite = true
	plots := m.CountSteps(64)
	return len(plots)
}

func newBoundsFromPts(pts []image.Point) image.Rectangle {
	var bounds image.Rectangle
	for _, pt := range pts {
		bounds = bounds.Union(image.Rectangle{pt, pt.Add(image.Point{1, 1})})
	}
	return bounds
}

func part2(input string) int {
	m := parseInput(input)
	m.Infinite = true

	const target = 26501365
	const fitPoints = 3
	size := m.Bounds.Max.X
	start := m.Start.X

	interests := make([]int, fitPoints)
	interestedPlots := make([]int, 0, fitPoints)
	for i := range interests {
		interests[i] = start + (size * i) // 65 + (131 * i)
	}

	for steps, plots := range m.TraverseSteps() {
		if slices.Contains(interests, steps) {
			log.Printf("steps: %d, plots: %d", steps, len(plots))
			interestedPlots = append(interestedPlots, len(plots))
			drawPlots(m, steps, plots)
		}
		if len(interestedPlots) == len(interests) {
			break
		}
	}

	regression := aocutil.Polyfit(
		aocutil.Range(0, float64(len(interests))).All(),
		aocutil.Map(interestedPlots, func(i int) float64 { return float64(i) }),
		2)
	regression = aocutil.RoundedRegression(regression)
	log.Printf("fitted polynomial function: %s", regression)

	x := (target - start) / size
	y := aocutil.CalculateRegression(regression, x)
	log.Printf("fitted result: F(%v) = %v", x, y)

	return y
}

func drawPlots(m Map, steps int, plots []image.Point) {
	bounds := newBoundsFromPts(plots)
	bounds.Min = bounds.Min.Sub(image.Point{1, 1})
	bounds.Max = bounds.Max.Add(image.Point{1, 1})
	steppedMap := aocutil.NewEmptyMap2D(bounds)
	for pt := range aocutil.PointsWithin(bounds) {
		steppedMap.Set(pt, m.At(pt))
	}
	for _, pt := range plots {
		steppedMap.Set(pt, 'O')
	}
	img := drawMap(steppedMap, map[byte]color.RGBA{
		'.': {0, 0, 0, 255},
		'#': {255, 0, 0, 255},
		'O': {255, 255, 255, 255},
	})
	f, _ := os.Create(fmt.Sprintf("output-%d.png", steps))
	png.Encode(f, img)
	f.Close()
}

func drawMap(m aocutil.Map2D, colorMap map[byte]color.RGBA) *image.Paletted {
	palette := make(color.Palette, 1, 1+len(colorMap))
	palette[0] = color.RGBA{0, 0, 0, 255}

	colorIx := make(map[byte]uint8)
	for k, v := range colorMap {
		palette = append(palette, v)
		colorIx[k] = uint8(len(palette) - 1)
	}

	img := image.NewPaletted(m.Bounds, palette)
	for pt, v := range m.All() {
		ix := colorIx[v]
		img.SetColorIndex(pt.X, pt.Y, ix)
	}

	return img
}
