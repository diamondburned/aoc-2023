package main

import (
	"fmt"
	"log"
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

func chunkToString(chunk int) string {
	switch chunk {
	case chunkSeeds:
		return "seeds"
	case chunkSeedToSoil:
		return "seed to soil"
	case chunkSoilToFertilizer:
		return "soil to fertilizer"
	case chunkFertilizerToWater:
		return "fertilizer to water"
	case chunkWaterToLight:
		return "water to light"
	case chunkLightToTemperature:
		return "light to temperature"
	case chunkTemperatureToHumidity:
		return "temperature to humidity"
	case chunkHumidityToLocation:
		return "humidity to location"
	default:
		return "unknown"
	}
}

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
	slices.SortFunc(rm, func(a, b RangeMapItem) int { return a.SourceStart - b.SourceStart })
}

func (rm RangeMap) SortByDestination() {
	slices.SortFunc(rm, func(a, b RangeMapItem) int { return a.DestinationStart - b.DestinationStart })
}

// SearchSource searches for the given value in the source range.
func (rm RangeMap) SearchSource(value int) (RangeMapItem, bool) {
	v, _ := slices.BinarySearchFunc(rm, value, func(r RangeMapItem, t int) int {
		// log.Printf("   testing %d in %v:", t, r)
		if r.SourceIsInRange(t) {
			// log.Printf("      in range")
			return 0
		}
		if t < r.SourceStart {
			// log.Printf("      too low")
			return 1
		}
		// log.Printf("      too high")
		return -1
	})
	if v < 0 || v >= len(rm) {
		// log.Printf("   not found: %d", v)
		return RangeMapItem{}, false
	}

	r := rm[v]
	if !r.SourceIsInRange(value) {
		// log.Printf("   not found: %d", v)
		return RangeMapItem{}, false
	}

	// log.Printf("   found: %d (%v)", v, r)
	return r, true
}

// MapSource is a convenience function that searches for the given value in the
// source range and returns the destination value. If the value is not found,
// the default value is returned.
func (rm RangeMap) MapSource(value int, defaultValue int) int {
	r, ok := rm.SearchSource(value)
	if !ok {
		return defaultValue
	}
	return r.SourceToDestination(value)
}

// MapSourceRange searches for the given range in the source range.
// The destination range in valueRange is ignored.
func (rm RangeMap) MapSourceRange(v Range) []Range {
	log.Printf("    searching for %v in range map", v)

	inputs := []Range{v}
	var ranges []Range

queueLoop:
	for len(inputs) > 0 {
		v := inputs[0]
		inputs = inputs[1:]

		for _, r := range rm {
			// Try to calculate the intersection of the source range and the given
			// range.
			intersection := r.SourceRange().Intersect(v)
			if intersection.Length == 0 {
				log.Printf("      no intersection between %v and %v",
					v,
					r.SourceRange())
				continue
			}

			// Some intersection was found. Calculate the destination range.
			dstStart := r.SourceToDestination(intersection.Start)
			dstEnd := r.SourceToDestination(intersection.Start + intersection.Length - 1)
			dstLen := dstEnd - dstStart + 1

			log.Printf("      intersection between %v and %v (->%v): %d-%d -> %d-%d",
				v,
				r.SourceRange(),
				r.DestinationRange(),
				intersection.Start, intersection.Start+intersection.Length-1,
				dstStart, dstEnd)

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

			// Sort the inputs by length, so that we can process the smallest
			// ranges first.
			slices.SortFunc(inputs, func(a, b Range) int { return a.Length - b.Length })

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
	// log.Printf("   converting %d in %v:", value, r)
	dst := r.DestinationStart + (value - r.SourceStart)
	// log.Printf("      %d -> %d", value, dst)
	return dst
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
			n = rangeMaps.MapSource(n, n)
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

	for i, rangeMap := range almanac.RangeMaps() {
		log.Printf("doing chunk %s", chunkToString(i+1))
		nextRanges := make([]Range, 0, len(currentRanges))

		for _, currentRange := range currentRanges {
			log.Printf("  for range %v", currentRange)

			nextRange := rangeMap.MapSourceRange(currentRange)
			if nextRange == nil {
				nextRanges = []Range{currentRange}
			}

			nextRanges = append(nextRanges, nextRange...)
		}

		currentRanges = nextRanges
	}

	minDist := math.MaxInt
	for _, range_ := range currentRanges {
		minDist = min(minDist, range_.Start)
	}

	return minDist
}
