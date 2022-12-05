package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	input := aocutil.ReadFile("input")

	var crateStacksBuffer strings.Builder
	var cratesParsed bool

	var crateStacks CrateStacks
	var rearrangementProcedures []RearrangementProcedure

	scanner := aocutil.NewBytesScanner(input)
	for scanner.Next() {
		token := scanner.Token()
		if !cratesParsed {
			if token != "" {
				crateStacksBuffer.WriteString(token)
				crateStacksBuffer.WriteByte('\n')
			} else {
				cratesParsed = true
				crateStacks = MustParseCrateStacks(crateStacksBuffer.String())
			}
			continue
		}

		if token == "" {
			continue
		}

		proc := MustParseRearrangementProcedure(token)
		rearrangementProcedures = append(rearrangementProcedures, proc)
	}

	move := func(part int) {
		is9001 := part == 2

		crateStacks := crateStacks.Copy()
		for _, proc := range rearrangementProcedures {
			crateStacks.Move(proc, is9001)
		}

		fmt.Printf("part %d: tops: %s", part, string(crateStacks.TopCrates()))
		fmt.Println()
	}

	move(1)
	move(2)
}

// https://regex101.com/r/HvS2Jx/1
var crateRowRe = regexp.MustCompile(`(?:\[([A-Z])\] ?|    ?)`)

// CrateStacks is a slice of crate stacks. It's aligned so that each CrateStack
// slice is one vertical stack:
//
//    [A] [B] [C] [D] [E]
//    [F] [G] [H] [I] [J]
//     1   2   3   4   5
//
// Becomes:
//
//    [0]: [A] [F]
//    [1]: [B] [G]
//    [2]: [C] [H]
//    [3]: [D] [I]
//    [4]: [E] [J]
//
type CrateStacks []CrateStack

// MustParseCrateStacks parses crate stacks from a string. It panics if the
// string is invalid.
func MustParseCrateStacks(input string) CrateStacks {
	input = strings.Trim(input, "\n")
	lines := strings.Split(input, "\n")

	// Parse the last line.
	lastFields := strings.Fields(lines[len(lines)-1])

	// Parse the last number in the last line. This is the number of columns.
	columns := aocutil.Atoi[int](lastFields[len(lastFields)-1])

	// The number of rows will just be the line count excluding the last line.
	rows := len(lines) - 1

	stacks := make([]CrateStack, rows)
	for i := range stacks {
		stacks[i] = make(CrateStack, columns)
	}

	for row, line := range lines[:len(lines)-1] {
		matches := crateRowRe.FindAllStringSubmatch(line, -1)
		aocutil.Assertf(len(matches) == columns,
			"invalid number of matches: %d != %d", len(matches), columns)

		for col, match := range matches {
			if match[1] != "" {
				stacks[col][rows-row-1] = Crate(match[1][0])
			}
		}
	}

	for i := range stacks {
		stacks[i] = aocutil.Trim(stacks[i], 0)
	}

	return stacks
}

// Move moves crates from one stack to another using the given procedure.
// If is9001 is true, then multiple crates are moved at once, preserving the
// order of the crates.
func (s CrateStacks) Move(p RearrangementProcedure, is9001 bool) {
	if is9001 {
		moving := s[p.From-1].Pop(p.Qty)
		s[p.To-1].Push(moving)
	} else {
		for i := 0; i < p.Qty; i++ {
			crates := s[p.From-1].Pop(1)
			s[p.To-1].Push(crates)
		}
	}
}

// TopCrates returns the top crates of each stack.
func (s CrateStacks) TopCrates() []Crate {
	crates := make([]Crate, len(s))
	for i, stack := range s {
		crates[i] = stack[len(stack)-1]
	}
	return crates
}

// Copy returns a deep copy of the stack.
func (s CrateStacks) Copy() CrateStacks {
	stacks := make([]CrateStack, len(s))
	for i, stack := range s {
		stacks[i] = append(CrateStack(nil), stack...)
	}
	return stacks
}

// CrateStack is a vertical stack of crates.
type CrateStack []Crate

// Pop pops the top n crates from the stack.
func (s *CrateStack) Pop(n int) []Crate {
	// Copy the top n crates to a new slice.
	crates := append([]Crate(nil), (*s)[len(*s)-n:]...)
	// Slice the top n crates off the stack.
	*s = (*s)[:len(*s)-n]
	return crates
}

// Push pushes crates onto the stack.
func (s *CrateStack) Push(crates []Crate) {
	*s = append(*s, crates...)
}

type Crate byte

type RearrangementProcedure struct {
	Qty  int
	From int
	To   int
}

// MustParseRearrangementProcedure parses a rearrangement procedure from a
// string. It panics if the string is invalid.
func MustParseRearrangementProcedure(input string) RearrangementProcedure {
	const f = `move %d from %d to %d`
	var proc RearrangementProcedure
	_, err := fmt.Sscanf(input, f, &proc.Qty, &proc.From, &proc.To)
	aocutil.E1(err)
	return proc
}
