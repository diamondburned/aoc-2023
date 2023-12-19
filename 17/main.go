package main

import (
	"image"
	"math"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	aocutil.Run(part1, part2)
}

type Map struct {
	aocutil.Map2D
}

func parseInput(input string) Map {
	m := aocutil.NewMap2D(input)
	return Map{m}
}

func (m Map) CostAt(pt image.Point) int {
	return int(m.At(pt) - '0')
}

func (m Map) Source() image.Point {
	return m.Bounds.Min
}

func (m Map) Sink() image.Point {
	return m.Bounds.Max.Sub(image.Pt(1, 1))
}

func turns(dir image.Point) [3]image.Point {
	return [3]image.Point{
		dir,
		{+dir.Y, -dir.X},
		{-dir.Y, +dir.X},
	}
}

func leastCostlyPath(m Map, straightMin, straightMax int) int {
	type node struct {
		pt       image.Point
		dir      image.Point
		straight int
	}

	start := m.Source()
	end := m.Sink()

	queue := make([]node, 0, 100)
	queue = append(queue,
		node{start, aocutil.VecRight, 0},
		node{start, aocutil.VecDown, 0},
	)

	costMap := map[node]int{}
	minCost := math.MaxInt

	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]

		if n.pt == end && n.straight >= straightMin {
			minCost = min(minCost, costMap[n])
			continue
		}

		for _, dir := range turns(n.dir) {
			if dir == n.dir && n.straight >= straightMax {
				continue
			}
			if dir != n.dir && n.straight < straightMin {
				continue
			}

			next := node{n.pt.Add(dir), dir, 1}
			if next.dir == n.dir {
				next.straight = n.straight + 1
			}
			if !next.pt.In(m.Bounds) {
				// Out of bounds or we're trying to go back.
				continue
			}

			cost := costMap[n] + m.CostAt(next.pt)
			if nextCost, ok := costMap[next]; ok && nextCost <= cost {
				// We've already visited this point with a lower cost.
				continue
			}

			costMap[next] = cost
			queue = append(queue, next)
		}
	}

	return minCost
}

func part1(input string) int {
	m := parseInput(input)
	c := leastCostlyPath(m, 1, 3)
	return c
}

func part2(input string) int {
	m := parseInput(input)
	c := leastCostlyPath(m, 4, 10)
	return c
}
