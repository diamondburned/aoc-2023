package main

import (
	"slices"
	"strings"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	aocutil.Run(part1, part2)
}

type Sequence []int

// Differences returns the differences between each element in the sequence.
// That is, the returned list contains [s[1]-s[0], s[2]-s[1], ...].
func (s Sequence) Differences() Sequence {
	diffs := make(Sequence, len(s)-1)
	for i := 0; i < len(s)-1; i++ {
		diffs[i] = s[i+1] - s[i]
	}
	return diffs
}

// AllZeroes returns true if all elements in the sequence are zero.
func (s Sequence) AllZeroes() bool {
	for _, n := range s {
		if n != 0 {
			return false
		}
	}
	return true
}

// AllDifferences continues to call Differences until the returned sequence
// contains all zeroes, indicating that it has found the differences between
// the elements of the sequence. The returned list contains s itself, followed
// by the differences between each element in the sequence, followed by the
// differences between each element in the previous sequence, and so on.
func (s Sequence) AllDifferences() []Sequence {
	diffs := []Sequence{s}
	for {
		s = s.Differences()
		diffs = append(diffs, s)
		if s.AllZeroes() {
			break
		}
	}
	return diffs
}

func parseInput(input string) []Sequence {
	var sequences []Sequence
	for _, line := range aocutil.SplitLines(input) {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		nums := aocutil.Atois[int](fields)
		sequences = append(sequences, nums)
	}

	return sequences
}

// extrapolateSequence extrapolates the differences between each element in
// the sequence so that there is one more element than the previous sequence in
// all sequences.
func extrapolateSequence(diffs []Sequence, left, right bool) []Sequence {
	// Copy so we can append to the last sequence.
	diffs = append([]Sequence(nil), diffs...)

	// Assert that the last sequence is all zeroes.
	zeroes := diffs[len(diffs)-1]
	if !zeroes.AllZeroes() {
		panic("last sequence is not all zeroes")
	}

	// This is just a bunch of zeroes. Add one more zero to the end.
	if right {
		zeroes = append(zeroes, 0)
	}
	if left {
		zeroes = slices.Insert(zeroes, 0, 0)
	}

	for i := len(diffs) - 2; i >= 0; i-- {
		lowerDiffs := diffs[i+1]
		currentDiffs := diffs[i]
		if right {
			// Extrapolate by getting the matching value of the lower sequence
			// and adding that with the current value.
			lowerDiff := lowerDiffs[len(lowerDiffs)-1]
			currentDiff := currentDiffs[len(currentDiffs)-1]

			extrapolated := currentDiff + lowerDiff
			diffs[i] = append(diffs[i], extrapolated)
		}
		if left {
			// Extrapolate by getting the matching value of the lower sequence
			// and subtracting that with the current value.
			lowerDiff := lowerDiffs[0]
			currentDiff := currentDiffs[0]

			extrapolated := currentDiff - lowerDiff
			diffs[i] = append(Sequence{extrapolated}, diffs[i]...)
		}
	}

	return diffs
}

func part1(stdin string) int {
	sequences := parseInput(stdin)

	var sum int
	for _, sequence := range sequences {
		diffs := sequence.AllDifferences()
		extra := extrapolateSequence(diffs, false, true)
		sum += extra[0][len(extra[0])-1]
	}

	return sum
}

func part2(stdin string) int {
	sequences := parseInput(stdin)

	var sum int
	for _, sequence := range sequences {
		diffs := sequence.AllDifferences()
		extra := extrapolateSequence(diffs, true, false)
		sum += extra[0][0]
	}

	return sum
}
