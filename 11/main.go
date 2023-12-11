package main

import (
	"bytes"
	"image"
	"slices"

	"github.com/diamondburned/aoc-2022/aocutil"
	"golang.org/x/exp/constraints"
	"gonum.org/v1/gonum/stat/combin"
)

func main() {
	aocutil.Run(part1, part2)
}

func manhattanDistance(a, b image.Point) int {
	return aocutil.Abs(a.X-b.X) + aocutil.Abs(a.Y-b.Y)
}

const (
	EmptySpace byte = '.'
	Galaxy     byte = '#'
)

type GalaxyImage aocutil.Set[image.Point]

func parseInput(input string) GalaxyImage {
	lines := aocutil.SplitLines(input)
	img := aocutil.NewSet[image.Point](0)
	for y, line := range lines {
		for x, b := range []byte(line) {
			if b == Galaxy {
				img.Add(image.Pt(x, y))
			}
		}
	}
	return GalaxyImage(img)
}

func (i GalaxyImage) Add(p image.Point) {
	i[p] = struct{}{}
}

// FurthestGalaxy returns the point that is furthest away from the origin.
func (i GalaxyImage) FurthestGalaxy() image.Point {
	rect := image.Rectangle{}
	for p := range i {
		rect = rect.Union(image.Rectangle{Max: p})
	}
	return rect.Max
}

func (i GalaxyImage) Expand(gap int) GalaxyImage {
	gap--

	maxPt := i.FurthestGalaxy()
	expanded := make(GalaxyImage, len(i))

	yGalaxies := make(map[int][]int, len(i))
	for p := range i {
		yGalaxies[p.Y] = append(yGalaxies[p.Y], p.X)
	}

	yGalaxies2 := make(map[int][]int, len(i))
	yOffset := 0
	for y := 0; y <= maxPt.Y; y++ {
		if xs, ok := yGalaxies[y]; ok {
			yGalaxies2[y+yOffset] = xs
			continue
		}
		yOffset += gap
	}
	yGalaxies = yGalaxies2

	xGalaxies := make(map[int][]int, len(i))
	for y, xs := range yGalaxies {
		for _, x := range xs {
			xGalaxies[x] = append(xGalaxies[x], y)
		}
	}

	xGalaxies2 := make(map[int][]int, len(i))
	xOffset := 0
	for x := 0; x <= maxPt.X; x++ {
		if ys, ok := xGalaxies[x]; ok {
			xGalaxies2[x+xOffset] = ys
			continue
		}
		xOffset += gap
	}
	xGalaxies = xGalaxies2

	for x, ys := range xGalaxies {
		for _, y := range ys {
			expanded.Add(image.Pt(x, y))
		}
	}

	return expanded
}

func mapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

func sortedKeys[K constraints.Ordered, V any](m map[K]V) []K {
	keys := mapKeys(m)
	slices.Sort(keys)
	return keys
}

func (i GalaxyImage) String() string {
	maxPt := i.FurthestGalaxy()
	img := make([][]byte, maxPt.Y+1)
	for y := range img {
		img[y] = make([]byte, maxPt.X+1)
		for x := range img[y] {
			img[y][x] = EmptySpace
		}
		for x := range img[y] {
			if _, ok := i[image.Pt(x, y)]; ok {
				img[y][x] = Galaxy
			}
		}
	}
	return string(bytes.Join(img, []byte{'\n'}))
}

func sumGalaxyDistances(img GalaxyImage) int {
	galaxies := mapKeys(img)
	var sum int
	for _, combination := range combin.Combinations(len(galaxies), 2) {
		a, b := galaxies[combination[0]], galaxies[combination[1]]
		dist := manhattanDistance(a, b)
		sum += dist
	}
	return sum
}

func part1(input string) int {
	img := parseInput(input)
	img = img.Expand(2)
	return sumGalaxyDistances(img)
}

func part2(input string) int {
	img := parseInput(input)
	img = img.Expand(1_000_000)
	return sumGalaxyDistances(img)
}
