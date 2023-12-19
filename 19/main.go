package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	aocutil.Run(part1, part2)
}

type RatingSystem struct {
	Workflows map[string][]WorkflowStep
	Parts     []Part
}

type PartCategory int

const (
	CategoryX PartCategory = iota
	CategoryM
	CategoryA
	CategoryS
	categoryMax
)

func ParsePartCategory(s string) PartCategory {
	s = strings.ToUpper(s)
	switch s {
	case "X":
		return CategoryX
	case "M":
		return CategoryM
	case "A":
		return CategoryA
	case "S":
		return CategoryS
	default:
		log.Panicf("invalid part category: %s", s)
		return categoryMax
	}
}

func (c PartCategory) String() string {
	switch c {
	case CategoryX:
		return "X"
	case CategoryM:
		return "M"
	case CategoryA:
		return "A"
	case CategoryS:
		return "S"
	default:
		log.Panicf("invalid part category: %d", c)
		return ""
	}
}

type Part [categoryMax]int

func (p Part) String() string {
	strs := make([]string, len(p))
	for i, v := range p {
		strs[i] = fmt.Sprintf("%s=%d", PartCategory(i), v)
	}
	return strings.Join(strs, ",")
}

type WorkflowStep struct {
	Category   PartCategory
	Comparison Comparison
	Value      int
	Action     Action
}

// IsUnconditional returns true if the step is unconditional, meaning that the
// action will always be taken.
func (s WorkflowStep) IsUnconditional() bool {
	return s.Category == categoryMax
}

// Check returns true if the step's condition is met.
func (s WorkflowStep) Check(part Part) bool {
	if s.IsUnconditional() {
		return true
	}
	return s.Comparison.Compare(part[s.Category], s.Value)
}

type Comparison byte

const (
	ComparisonLessThan    Comparison = '<'
	ComparisonGreaterThan Comparison = '>'
)

func (c Comparison) Compare(a, b int) bool {
	switch c {
	case ComparisonLessThan:
		return a < b
	case ComparisonGreaterThan:
		return a > b
	default:
		log.Panicf("invalid comparison: %c", c)
		return false
	}
}

type Action string

const (
	ActionAccept Action = "A"
	ActionReject Action = "R"
	// anything else is a workflow name
)

func parseInput(input string) RatingSystem {
	blocks := aocutil.SplitBlocks(input)
	return RatingSystem{
		Workflows: parseWorkflows(blocks[0]),
		Parts:     parseParts(blocks[1]),
	}
}

func parseWorkflows(input string) map[string][]WorkflowStep {
	lines := aocutil.SplitLines(input)
	workflows := make(map[string][]WorkflowStep, len(lines))
	for _, line := range lines {
		name, value, _ := strings.Cut(line, "{")
		value = strings.TrimSuffix(value, "}")

		checks := strings.Split(value, ",")
		steps := make([]WorkflowStep, len(checks))

		for i, check := range checks {
			var step WorkflowStep

			condition, action, ok := strings.Cut(check, ":")
			if ok {
				step.Category = ParsePartCategory(string(condition[0]))
				step.Comparison = Comparison(condition[1])
				step.Value = aocutil.Atoi[int](string(condition[2:]))
				step.Action = Action(action)
			} else {
				step.Category = categoryMax
				step.Action = Action(check)
			}

			steps[i] = step
		}
		workflows[name] = steps
	}
	return workflows
}

func parseParts(input string) []Part {
	lines := aocutil.SplitLines(input)
	parts := make([]Part, len(lines))
	for i, line := range lines {
		line = strings.Trim(line, "{}")
		vars := strings.Split(line, ",")
		for _, variable := range vars {
			name, value, _ := strings.Cut(variable, "=")
			variable := ParsePartCategory(name)
			parts[i][variable] = aocutil.Atoi[int](value)
		}
	}
	return parts
}

func (s RatingSystem) IsAccepted(part Part) bool {
	workflow := s.Workflows["in"]

workflowEval:
	for {
		for _, step := range workflow {
			if !step.Check(part) {
				continue
			}

			switch step.Action {
			case ActionAccept:
				return true
			case ActionReject:
				return false
			}

			workflow = s.Workflows[string(step.Action)]
			continue workflowEval
		}
	}
}

func part1(input string) int {
	rs := parseInput(input)

	var totalRating int
	for _, part := range rs.Parts {
		if rs.IsAccepted(part) {
			totalRating += aocutil.Sum[int](part[:])
		}
	}

	return totalRating
}

type Interval struct {
	Min int // inclusive
	Max int // exclusive
}

func (iv Interval) Count() int {
	return iv.Max - iv.Min
}

func (iv Interval) In(v int) bool {
	return v >= iv.Min && v < iv.Max
}

type PartIntervals [categoryMax]Interval

// Count returns the number of possible combinations of parts.
func (p PartIntervals) Count() int {
	count := 1
	for _, iv := range p {
		count *= iv.Count()
	}
	return count
}

func (p PartIntervals) String() string {
	strs := make([]string, len(p))
	for i, iv := range p {
		strs[i] = fmt.Sprintf("%s=[%d,%d)", PartCategory(i), iv.Min, iv.Max)
	}
	return strings.Join(strs, ",")
}

func allAcceptedCombinations(rs RatingSystem) int {
	return acceptedCombinations(rs, Action("in"), PartIntervals{
		{1, 4001},
		{1, 4001},
		{1, 4001},
		{1, 4001},
	})
}

func acceptedCombinations(rs RatingSystem, action Action, partIntervals PartIntervals) int {
	switch action {
	case ActionAccept:
		return partIntervals.Count()
	case ActionReject:
		return 0
	}

	var count int
	for _, step := range rs.Workflows[string(action)] {
		if step.IsUnconditional() {
			count += acceptedCombinations(rs, step.Action, partIntervals)
			break
		}

		l := partIntervals
		r := partIntervals

		switch step.Comparison {
		case ComparisonLessThan:
			l[step.Category].Max = step.Value
			r[step.Category].Min = step.Value

			partIntervals = r
			count += acceptedCombinations(rs, step.Action, l)

		case ComparisonGreaterThan:
			l[step.Category].Min = step.Value + 1
			r[step.Category].Max = step.Value + 1

			partIntervals = r
			count += acceptedCombinations(rs, step.Action, l)
		}
	}

	return count
}

/*
func allAcceptedCombinations(rs RatingSystem) int {
	type queueItem struct {
		workflow  string
		intervals PartIntervals
	}

	biggestPartIntervals := PartIntervals{
		{1, 4001},
		{1, 4001},
		{1, 4001},
		{1, 4001},
	}

	// var acceptedIntervals []PartIntervals
	var count int
	stack := []queueItem{{"in", biggestPartIntervals}}

	debug := func() {
		log.Println("debug:")
		// log.Println("  accepted:")
		// for _, partIntervals := range acceptedIntervals {
		// 	log.Print("    ", partIntervals)
		// }
		log.Println("  queue:")
		for _, item := range stack {
			log.Printf("    %s: %s", item.workflow, item.intervals)
		}
	}

	for len(stack) > 0 {
		debug()
		item := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

	}

	// log.Printf("acceptedIntervals: %d", len(acceptedIntervals))
	//
	// var count int
	// for i, accepted := range acceptedIntervals {
	// 	// log.Printf("%s (%d)", accepted, accepted.Count())
	// 	// log.Println()
	// 	count += accepted.Count()
	// 	for _, accepted2 := range acceptedIntervals[:i] {
	// 		// log.Printf("    %s", accepted)
	// 		// log.Printf("  ∩ %s", accepted2)
	// 		// log.Printf("    %s (%d)", accepted.Intersect(accepted2), accepted.Intersect(accepted2).Count())
	// 		// log.Println()
	// 		count -= accepted.Intersect(accepted2).Count()
	// 	}
	// }

	return count
}
*/

func part2(input string) int {
	rs := parseInput(input)
	return allAcceptedCombinations(rs)
}
