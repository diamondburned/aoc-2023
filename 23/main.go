package main

import (
	"fmt"
	"image"
	"log"
	"slices"
	"strings"

	. "libdb.so/aoc-2023/aocutil"
)

func main() {
	Run(part1, part2)
}

const (
	Path       = '.'
	Forest     = '#'
	SlopeUp    = '^'
	SlopeDown  = 'v'
	SlopeLeft  = '<'
	SlopeRight = '>'
)

type HikingMap struct {
	Map2D
	Entrance image.Point
	Exit     image.Point
}

func parseInput(input string) HikingMap {
	m := NewMap2D(input)
	return HikingMap{
		Map2D:    m,
		Entrance: image.Point{X: 1, Y: 0},
		Exit:     image.Point{X: m.Bounds.Dx() - 2, Y: m.Bounds.Dy() - 1},
	}
}

func (m HikingMap) neighboringPoints(pt image.Point, ignoreSlopes bool) Iter[image.Point] {
	directions := CardinalDirections
	if !ignoreSlopes {
		switch m.At(pt) {
		case SlopeUp:
			directions = []image.Point{VecUp}
		case SlopeDown:
			directions = []image.Point{VecDown}
		case SlopeLeft:
			directions = []image.Point{VecLeft}
		case SlopeRight:
			directions = []image.Point{VecRight}
		}
	}
	return func(yield func(image.Point) bool) {
		for _, dir := range directions {
			pt := pt.Add(dir)
			at := m.At(pt)
			if at == 0 || at == Forest {
				continue
			}
			if !yield(pt) {
				return
			}
		}
	}
}

type mazeGraph map[image.Point][]mazeConjunction

type mazeConjunction struct {
	pt   image.Point
	dist int
}

func extractMazeAsGraph(m HikingMap, ignoreSlopes bool) mazeGraph {
	type node struct {
		pt     image.Point
		prev   image.Point
		origin image.Point
		steps  int
	}

	type seenKey struct {
		pt     image.Point
		origin image.Point
	}

	maze := make(mazeGraph)
	seen := NewSet[seenKey](0)
	start := node{
		pt:     m.Entrance,
		origin: m.Entrance,
		prev:   image.Point{-1, -1},
	}

	BFS(start, func(n node) []node {
		pts := m.neighboringPoints(n.pt, ignoreSlopes).
			Filter(func(pt image.Point) bool { return n.prev != pt }).
			All()
		if len(pts) == 1 {
			// Just a normal path. Keep following it.
			return []node{{
				pt:     pts[0],
				prev:   n.pt,
				steps:  n.steps + 1,
				origin: n.origin,
			}}
		}

		// Conjunction
		if !seen.Add(seenKey{n.pt, n.origin}) {
			return nil
		}

		maze[n.origin] = append(maze[n.origin], mazeConjunction{
			pt:   n.pt,
			dist: n.steps,
		})

		queue := make([]node, 0, len(pts))
		for _, pt := range pts {
			queue = append(queue, node{
				pt:     pt,
				prev:   n.pt,
				origin: n.pt,
				steps:  1,
			})
		}

		return queue
	}).All()

	return maze
}

func mazeConjunctionsToGraphviz(maze mazeGraph, entrance, exit image.Point) string {
	ptName := func(pt image.Point) string {
		switch pt {
		case entrance:
			return "entrance"
		case exit:
			return "exit"
		default:
			return pt.String()
		}
	}
	var b strings.Builder
	b.WriteString("digraph G {\n")
	for origin, pts := range maze {
		for _, pt := range pts {
			fmt.Fprintf(&b,
				`  %q -> %q [label=%d]`+"\n",
				ptName(origin), ptName(pt.pt), pt.dist)
		}
	}
	b.WriteString("}\n")
	return b.String()
}

func countLongestPathInGraph(maze mazeGraph, entrance, exit image.Point) int {
	type walkingState struct {
		path  []image.Point
		steps int
	}

	var maxSteps int
	stack := []walkingState{{
		path: []image.Point{entrance},
	}}

	for len(stack) > 0 {
		w := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		pt := w.path[len(w.path)-1]
		if pt == exit {
			if w.steps > maxSteps {
				maxSteps = w.steps
				log.Printf("found a longer path: %d steps", maxSteps)
			}
			continue
		}

		for _, next := range maze[pt] {
			if slices.Contains(w.path, next.pt) {
				continue
			}
			stack = append(stack, walkingState{
				path:  append(slices.Clone(w.path), next.pt),
				steps: w.steps + next.dist,
			})
		}
	}

	return maxSteps
}

func part1(input string) int {
	m := parseInput(input)
	g := extractMazeAsGraph(m, false)
	return countLongestPathInGraph(g, m.Entrance, m.Exit)
}

func part2(input string) int {
	m := parseInput(input)
	g := extractMazeAsGraph(m, true)
	// fmt.Print(mazeConjunctionsToGraphviz(g, m.Entrance, m.Exit))
	return countLongestPathInGraph(g, m.Entrance, m.Exit)
}
