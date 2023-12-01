package main

import (
	"fmt"
	"strings"

	"github.com/diamondburned/aoc-2022/aocutil"
)

var digits = []string{
	"1", "2", "3", "4", "5", "6", "7", "8", "9",
	"one", "two", "three", "four", "five", "six", "seven", "eight", "nine",
}

type searchDirection int

const (
	start searchDirection = iota
	end
)

func searchNumber(str string, direction searchDirection) (int, bool) {
	d := -1
	n := len(str)
	if direction == end {
		n = -1
	}

	for i, digit := range digits {
		switch direction {
		case start:
			m := strings.Index(str, digit)
			if m == -1 || m > n {
				continue
			}
			n = m
		case end:
			m := strings.LastIndex(str, digit)
			if m == -1 || m < n {
				continue
			}
			n = m
		}
		if i < 9 {
			d = i + 1
		} else {
			d = i - 8
		}
	}

	if d == -1 {
		return 0, false
	}

	return d, true
}

func main() {
	input := aocutil.ReadStdin()
	lines := aocutil.SplitLines(input)

	var sum int
	for _, line := range lines {
		d1, ok1 := searchNumber(line, start)
		d2, ok2 := searchNumber(line, end)
		if !ok1 || !ok2 {
			continue
		}
		sum += d1*10 + d2
	}

	fmt.Println(sum)
	// part1()
	// part2()
}
