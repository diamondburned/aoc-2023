package main

import (
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

type Part [categoryMax]int

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
		panic("invalid comparison")
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

type PartIntervals [categoryMax]Interval

// Count returns the number of possible combinations of parts.
func (p PartIntervals) Count() int {
	count := 1
	for _, iv := range p {
		count *= iv.Count()
	}
	return count
}

func allAcceptedCombinations(rs RatingSystem) int {
	type queueItem struct {
		action    Action
		intervals PartIntervals
	}

	var count int
	stack := []queueItem{{"in", PartIntervals{
		{1, 4001},
		{1, 4001},
		{1, 4001},
		{1, 4001},
	}}}

	for len(stack) > 0 {
		item := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// Last step.
		switch item.action {
		case ActionAccept:
			count += item.intervals.Count()
		case ActionReject:
		}

		for _, step := range rs.Workflows[string(item.action)] {
			if step.IsUnconditional() {
				stack = append(stack, queueItem{step.Action, item.intervals})
				break
			}

			l := item.intervals
			r := item.intervals
			switch step.Comparison {
			case ComparisonLessThan:
				l[step.Category].Max = step.Value
				r[step.Category].Min = step.Value
			case ComparisonGreaterThan:
				l[step.Category].Min = step.Value + 1
				r[step.Category].Max = step.Value + 1
			}

			stack = append(stack, queueItem{step.Action, l})
			item.intervals = r
		}
	}

	return count
}

func part2(input string) int {
	rs := parseInput(input)
	return allAcceptedCombinations(rs)
}
