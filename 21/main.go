package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
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
		pt = image.Point{
			X: aocutil.PositiveMod(pt.X, m.Bounds.Dx()),
			Y: aocutil.PositiveMod(pt.Y, m.Bounds.Dy()),
		}
	}
	return m.Map2D.At(pt)
}

func (m Map) TraverseSteps() aocutil.Iter2[int, []image.Point] {
	type nodeType struct {
		pt    image.Point
		steps int
	}

	type state uint8
	calculateState := func(steps int) state { return state(steps % 2) }

	type visitedKey struct {
		pt    image.Point
		state state
	}

	return func(yield func(int, []image.Point) bool) {
		queue := []nodeType{{m.Start, 0}}

		var plotSteps int
		var plots []image.Point
		plotMap := aocutil.NewSet[image.Point](0)
		visited := aocutil.NewSet[visitedKey](0)

		for len(queue) > 0 {
			node := queue[0]
			queue = queue[1:]
			state := calculateState(node.steps)

			if !visited.Add(visitedKey{node.pt, state}) {
				continue
			}

			if node.steps != plotSteps {
				// Finish up plot steps.
				for k := range visited {
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
				clear(plotMap)
			}

			if node.steps == plotSteps {
				plots = append(plots, node.pt)
				plotMap.Add(node.pt)
			}

			for _, dir := range aocutil.CardinalDirections {
				pt := node.pt.Add(dir)
				at := m.At(pt)
				if at == GardenPlot || at == Start {
					next := nodeType{pt, node.steps + 1}
					queue = append(queue, next)
				}
			}
		}
	}
}

func part1(input string) int {
	m := parseInput(input)
	for steps, plots := range m.TraverseSteps() {
		if steps == 64 {
			return len(plots)
		}
	}
	panic("unreachable")
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

	x := (target - start) / size // start = 65 + (131 * x)
	y := aocutil.CalculateRegression(regression, x)
	log.Printf("fitted result: f(%v) = %v", x, y)

	return y
}

func drawPlots(m Map, steps int, plots []image.Point) {
	bounds := aocutil.RectangleContainingPoints(plots)
	bounds.Min = bounds.Min.Sub(image.Point{1, 1})
	bounds.Max = bounds.Max.Add(image.Point{1, 1})

	steppedMap := aocutil.NewEmptyMap2D(bounds)
	for pt := range aocutil.PointsWithin(bounds) {
		steppedMap.Set(pt, m.At(pt))
	}
	for _, pt := range plots {
		steppedMap.Set(pt, 'O')
	}

	img := steppedMap.Draw(map[byte]color.RGBA{
		'.': {0, 0, 0, 255},
		'#': {255, 0, 0, 255},
		'O': {255, 255, 255, 255},
	})
	aocutil.SaveImage(img, fmt.Sprintf("output-%d.png", steps))
}
