package main

import (
	"strings"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	aocutil.Run(part1, part2)
}

type Direction = byte

const (
	Left  Direction = 'L'
	Right Direction = 'R'
)

type MapNode [2]string

type Map struct {
	Directions []Direction
	Nodes      map[string]MapNode
	ANodes     []string
}

func parseInput(input string) Map {
	blocks := strings.Split(input, "\n\n")
	blocks = aocutil.FilterEmptyStrings(blocks)

	m := Map{Directions: []Direction(blocks[0])}

	nodeLines := aocutil.SplitLines(blocks[1])
	m.Nodes = make(map[string]MapNode, len(nodeLines))

	for _, line := range nodeLines {
		name, values, _ := strings.Cut(line, " = ")
		value1, value2, _ := strings.Cut(strings.Trim(values, "()"), ", ")
		m.Nodes[name] = MapNode{value1, value2}

		if strings.HasSuffix(name, "A") {
			m.ANodes = append(m.ANodes, name)
		}
	}

	return m
}

func part1(stdin string) int {
	return -1

	input := parseInput(stdin)
	_ = input

	const start = "AAA"
	const end = "ZZZ"

	var steps int

	current := start
	for current != end {
		for _, dir := range input.Directions {
			possible := input.Nodes[current]
			switch dir {
			case Left:
				current = possible[0]
			case Right:
				current = possible[1]
			}
			steps++
		}
	}

	return steps
}

func part2(stdin string) int {
	input := parseInput(stdin)
	_ = input

	aNodes := input.ANodes

	// pathFinished returns true if the path has reached the end.
	pathFinished := func(path []string) bool {
		return strings.HasSuffix(path[len(path)-1], "Z")
	}

	// nodesEnded := func(path []string) bool {
	// 	for _, node := range aNodes {
	// 		if !strings.HasSuffix(node, "Z") {
	// 			return false
	// 		}
	// 	}
	// 	return true
	// }

	// Keep track of the paths that each A node has taken.
	// If we detect a cycle, then we should reuse the path.
	// The steps is the length of the path.
	aPaths := make(map[string][]string, len(aNodes))

	for _, aNode := range aNodes {
		aPath := []string{aNode}
		// Find the path for all A nodes.
		for !pathFinished(aPath) {
			for _, dir := range input.Directions {
				node := aPath[len(aPath)-1]

				possible := input.Nodes[node]
				switch dir {
				case Left:
					node = possible[0]
				case Right:
					node = possible[1]
				}

				aPath = append(aPath, node)
			}
		}
		aPaths[aNode] = aPath[1:]
	}

	aSteps := make([]int, len(aNodes))
	for i, aNode := range aNodes {
		aSteps[i] = len(aPaths[aNode])
	}

	lcm := LCM(aSteps...)
	return lcm
}

// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(integers ...int) int {
	a := integers[0]
	b := integers[1]
	integers = integers[2:]

	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}
