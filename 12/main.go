package main

import (
	"log"
	"slices"
	"strings"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	aocutil.Run(part1, part2)
}

// SpringCondition describes the operating condition of a spring.
type SpringCondition = byte

const (
	Operational SpringCondition = '.'
	Damaged     SpringCondition = '#'
	Unknown     SpringCondition = '?'
)

// SpringsConditions describes the conditions of a row of springs.
type SpringsConditions string

// SpringsRecord describes a row of springs.
type SpringsRecord struct {
	Conditions SpringsConditions
	// DamagedRuns describes runs of damaged springs. A group of [1, 1, 3]
	// may look like '#....#.###'.
	DamagedRuns []int
}

func parseInput(input string) []SpringsRecord {
	lines := aocutil.SplitLines(input)
	conditions := make([]SpringsRecord, 0, len(lines))
	for _, line := range lines {
		s1, s2, _ := strings.Cut(line, " ")
		conditions = append(conditions, SpringsRecord{
			Conditions:  SpringsConditions(s1),
			DamagedRuns: aocutil.Atois[int](strings.Split(s2, ",")),
		})
	}
	return conditions
}

// Unfold returns a copy of the SpringRow with the conditions repeated n times.
func (r SpringsRecord) Unfold(repeats int) SpringsRecord {
	new := r
	new.Conditions = SpringsConditions(strings.Repeat(string(r.Conditions)+"?", repeats))
	new.Conditions = new.Conditions[:len(new.Conditions)-1]
	new.DamagedRuns = make([]int, len(r.DamagedRuns)*repeats)
	for i := 0; i < repeats; i++ {
		copy(new.DamagedRuns[i*len(r.DamagedRuns):], r.DamagedRuns)
	}
	return new
}

func totalValid(springs SpringsRecord) int {
	type countFunc func(springs SpringsConditions, n int, damagedRuns []int) int
	var count countFunc
	var countActual countFunc

	type cacheKey struct {
		Remaining   SpringsConditions
		DamagedRuns []int
	}

	cache := aocutil.NewAnyMap[cacheKey, int]()

	count = func(springs SpringsConditions, n int, damagedRuns []int) int {
		var count int

		if !aocutil.IsSilent() {
			log.Printf("%q %v", springs, damagedRuns)
			oldPrefix := log.Prefix()
			log.SetPrefix(oldPrefix + "⎸")
			defer func() {
				log.Printf("%d <- %q", count, springs)
				log.SetPrefix(oldPrefix)
			}()
		}

		var ok bool
		key := cacheKey{springs[n:], damagedRuns}

		if count, ok = cache.Get(key); ok {
			return count
		}

		count = countActual(springs, n, damagedRuns)
		cache.Set(key, count)

		return count
	}

	countActual = func(springs SpringsConditions, n int, damagedRuns []int) int {
		// Base case.
		if n == len(springs) {
			switch {
			case len(damagedRuns) == 0:
			case len(damagedRuns) == 1 && damagedRuns[0] == 0:
			default:
				// We're still expecting more damaged runs.
				// This is not valid.
				return 0
			}

			// We have nothing left. This is fine.
			return 1
		}

		switch springs[n] {
		case Operational:
			if len(damagedRuns) > 0 {
				if damagedRuns[0] == 0 {
					damagedRuns = damagedRuns[1:]
				} else if n != 0 && springs[n-1] == Damaged {
					// We're still expecting more damaged runs.
					// This is not valid.
					// ⎸⎸⎸⎸⎸⎸⎸⎸".###.#??????" [2 1]
					return 0
				}
			}

			// We're still expecting more operational runs.
			return count(springs, n+1, damagedRuns)

		case Damaged:
			if len(damagedRuns) == 0 || damagedRuns[0] == 0 {
				// We're not expecting any more damaged runs.
				// This is not valid.
				return 0
			}

			damagedRuns[0]--
			return count(springs, n+1, damagedRuns)

		case Unknown:
			// Expect either operational or damaged.
			springs = aocutil.ReplaceStringIndex(springs, n, SpringsConditions(Operational))
			a := count(springs, n, slices.Clone(damagedRuns))
			springs = aocutil.ReplaceStringIndex(springs, n, SpringsConditions(Damaged))
			b := count(springs, n, slices.Clone(damagedRuns))
			return a + b

		default:
			panic("unreachable")
		}
	}

	log.Printf("totalValid(%s, %v)\n", springs.Conditions, springs.DamagedRuns)
	valid := count(springs.Conditions, 0, springs.DamagedRuns)
	log.Printf("  = %d\n", valid)
	return valid
}

func part1(input string) int {
	var total int
	for _, row := range parseInput(input) {
		total += totalValid(row)
	}
	return total
}

func part2(input string) int {
	var total int
	for _, row := range parseInput(input) {
		row = row.Unfold(5)
		total += totalValid(row)
	}
	return total
}
