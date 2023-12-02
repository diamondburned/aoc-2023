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

		parts := strings.Split(match[1], ";")
		var game Game
		for _, part := range parts {
			part = strings.TrimSpace(part)
			balls := strings.Split(part, ",")
			var pull Pull
			for _, ball := range balls {
				ball = strings.TrimSpace(ball)
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

	wantRed := 12
	wantGreen := 13
	wantBlue := 14

	possibleGames := make([]Game, 0, len(games))
	var idSum int
searchLoop:
	for i, game := range games {
		for _, pull := range game.Pulls {
			if pull.Blue > wantBlue || pull.Red > wantRed || pull.Green > wantGreen {
				continue searchLoop
			}
		}
		possibleGames = append(possibleGames, game)
		idSum += i + 1
	}

	fmt.Println(idSum)
}

func part2(input string) {
	games := parseGames(input)
	//
	// wantRed := 12
	// wantGreen := 13
	// wantBlue := 14

	// 	possibleGames := make([]Game, 0, len(games))
	// searchLoop:
	// 	for _, game := range games {
	// 		for _, pull := range game.Pulls {
	// 			if pull.Blue > wantBlue || pull.Red > wantRed || pull.Green > wantGreen {
	// 				continue searchLoop
	// 			}
	// 		}
	// 		possibleGames = append(possibleGames, game)
	// 	}

	type maxGame struct {
		MaxBlue  int
		MaxRed   int
		MaxGreen int
		Power    int
	}
	var powerSum int
	for i, game := range games {
		var maxGame maxGame
		for _, pull := range game.Pulls {
			if pull.Blue > maxGame.MaxBlue {
				maxGame.MaxBlue = pull.Blue
			}
			if pull.Red > maxGame.MaxRed {
				maxGame.MaxRed = pull.Red
			}
			if pull.Green > maxGame.MaxGreen {
				maxGame.MaxGreen = pull.Green
			}
		}
		maxGame.Power = maxGame.MaxBlue * maxGame.MaxRed * maxGame.MaxGreen
		log.Printf("Game %d: power = %d", i+1, maxGame.Power)
		powerSum += maxGame.Power
	}

	fmt.Println(powerSum)
}
