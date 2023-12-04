package main

import (
	"strconv"
	"strings"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	aocutil.ParseAndRun(parseCards, part1, part2)
}

type Card struct {
	WinningNumbers []int
	MyNumbers      []int
}

func parseCards(input string) []Card {
	lines := aocutil.SplitLines(input)
	cards := make([]Card, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			continue
		}

		_, tail, _ := strings.Cut(line, ": ")
		left, right, _ := strings.Cut(tail, " | ")

		var card Card
		for _, word := range strings.Fields(left) {
			n := aocutil.E2(strconv.Atoi(word))
			card.WinningNumbers = append(card.WinningNumbers, n)
		}
		for _, word := range strings.Fields(right) {
			n := aocutil.E2(strconv.Atoi(word))
			card.MyNumbers = append(card.MyNumbers, n)
		}

		cards = append(cards, card)
	}

	return cards
}

func part1(input []Card) int {
	var sum int

	for _, card := range input {
		winningSet := aocutil.NewSetFromSlice(card.WinningNumbers)

		var pts int
		for _, n := range card.MyNumbers {
			if !winningSet.Has(n) {
				continue
			}
			if pts == 0 {
				pts = 1
			} else {
				pts *= 2
			}
		}

		sum += pts
	}

	return sum
}

func cardAt(cards []Card, number int) *Card {
	if number < 1 || number > len(cards) {
		return nil
	}
	return &cards[number-1]
}

func part2(input []Card) int {
	var sum int

	type WonCard struct {
		Card
		Times int
	}

	wonCards := make([]WonCard, len(input))
	for i, cards := range input {
		wonCards[i] = WonCard{Card: cards, Times: 1}
	}

	for i, card := range wonCards {
		winningSet := aocutil.NewSetFromSlice(card.WinningNumbers)

		var won int
		for _, n := range card.MyNumbers {
			if winningSet.Has(n) {
				won++
			}
		}

		for x := i + 1; x < i+1+won; x++ {
			wonCards[x].Times += wonCards[i].Times
		}
	}

	for _, won := range wonCards {
		sum += won.Times
	}

	return sum
}
