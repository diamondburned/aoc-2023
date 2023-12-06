package main

import (
	"log"
	"strings"

	"github.com/diamondburned/aoc-2022/aocutil"
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
	for i, record := range records {
		log.Printf("for race[%d], record=%v", i, record)
		var totalWays int
		for i := 1; i < record.Time; i++ {
			if beatsRace(record, i) {
				totalWays++
			}
			// if beatsRace(record, i) {
			// 	log.Printf("  holding for %dms beats the race", i)
			// } else {
			// 	log.Printf("  holding for %dms does not work", i)
			// }
		}
		totalMul *= totalWays
	}

	return totalMul
}

func part2(input string) int {
	records := parseRecords(input, 2)
	totalMul := 1
	for i, record := range records {
		log.Printf("for race[%d], record=%v", i, record)
		var totalWays int
		for i := 1; i < record.Time; i++ {
			if beatsRace(record, i) {
				totalWays++
			}
			// if beatsRace(record, i) {
			// 	log.Printf("  holding for %dms beats the race", i)
			// } else {
			// 	log.Printf("  holding for %dms does not work", i)
			// }
		}
		totalMul *= totalWays
	}

	return totalMul
}
