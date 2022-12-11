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
	blocks := strings.Split(input, "\n\n")
	for _, block := range blocks {
		if block == "" {
			continue
		}
		monkeys = append(monkeys, MustParseMonkey(block))
	}

	// OH MY GOD I AM FUCKING STUPID. part1(monkeys) MUTATES THE SLICE. THAT WAS
	// WHY PART 2 IS WRONG. WHAT A DUMB MISTAKE!! CONST WHEN?!?!?!?!
	part1(aocutil.Clone(monkeys))
	part2(aocutil.Clone(monkeys))
}

func part1(monkeys []Monkey) {
	inspectTimes := make(map[int]int)

	for round := 0; round < 20; round++ {
		fmt.Print("round ", round+1, ": ")
		fmt.Println()

		for i := range monkeys {
			monkey := &monkeys[i]
			for _, item := range monkey.StartingItems {
				// Log inspect time.
				inspectTimes[i]++

				// Apply worry level to item.
				item = monkey.Operation.Do(item)
				// Monkey is now bored, so div 3.
				item = item / 3
				// Check for next monkey.
				nextMonkeyIx := monkey.Test.Do(int(item))
				nextMonkey := &monkeys[nextMonkeyIx]
				// Throw to next monkey.
				nextMonkey.StartingItems = append(nextMonkey.StartingItems, item)
				// Remove item from current monkey.
				monkey.StartingItems = monkey.StartingItems[1:]
			}
		}

		for i, monkey := range monkeys {
			fmt.Print("  monkey ", i, ": ", monkey.StartingItems)
			fmt.Println()
		}
	}

	fmt.Println("inspect times:", inspectTimes)
	inspectTimesList := aocutil.MapPairs(inspectTimes)
	sort.Slice(inspectTimesList, func(i, j int) bool {
		return inspectTimesList[i].V > inspectTimesList[j].V
	})
	fmt.Println("monkey business:", inspectTimesList[0].V*inspectTimesList[1].V)
}

func part2(monkeys []Monkey) {
	inspectTimes := make(map[int]int)

	careRounds := map[int]bool{
		1:     true,
		20:    true,
		1000:  true,
		2000:  true,
		3000:  true,
		4000:  true,
		5000:  true,
		6000:  true,
		7000:  true,
		8000:  true,
		9000:  true,
		10000: true,
	}

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

	for round := 0; round < 10000; round++ {
		for i := range monkeys {
			monkey := &monkeys[i]
			for _, item := range monkey.StartingItems {
				// Log inspect time.
				inspectTimes[i]++

				// Apply worry level to item.
				item = monkey.Operation.Do(item)
				// We're not dividing by 3 anymore.
				// Minimize the item number so that it doesn't overflow.
				item = item % product

				// Check for next monkey.
				nextMonkeyIx := monkey.Test.Do(int(item))
				nextMonkey := &monkeys[nextMonkeyIx]
				// Throw to next monkey.
				nextMonkey.StartingItems = append(nextMonkey.StartingItems, item)
				// Remove item from current monkey.
				monkey.StartingItems = monkey.StartingItems[1:]
			}
		}

		if careRounds[round+1] {
			fmt.Print("round ", round+1, ": ")
			fmt.Println()
			fmt.Println("  inspect times:", inspectTimes)
			for i, monkey := range monkeys {
				fmt.Print("  monkey ", i, ": ", monkey.StartingItems)
				fmt.Println()
			}
		}
	}

	fmt.Println("inspect times:", inspectTimes)
	inspectTimesList := aocutil.MapPairs(inspectTimes)
	sort.Slice(inspectTimesList, func(i, j int) bool {
		return inspectTimesList[i].V > inspectTimesList[j].V
	})
	fmt.Println("monkey business:", inspectTimesList[0].V*inspectTimesList[1].V)
}

type Operator uint8

const (
	_ Operator = iota
	AddOp
	MulOp
)

// Monkey represents a monkey.
type Monkey struct {
	// StartingItems is a list of items (worry levels).
	StartingItems []int
	Operation     Operation
	// Test shows how the monkey uses the worry level to decide which monkey to
	// throw next.
	Test MonkeyTest
}

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

type Operation struct {
	// Op is the operator used in the Operation function.
	Op Operator
	// Rhs is the right hand side of the operation. If null, then the worry
	// level is powered by 2.
	Rhs *int
}

func (op Operation) Do(old int) int {
	new := old
	if op.Rhs != nil {
		new = *op.Rhs
	}

	switch op.Op {
	case AddOp:
		new = old + new
	case MulOp:
		new = old * new
	default:
		aocutil.Assertf(false, "unknown operation %q", op.Op)
	}

	return new
}

type MonkeyTest struct {
	DivisibleBy   int
	MonkeyIfTrue  int
	MonkeyIfFalse int
}

// Do tests the number n and returns the next monkey to throw to.
func (t MonkeyTest) Do(n int) int {
	if n%t.DivisibleBy == 0 {
		return t.MonkeyIfTrue
	}
	return t.MonkeyIfFalse
}
