package main

import (
	"fmt"
	"io"
	"math"
	"sort"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	input := aocutil.InputString()

	m := ParseMap(input)
	part1(m)
	fmt.Println()
	part2(m)
}

func part1(m Map) {
	fmt.Println("Part 1:")
	nav := m.NavigateFromMe()
	for i, route := range nav.Routes {
		fmt.Printf("route %d (length %d)\n", i, len(route))
		route.Print(m, nav.Source, aocutil.PrefixedStdout("  "))
	}
}

func part2(m Map) {
	fmt.Println("Part 2:")

	navs := m.NavigateFrom('a')
	sort.Slice(navs, func(i, j int) bool {
		return len(navs[i].Routes[0]) < len(navs[j].Routes[0])
	})

	nav := navs[0]
	fmt.Printf("nearest point found at %v\n", nav.Source)
	for j, route := range nav.Routes {
		fmt.Printf("route %d: length %d\n", j, len(route))
		route.Print(m, nav.Source, aocutil.PrefixedStdout("  "))
	}
}

// Coordinate is a coordinate in the map.
type Coordinate struct{ X, Y int }

// Add adds the coordinates together and returns the result.
func (c Coordinate) Add(other Coordinate) Coordinate {
	c.X += other.X
	c.Y += other.Y
	return c
}

// MapValue is a value on the map.
type MapValue byte

const (
	MyLocation MapValue = 'S' // elevation a
	BestSignal MapValue = 'E' // elevation z
)

// Elevation returns the elevation of the point.
func (p MapValue) Elevation() int {
	if p == MyLocation {
		p = 'a'
	}
	if p == BestSignal {
		p = 'z'
	}
	if 'a' <= p && p <= 'z' {
		return int(p - 'a')
	}
	return math.MaxInt // invalid
}

// Map describes a map.
type Map [][]MapValue

// ParseMap parses the map.
func ParseMap(block string) Map {
	lines := aocutil.SplitLines(block)

	w := len(lines[0])
	h := len(lines)

	m := make(Map, h)
	for y, line := range lines {
		m[y] = make([]MapValue, w)
		for x, r := range line {
			m[y][x] = MapValue(r)
		}
	}

	return m
}

// At returns the value at the coordinate.
func (m Map) At(coord Coordinate) MapValue {
	if coord.Y < 0 || coord.Y >= len(m) {
		return 0
	}
	if coord.X < 0 || coord.X >= len(m[coord.Y]) {
		return 0
	}
	return m[coord.Y][coord.X]
}

// Width returns the width of the map.
func (m Map) Width() int { return len(m[0]) }

// Height returns the height of the map.
func (m Map) Height() int { return len(m) }

// Direction is a direction enum.
type Direction byte

const (
	Up    Direction = '^'
	Down  Direction = 'v'
	Left  Direction = '<'
	Right Direction = '>'
)

var directionDeltas = map[Direction]Coordinate{
	Up:    {0, -1},
	Down:  {0, 1},
	Left:  {-1, 0},
	Right: {1, 0},
}

// Route is a route. It consists of directions, each of them relative to the
// previous location.
type Route []Direction

// Add adds a direction to a new route and returns it.
func (r Route) Add(direction Direction) Route {
	cpy := make(Route, len(r)+1)
	copy(cpy, r)
	cpy[len(r)] = direction
	return cpy
}

// Print prints the route using the given map and starting coordinate. The
// output is written to the given writer.
func (r Route) Print(m Map, start Coordinate, out io.Writer) {
	buf := make([][]byte, m.Height())
	for i := range buf {
		buf[i] = make([]byte, m.Width())
		for j := range buf[i] {
			buf[i][j] = '.'
		}
	}

	current := start
	for _, dir := range r {
		buf[current.Y][current.X] = byte(dir)
		current = current.Add(directionDeltas[dir])
	}
	buf[current.Y][current.X] = 'E'

	for _, row := range buf {
		out.Write(row)
		io.WriteString(out, "\n")
	}
}

// Navigation describes a navigation.
type Navigation struct {
	Source Coordinate
	Routes []Route
}

// NavigateFromMe navigates from the current location. The first location found
// on the map is assumed.
func (m Map) NavigateFromMe() Navigation {
	return m.NavigateFrom(MyLocation)[0]
}

// NavigateFrom navigates from the given value. A list of navigations is
// returned, with each one starting from each location of the given value found
// on the map.
func (m Map) NavigateFrom(val MapValue) []Navigation {
	var navs []Navigation
	for y, row := range m {
		for x, point := range row {
			if point == val {
				navigator := newNavigator(m, Coordinate{X: x, Y: y})
				navigator.navigate(Coordinate{X: x, Y: y})
				if len(navigator.routes) > 0 {
					navs = append(navs, Navigation{
						Source: Coordinate{X: x, Y: y},
						Routes: navigator.routes,
					})
				}
			}
		}
	}
	return navs
}

type navigator struct {
	m       Map
	dst     MapValue
	routes  []Route
	visited [][]bool
}

func newNavigator(m Map, source Coordinate) navigator {
	navigator := navigator{
		m:       m,
		dst:     BestSignal,
		visited: make([][]bool, len(m)),
	}
	for i := range navigator.visited {
		navigator.visited[i] = make([]bool, len(m[i]))
	}
	return navigator
}

func (n *navigator) navigate(srcPos Coordinate) {
	type queueItem struct {
		pos   Coordinate
		route Route
	}

	queue := make([]queueItem, 0, 4)
	queue = append(queue, queueItem{srcPos, nil})

	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:] // pop

		if n.m.At(item.pos) == n.dst {
			n.routes = append(n.routes, item.route) // save route
			return
		}

		// Recurse in 4 directions.
		for _, dir := range []Direction{Up, Down, Left, Right} {
			delta := directionDeltas[dir]
			pnext := item.pos.Add(delta)

			if !n.canClimb(item.pos, pnext) || n.visited[pnext.Y][pnext.X] {
				continue
			}

			queue = append(queue, queueItem{pnext, item.route.Add(dir)})
			n.visited[pnext.Y][pnext.X] = true
		}
	}

	sort.Slice(n.routes, func(i, j int) bool {
		return len(n.routes[i]) < len(n.routes[j])
	})
}

func (n *navigator) canClimb(cur, dst Coordinate) bool {
	curValue := n.m.At(cur)
	dstValue := n.m.At(dst)
	if curValue == 0 || dstValue == 0 {
		return false
	}
	return dstValue.Elevation()-curValue.Elevation() <= 1
}
