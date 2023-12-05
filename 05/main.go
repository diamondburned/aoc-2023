package main

import (
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	aocutil.Run(part1, part2)
}

const (
	chunkSeeds = iota
	chunkSeedToSoil
	chunkSoilToFertilizer
	chunkFertilizerToWater
	chunkWaterToLight
	chunkLightToTemperature
	chunkTemperatureToHumidity
	chunkHumidityToLocation
	chunkCount
)

type Almanac struct {
	Seeds                 [][2]int
	SeedToSoil            RangeMap
	SoilToFertilizer      RangeMap
	FertilizerToWater     RangeMap
	WaterToLight          RangeMap
	LightToTemperature    RangeMap
	TemperatureToHumidity RangeMap
	HumidityToLocation    RangeMap
}

func (a Almanac) RangeMaps() []RangeMap {
	return []RangeMap{
		a.SeedToSoil,
		a.SoilToFertilizer,
		a.FertilizerToWater,
		a.WaterToLight,
		a.LightToTemperature,
		a.TemperatureToHumidity,
		a.HumidityToLocation,
	}
}

type Range struct {
	Start  int
	Length int
}

// String returns a string representation of the range.
func (r Range) String() string {
	return fmt.Sprintf("[%d, %d)", r.Start, r.Start+r.Length)
}

// Intersect returns the intersection of the two ranges.
func (r Range) Intersect(other Range) Range {
	start := max(r.Start, other.Start)
	end := min(r.Start+r.Length, other.Start+other.Length)
	if start >= end {
		return Range{}
	}
	return Range{
		Start:  start,
		Length: end - start,
	}
}

// Except returns the ranges that are not in the other range.
// It may return a maximum of two ranges, one on the left and one on the right.
func (r Range) Except(other Range) []Range {
	if other.Start <= r.Start && other.Start+other.Length >= r.Start+r.Length {
		return nil
	}

	var ranges []Range
	if other.Start > r.Start {
		ranges = append(ranges, Range{
			Start:  r.Start,
			Length: other.Start - r.Start,
		})
	}
	if other.Start+other.Length < r.Start+r.Length {
		ranges = append(ranges, Range{
			Start:  other.Start + other.Length,
			Length: r.Start + r.Length - (other.Start + other.Length),
		})
	}

	return ranges
}

// RangeMap is a map of ranges. It is stored as a sorted slice of RangeMapItems.
// This slice can be binary searched to find the correct range.
type RangeMap []RangeMapItem

func (rm RangeMap) SortBySource() {
	slices.SortFunc(rm, func(a, b RangeMapItem) int {
		return a.SourceStart - b.SourceStart
	})
}

func (rm RangeMap) SortByDestination() {
	slices.SortFunc(rm, func(a, b RangeMapItem) int {
		return a.DestinationStart - b.DestinationStart
	})
}

// MapSourceRange searches for the given range in the source range.
// The destination range in valueRange is ignored.
func (rm RangeMap) MapSourceRange(v Range) []Range {
	inputs := []Range{v}
	var ranges []Range

queueLoop:
	for len(inputs) > 0 {
		v := inputs[0]
		inputs = inputs[1:]

		for _, r := range rm {
			// Try to calculate the intersection of the source range and the
			// given range.
			intersection := r.SourceRange().Intersect(v)
			if intersection.Length == 0 {
				continue
			}

			// Some intersection was found. Calculate the destination range.
			dstStart := r.SourceToDestination(intersection.Start)
			dstEnd := r.SourceToDestination(intersection.Start + intersection.Length - 1)
			dstLen := dstEnd - dstStart + 1

			ranges = append(ranges, Range{
				Start:  dstStart,
				Length: dstLen,
			})

			// If the intersection is the same as the source range, we can stop
			// searching.
			if intersection.Start == v.Start && intersection.Length == v.Length {
				continue queueLoop
			}

			// Otherwise, we need to search for the remaining parts.
			inputs = append(inputs, v.Except(intersection)...)

			// We can stop searching for more ranges in the current range.
			continue queueLoop
		}

		// No intersection was found. The given range is not mapped.
		ranges = append(ranges, v)
	}

	return ranges
}

type RangeMapItem struct {
	// SourceStart is the starting value of the source range.
	SourceStart int
	// DestinationStart is the starting value of the destination range.
	DestinationStart int
	// Length is the length of the source and destination.
	Length int
}

func (r RangeMapItem) SourceRange() Range {
	return Range{r.SourceStart, r.Length}
}

func (r RangeMapItem) DestinationRange() Range {
	return Range{r.DestinationStart, r.Length}
}

// SourceIsInRange returns true if the given value is in the source range.
func (r RangeMapItem) SourceIsInRange(value int) bool {
	return value >= r.SourceStart && value < r.SourceStart+r.Length
}

// SourceToDestination converts the given value from the source range to the
// destination range.
func (r RangeMapItem) SourceToDestination(value int) int {
	if !r.SourceIsInRange(value) {
		panic("value is not in range")
	}
	return r.DestinationStart + (value - r.SourceStart)
}

func parseAlmanac(input string, partNum int) Almanac {
	chunks := strings.Split(input, "\n\n")[:chunkCount]
	chunks = aocutil.Filter(chunks, func(s string) bool { return s != "" })
	for i, chunk := range chunks {
		if i == 0 {
			_, value, _ := strings.Cut(chunk, ": ")
			chunks[i] = value
		} else {
			_, value, _ := strings.Cut(chunk, "\n")
			chunks[i] = value
		}
	}

	var almanac Almanac

	seedsLine := aocutil.Atois[int](strings.Fields(chunks[chunkSeeds]))
	switch partNum {
	case 1:
		almanac.Seeds = make([][2]int, len(seedsLine))
		for i, seed := range seedsLine {
			almanac.Seeds[i] = [2]int{seed, seed}
		}
	case 2:
		almanac.Seeds = make([][2]int, len(seedsLine)/2)
		for i := 0; i < len(seedsLine); i += 2 {
			start := seedsLine[i]
			almanac.Seeds[i/2] = [2]int{start, start + seedsLine[i+1]}
		}
	}

	parseMap := func(chunk int) RangeMap {
		lines := aocutil.SplitLines(chunks[chunk])
		rmap := make(RangeMap, len(lines))
		for i, line := range lines {
			nums := aocutil.Atois[int](strings.Fields(line))
			rmap[i] = RangeMapItem{
				SourceStart:      nums[1],
				DestinationStart: nums[0],
				Length:           nums[2],
			}
		}
		rmap.SortBySource()
		return rmap
	}

	almanac.SeedToSoil = parseMap(chunkSeedToSoil)
	almanac.SoilToFertilizer = parseMap(chunkSoilToFertilizer)
	almanac.FertilizerToWater = parseMap(chunkFertilizerToWater)
	almanac.WaterToLight = parseMap(chunkWaterToLight)
	almanac.LightToTemperature = parseMap(chunkLightToTemperature)
	almanac.TemperatureToHumidity = parseMap(chunkTemperatureToHumidity)
	almanac.HumidityToLocation = parseMap(chunkHumidityToLocation)

	return almanac
}

func part1(input string) int {
	almanac := parseAlmanac(input, 1)
	minDist := math.MaxInt

	for _, seed := range almanac.Seeds {
		n := seed[0]
		for _, rangeMaps := range almanac.RangeMaps() {
			r := rangeMaps.MapSourceRange(Range{n, 1})
			n = r[0].Start
		}
		minDist = min(minDist, n)
	}

	return minDist
}

func part2(input string) int {
	almanac := parseAlmanac(input, 2)

	currentRanges := make([]Range, len(almanac.Seeds))
	for i, seedRange := range almanac.Seeds {
		currentRanges[i] = Range{
			Start:  seedRange[0],
			Length: seedRange[1] - seedRange[0],
		}
	}

	for _, rangeMap := range almanac.RangeMaps() {
		nextRanges := make([]Range, 0, len(currentRanges))
		for _, currentRange := range currentRanges {
			nextRange := rangeMap.MapSourceRange(currentRange)
			if nextRange == nil {
				nextRanges = []Range{currentRange}
			}
			nextRanges = append(nextRanges, nextRange...)
		}
		currentRanges = nextRanges
	}

	minDist := math.MaxInt
	for _, r := range currentRanges {
		minDist = min(minDist, r.Start)
	}

	return minDist
}
