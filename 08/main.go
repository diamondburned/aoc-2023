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

	log.Println("counted", nVisible, "trees from the edges")

	// Check all the trees in-between.
	for y := 1; y < mapY-1; y++ {
		for x := 1; x < mapX-1; x++ {
			vis := treeMap.TreeVisibility(x, y)
			if vis == 0 {
				continue
			}
			nVisible++
		}
	}

	fmt.Println("part 1: number of visible trees:", nVisible)

	var maxScenicScore int

	// Check all the trees in-between.
	for y := 1; y < mapY-1; y++ {
		for x := 1; x < mapX-1; x++ {
			score := treeMap.TreeScenicScore(x, y)
			if score > maxScenicScore {
				maxScenicScore = score
			}
		}
	}

	fmt.Println("part 2: max scenic score:", maxScenicScore)
}

type VisibilityAngle uint8

var VisibilityAngles = []VisibilityAngle{
	Left,
	Right,
	Top,
	Bottom,
}

const (
	Left VisibilityAngle = 1 << iota
	Right
	Top
	Bottom
)

func (v VisibilityAngle) String() string {
	if v == 0 {
		return "<not visible>"
	}

	var vals []string
	if v&Left != 0 {
		vals = append(vals, "Left")
	}
	if v&Right != 0 {
		vals = append(vals, "Right")
	}
	if v&Top != 0 {
		vals = append(vals, "Top")
	}
	if v&Bottom != 0 {
		vals = append(vals, "Bottom")
	}
	return strings.Join(vals, "+")
}

type TreeMap [][]int8

func (m TreeMap) At(x, y int) int8 {
	return m[y][x]
}

func (m TreeMap) TreeVisibility(x, y int) VisibilityAngle {
	tree := m.At(x, y)
	if tree == 0 {
		return 0
	}

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
	}

	// Keep going in the direction of the angle until we hit a tree taller or
	// equal to the height.
	for {
		x += deltaX
		y += deltaY

		if x < 0 || y < 0 || y >= len(m) || x >= len(m[y]) {
			return true
		}

		if m.At(x, y) >= height {
			return false
		}
	}
}

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

	log.Printf("tree %d at (%d, %d) has scenic score %d (%v)", tree, x, y, score, scores)

	return score
}

func (m TreeMap) angleScenicScore(x, y int, height int8, angle VisibilityAngle) int {
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
	}

	score := 0
	// Keep going in the direction of the angle until we hit a tree taller or
	// equal to the height.
	for {
		x += deltaX
		y += deltaY

		if x < 0 || y < 0 || y >= len(m) || x >= len(m[y]) {
			break
		}

		otherHeight := m.At(x, y)

		score++
		log.Printf("%d: angle %s: tree %d at (%d, %d) is visible from (%d, %d)", height, angle, m.At(x, y), x, y, x-deltaX, y-deltaY)

		if otherHeight >= height {
			break
		}
	}

	return score
}
