package main

import (
	"image"
	"log"
	"strconv"
	"strings"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	aocutil.Run(part1, part2)
}

func parseInput(input string) excavationPlan {
	var path []image.Point
	var pos image.Point

	for _, line := range aocutil.SplitLines(input) {
		parts := strings.Fields(line)
		units := aocutil.Atoi[int](parts[1])
		var direction image.Point
		switch parts[0] {
		case "U":
			direction = aocutil.VecUp
		case "D":
			direction = aocutil.VecDown
		case "L":
			direction = aocutil.VecLeft
		case "R":
			direction = aocutil.VecRight
		}
		path = append(path, pos)
		pos = pos.Add(direction.Mul(units))
	}

	return excavationPlan{
		path: path,
	}
}

type excavationPlan struct {
	path []image.Point
}

func parseInputPart2(input string) excavationPlan {
	var plan excavationPlan
	var pos image.Point

	for _, line := range aocutil.SplitLines(input) {
		parts := strings.Fields(line)
		magic := strings.Trim(parts[2], "()")
		magic = magic[1:]

		units, _ := strconv.ParseInt(magic[:5], 16, 64)

		instruction := magic[5:]
		var direction image.Point
		switch instruction {
		case "0":
			direction = aocutil.VecRight
		case "1":
			direction = aocutil.VecDown
		case "2":
			direction = aocutil.VecLeft
		case "3":
			direction = aocutil.VecUp
		}

		plan.path = append(plan.path, pos)
		pos = pos.Add(direction.Mul(int(units)))
	}

	return plan
}

func countPointsBetween(a, b image.Point) int {
	if a.X == b.X {
		return aocutil.Abs(a.Y - b.Y)
	}
	if a.Y == b.Y {
		return aocutil.Abs(a.X - b.X)
	}
	return aocutil.GCD(aocutil.Abs(a.X-b.X), aocutil.Abs(a.Y-b.Y))
}

// area calculates the area of the polygon defined by the given path.
// It uses the Shoelace formula.
func area(path []image.Point) int {
	var sum int
	for i := 0; i < len(path); i++ {
		j := (i + 1) % len(path)
		log.Printf("shoelace: %d, %d", i, j)
		sum += path[i].X * path[j].Y
		sum -= path[j].X * path[i].Y
	}
	return aocutil.Abs(sum) / 2
}

func countExteriorPoints(path []image.Point) int {
	var count int
	for i := 0; i < len(path); i++ {
		count += countPointsBetween(path[i], path[(i+1)%len(path)])
	}
	return count
}

func countInteriorPoints(path []image.Point) int {
	exteriorPoints := countExteriorPoints(path)
	log.Printf("exterior points: %d", exteriorPoints)
	area := area(path)
	log.Printf("area: %d", area)
	return (2*area - exteriorPoints + 2) / 2
}

func part1(input string) int {
	plan := parseInput(input)
	log.Printf("plan: %v", plan)
	return countExteriorPoints(plan.path) + countInteriorPoints(plan.path)
}

func part2(input string) int {
	plan := parseInputPart2(input)
	log.Printf("plan: %v", plan)
	return countExteriorPoints(plan.path) + countInteriorPoints(plan.path)
}
