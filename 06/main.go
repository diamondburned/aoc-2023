package main

import (
	"strings"

	"libdb.so/aoc-2023/aocutil"
)

func main() {
	aocutil.Run(part1, part2)
}

type RaceRecord struct {
	Time     int
	Distance int
}

func parseRecords(input string, partNum int) []RaceRecord {
	lines := strings.Split(input, "\n")
	lineTime := lines[0]
	lineDistance := lines[1]

	_, valueTime, _ := strings.Cut(lineTime, ": ")
	_, valueDistance, _ := strings.Cut(lineDistance, ": ")

	if partNum == 2 {
		// bad kerning lmfao
		valueTime = strings.ReplaceAll(valueTime, " ", "")
		valueDistance = strings.ReplaceAll(valueDistance, " ", "")
	}

	timeValues := aocutil.Atois[int](strings.Fields(valueTime))
	distanceValues := aocutil.Atois[int](strings.Fields(valueDistance))

	records := make([]RaceRecord, len(timeValues))
	for i := range records {
		records[i].Time = timeValues[i]
		records[i].Distance = distanceValues[i]
	}

	return records
}

func beatsRace(race RaceRecord, holdingMs int) bool {
	time := race.Time
	time -= holdingMs // penalty for holding
	if time <= 0 {
		return false // L rizz
	}

	distance := time * holdingMs
	return race.Distance < distance
}

func part1(input string) int {
	records := parseRecords(input, 1)
	totalMul := 1
	for _, record := range records {
		var totalWays int
		for i := 1; i < record.Time; i++ {
			if beatsRace(record, i) {
				totalWays++
			}
		}
		totalMul *= totalWays
	}

	return totalMul
}

func part2(input string) int {
	records := parseRecords(input, 2)
	totalMul := 1
	for _, record := range records {
		var totalWays int
		for i := 1; i < record.Time; i++ {
			if beatsRace(record, i) {
				totalWays++
			}
		}
		totalMul *= totalWays
	}
	return totalMul
}
