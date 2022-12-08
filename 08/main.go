package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	input := aocutil.InputString()
	lines := aocutil.SplitLines(input)

	mapY := len(lines)
	mapX := len(lines[0])

	treeMap := make(TreeMap, mapY)
	for i, line := range lines {
		treeMap[i] = make([]int8, mapX)
		for j, c := range line {
			treeMap[i][j] = aocutil.Atoi[int8](string(c))
		}
	}

	nVisible := 0
	// Top and bottom trees are all visible.
	// Left and right trees are all visible.
	nVisible += mapX + mapX + mapY - 2 + mapY - 2
	// Check all the trees in-between.
	each2D(1, 1, mapX-2, mapY-2, func(x, y int) {
		if treeMap.TreeVisibility(x, y) > 0 {
			nVisible++
		}
	})

	fmt.Println("part 1: number of visible trees:", nVisible)

	var maxScenicScore int
	// Check all the trees in-between.
	each2D(1, 1, mapX-2, mapY-2, func(x, y int) {
		score := treeMap.TreeScenicScore(x, y)
		maxScenicScore = aocutil.Max2(maxScenicScore, score)
	})

	fmt.Println("part 2: max scenic score:", maxScenicScore)
}

// each2D calls f for each x, y in the given range.
func each2D(x1, y1, x2, y2 int, f func(x, y int)) {
	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			f(x, y)
		}
	}
}

// VisibilityAngle is a bitmask of the 4 directions of visibility.
type VisibilityAngle uint8

const (
	Left VisibilityAngle = 1 << iota
	Right
	Top
	Bottom
)

// VisibilityAngles is a list of all the visibility angles.
var VisibilityAngles = []VisibilityAngle{
	Left,
	Right,
	Top,
	Bottom,
}

var visibilityAngleNames = map[VisibilityAngle]string{
	Left:   "left",
	Right:  "right",
	Top:    "top",
	Bottom: "bottom",
}

func (v VisibilityAngle) String() string {
	if v == 0 {
		return "<not visible>"
	}
	var vals []string
	for _, angle := range VisibilityAngles {
		if v&angle != 0 {
			vals = append(vals, visibilityAngleNames[angle])
		}
	}
	return strings.Join(vals, "+")
}

// TreeMap is a 2D array of trees.
type TreeMap [][]int8

// At returns the tree height at the given coordinates.
func (m TreeMap) At(x, y int) int8 {
	return m[y][x]
}

// TreeVisibility returns the visibility angle of the tree at the given
// coordinates.
func (m TreeMap) TreeVisibility(x, y int) VisibilityAngle {
	tree := m.At(x, y)

	var vis VisibilityAngle
	// Check the 4 neighbors.
	for _, angle := range VisibilityAngles {
		if m.angleIsVisible(x, y, tree, angle) {
			vis |= angle
		}
	}

	return vis
}

func (m TreeMap) angleIsVisible(x, y int, height int8, angle VisibilityAngle) bool {
	visible := true
	m.walkAngle(x, y, angle, func(x, y int) bool {
		// break if tree at (x, y) is taller than the given height.
		if m.At(x, y) >= height {
			visible = false
			return false
		}
		return true
	})
	return visible
}

// TreeScenicScore calculates the scenic score of the tree at the given
// coordinates.
func (m TreeMap) TreeScenicScore(x, y int) int {
	tree := m.At(x, y)

	var scores [4]int
	// Check the 4 neighbors.
	for i, angle := range VisibilityAngles {
		scores[i] = m.angleScenicScore(x, y, tree, angle)
	}

	score := scores[0]
	for _, s := range scores[1:] {
		score *= s
	}

	return score
}

func (m TreeMap) angleScenicScore(x, y int, height int8, angle VisibilityAngle) int {
	var score int
	m.walkAngle(x, y, angle, func(x, y int) bool {
		score++
		// break if tree at (x, y) is taller than the given height.
		return m.At(x, y) < height
	})
	return score
}

// walkAngle walks in the given direction until fn returns false.
func (m TreeMap) walkAngle(x, y int, angle VisibilityAngle, fn func(x, y int) bool) {
	var deltaX int
	var deltaY int

	switch angle {
	case Left:
		deltaX = -1
	case Right:
		deltaX = 1
	case Top:
		deltaY = -1
	case Bottom:
		deltaY = 1
	default:
		log.Panicf("walkAngle: invalid angle: %v", angle)
	}

	for {
		x += deltaX
		y += deltaY

		if x < 0 || y < 0 || y >= len(m) || x >= len(m[y]) {
			break
		}

		if !fn(x, y) {
			break
		}
	}
}
