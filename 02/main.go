package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	input := aocutil.ReadStdin()

	part1(input)
	part2(input)
}

type Game struct {
	Pulls []Pull
}

type Pull struct {
	Blue  int
	Red   int
	Green int
}

var gameLine = regexp.MustCompile(`Game \d*: (.*;? *)`)

func parseGames(input string) []Game {
	lines := aocutil.SplitLines(input)
	games := make([]Game, 0, len(lines))

	for _, line := range lines {
		match := gameLine.FindStringSubmatch(line)
		if match == nil {
			continue
		}

		parts := strings.Split(match[1], "; ")
		var game Game

		for _, part := range parts {
			balls := strings.Split(part, ", ")
			var pull Pull

			for _, ball := range balls {
				value, name, _ := strings.Cut(ball, " ")
				valueInt, _ := strconv.Atoi(value)

				switch name {
				case "blue":
					pull.Blue = valueInt
				case "red":
					pull.Red = valueInt
				case "green":
					pull.Green = valueInt
				default:
					log.Panicf("unknown ball color %q", name)
				}
			}

			game.Pulls = append(game.Pulls, pull)
		}

		games = append(games, game)
	}

	return games
}

func part1(input string) {
	games := parseGames(input)

	const wantRed = 12
	const wantGreen = 13
	const wantBlue = 14

	var idSum int
searchLoop:
	for i, game := range games {
		for _, pull := range game.Pulls {
			if pull.Blue > wantBlue || pull.Red > wantRed || pull.Green > wantGreen {
				continue searchLoop
			}
		}
		idSum += i + 1
	}

	fmt.Println(idSum)
}

func part2(input string) {
	games := parseGames(input)
	var powerSum int
	for _, game := range games {
		maxBlue := aocutil.Maxs(
			aocutil.Map(game.Pulls, func(p Pull) int { return p.Blue })...)
		maxRed := aocutil.Maxs(
			aocutil.Map(game.Pulls, func(p Pull) int { return p.Red })...)
		maxGreen := aocutil.Maxs(
			aocutil.Map(game.Pulls, func(p Pull) int { return p.Green })...)
		power := maxBlue * maxRed * maxGreen
		powerSum += power
	}

	fmt.Println(powerSum)
}
