package main

import (
	"strings"

	"libdb.so/aoc-2023/aocutil"
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
	HasZZZNode bool
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

		if value1 == "ZZZ" || value2 == "ZZZ" {
			m.HasZZZNode = true
		}

		if strings.HasSuffix(name, "A") {
			m.ANodes = append(m.ANodes, name)
		}
	}

	return m
}

func part1(stdin string) int {
	input := parseInput(stdin)
	if !input.HasZZZNode {
		// Part 2 input, so we can't do part 1.
		return -1
	}

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
	aNodes := input.ANodes

	isEndNode := func(node string) bool {
		return strings.HasSuffix(node, "Z")
	}

	aSteps := make([]int, len(aNodes))
	for i, node := range aNodes {
		// Find the path for all A nodes.
		for !isEndNode(node) {
			for _, dir := range input.Directions {
				possible := input.Nodes[node]
				switch dir {
				case Left:
					node = possible[0]
				case Right:
					node = possible[1]
				}
				aSteps[i]++
			}
		}
	}

	lcm := aocutil.LCM(aSteps...)
	return lcm
}
