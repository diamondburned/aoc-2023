package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	input, _ := os.ReadFile("input")

	var overlaps int

	scanner := aocutil.NewBytesScanner(input)
	scanner.SetSplitter(bufio.ScanLines)
	for scanner.Next() {
		line := scanner.Token()

		elfPair := MustParseSectionRangePair(line)
		if elfPair[0].Overlaps(elfPair[1]) || elfPair[1].Overlaps(elfPair[0]) {
			overlaps++
		}
	}

	fmt.Println(overlaps)
}

type SectionRange struct{ lo, hi uint8 }

func MustParseSectionRange(input string) SectionRange {
	parts := aocutil.SplitN(input, "-", 2)
	i1 := aocutil.E2(strconv.ParseUint(parts[0], 10, 8))
	i2 := aocutil.E2(strconv.ParseUint(parts[1], 10, 8))
	return SectionRange{uint8(i1), uint8(i2)}
}

func MustParseSectionRangePair(input string) [2]SectionRange {
	parts := aocutil.SplitN(input, ",", 2)
	return [2]SectionRange{
		MustParseSectionRange(parts[0]),
		MustParseSectionRange(parts[1]),
	}
}

// Overlaps returns true if the given range partially or fully overlaps with the
// receiver.
func (r1 SectionRange) Overlaps(r2 SectionRange) bool {
	return r1.lo <= r2.hi && r2.lo <= r1.hi
	// for i := r.lo; i <= r.hi; i++ {
	// 	if i >= r2.lo && i <= r2.hi {
	// 		return true
	// 	}
	// }
	// return false
}

// FullyOverlaps returns true if the given range is fully contained within the
// receiver.
func (r1 SectionRange) FullyOverlaps(r2 SectionRange) bool {
	return r1.lo <= r2.lo && r2.hi <= r1.hi
}
