package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	var monkeys []Monkey
	input := aocutil.InputString()
	for _, block := range aocutil.SplitBlocks(input) {
		monkeys = append(monkeys, MustParseMonkey(block))
	}

	// OH MY GOD I AM FUCKING STUPID. part1(monkeys) MUTATES THE SLICE. THAT WAS
	// WHY PART 2 IS WRONG. WHAT A DUMB MISTAKE!! CONST WHEN?!?!?!?!
	part1(aocutil.Clone(monkeys))
	part2(aocutil.Clone(monkeys))
}

func part1(monkeys Monkeys) {
	inspectTimes := make(map[int]int)

	for round := 1; round <= 20; round++ {
		for i := range monkeys {
			monkey := &monkeys[i]
			for _, item := range monkey.StartingItems {
				// Log inspect time.
				inspectTimes[i]++
				// Apply worry level to item.
				item = monkey.Operation.Do(item)
				// Monkey is now bored, so div 3.
				item = item / 3
				// Pass the item.
				monkeys.Pass(monkey, item)
			}
		}

		printState(round+1, monkeys, inspectTimes)
	}

	printMonkeyBusiness(inspectTimes)
}

func part2(monkeys Monkeys) {
	inspectTimes := make(map[int]int)

	// Bro this shit is tough as FUCK. My initial solution was to do
	// item %= divisibleBy to limit the item number to the range of the
	// monkey test. However, this doesn't work because we're passing the
	// item to the next monkey, and the next monkey might have a different
	// divisibleBy.
	//
	// There's apparently a theorem where we can do `n % (a*b)` which will yield
	// us the same result as `n % a % b`. We can extend this to our monkeys:
	// assuming we have n monkeys, we can modulo the item by the product of all
	// the monkeys' divisibleBy. This will yield the same result as if we
	// eventually modulo'd the item by each monkey's divisibleBy.
	var product int = 1
	for _, monkey := range monkeys {
		product *= monkey.Test.DivisibleBy
	}

	for round := 1; round <= 10000; round++ {
		for i := range monkeys {
			monkey := &monkeys[i]
			for _, item := range monkey.StartingItems {
				// Log inspect time.
				inspectTimes[i]++
				// Apply worry level to item.
				// Minimize the item number so that it doesn't overflow.
				item = monkey.Operation.Do(item) % product
				// Pass the item.
				monkeys.Pass(monkey, item)
			}
		}

		if round == 1 || round == 20 || (round%1000) == 0 {
			printState(round, monkeys, inspectTimes)
		}
	}

	printMonkeyBusiness(inspectTimes)
}

func printState(round int, monkeys Monkeys, inspectTimes map[int]int) {
	fmt.Print("round ", round, ": ")
	fmt.Println()
	fmt.Println("  inspect times:", inspectTimes)
	for i, monkey := range monkeys {
		fmt.Print("  monkey ", i, ": ", monkey.StartingItems)
		fmt.Println()
	}
}

func printMonkeyBusiness(inspectTimes map[int]int) {
	fmt.Println("inspect times:", inspectTimes)
	inspectTimesList := aocutil.MapPairs(inspectTimes)
	sort.Slice(inspectTimesList, func(i, j int) bool {
		return inspectTimesList[i].V > inspectTimesList[j].V
	})
	fmt.Println("monkey business:", inspectTimesList[0].V*inspectTimesList[1].V)
}

// Operator is a mathematical operator.
type Operator uint8

const (
	_ Operator = iota
	AddOp
	MulOp
)

// Monkeys represents a bunch of monkeys.
type Monkeys []Monkey

// Pass passes the item to the next monkey based on the item's value. It returns
// the next monkey's index.
func (ms Monkeys) Pass(monkey *Monkey, item int) int {
	// Check for next monkey.
	nextMonkeyIx := monkey.Test.NextMonkey(item)
	nextMonkey := &ms[nextMonkeyIx]
	// Throw to next monkey.
	nextMonkey.StartingItems = append(nextMonkey.StartingItems, item)
	// Remove item from current monkey.
	monkey.StartingItems = monkey.StartingItems[1:]
	return nextMonkeyIx
}

// Monkey represents a monkey.
type Monkey struct {
	// StartingItems is a list of items (worry levels).
	StartingItems []int
	Operation     Operation
	// Test shows how the monkey uses the worry level to decide which monkey to
	// throw next.
	Test MonkeyTest
}

// MustParseMonkey parses a monkey from the given text block.
func MustParseMonkey(input string) Monkey {
	var m Monkey

	lines := strings.Split(input, "\n")
	aocutil.Assertf(len(lines) >= 6, "not enough lines")

	aocutil.Assertf(strings.HasPrefix(lines[0], "Monkey "), "input must contain 'Monkey'")

	startingItems := strings.TrimPrefix(lines[1], "  Starting items: ")
	for _, s := range strings.Split(startingItems, ", ") {
		i := aocutil.Atoi[int](s)
		m.StartingItems = append(m.StartingItems, i)
	}

	operation := strings.TrimPrefix(lines[2], "  Operation: new =")
	operationParams := strings.Fields(operation)

	if operationParams[2] == "old" {
		m.Operation.Rhs = nil
	} else {
		v := aocutil.Atoi[int](operationParams[2])
		m.Operation.Rhs = &v
	}

	switch operationParams[1] {
	case "+":
		m.Operation.Op = AddOp
	case "*":
		m.Operation.Op = MulOp
	default:
		aocutil.Assertf(false, "unknown operation %q", operationParams[0])
	}

	var testDivisibleBy int
	aocutil.Sscanf(lines[3], "  Test: divisible by %d", &testDivisibleBy)
	m.Test.DivisibleBy = testDivisibleBy

	var testMonkeyIfTrue, testMonkeyIfFalse int
	aocutil.Sscanf(lines[4], "    If true: throw to monkey %d", &testMonkeyIfTrue)
	aocutil.Sscanf(lines[5], "    If false: throw to monkey %d", &testMonkeyIfFalse)
	m.Test.MonkeyIfTrue = testMonkeyIfTrue
	m.Test.MonkeyIfFalse = testMonkeyIfFalse

	return m
}

// Operation represents a basic mathematical operation.
type Operation struct {
	// Op is the operator used in the Operation function.
	Op Operator
	// Rhs is the right hand side of the operation. If null, then the worry
	// level is powered by 2.
	Rhs *int
}

// Do applies the operation to the given number.
func (op Operation) Do(n int) int {
	rhs := n
	if op.Rhs != nil {
		rhs = *op.Rhs
	}

	switch op.Op {
	case AddOp:
		return n + rhs
	case MulOp:
		return n * rhs
	default:
		aocutil.Assertf(false, "unknown operation %q", op.Op)
		return 0
	}
}

// MonkeyTest represents a test that a monkey performs on a worry level.
type MonkeyTest struct {
	DivisibleBy   int
	MonkeyIfTrue  int
	MonkeyIfFalse int
}

// NextMonkey tests the number n and returns the next monkey to throw to.
func (t MonkeyTest) NextMonkey(n int) int {
	if n%t.DivisibleBy == 0 {
		return t.MonkeyIfTrue
	}
	return t.MonkeyIfFalse
}
