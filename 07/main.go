package main

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	aocutil.Run(part1, part2)
}

// Card is a card in a deck.
type Card struct {
	Label    byte
	WasJoker bool
}

func (c Card) String() string {
	if !c.WasJoker {
		return " " + string(c.Label)
	}
	return fmt.Sprintf(">%v", string(c.Label))
}

// Less returns true if c is less than other.
func (c Card) Compare(other Card) int {
	if c.WasJoker == other.WasJoker {
		return c.compare(other)
	}
	if c.WasJoker {
		c.Label = 'J'
	}
	if other.WasJoker {
		other.Label = 'J'
	}
	return c.compare(other)
}

func (c Card) compare(other Card) int {
	v1 := c.toValue()
	v2 := other.toValue()
	if v1 < v2 {
		return -1
	}
	if v1 > v2 {
		return 1
	}
	return 0
}

func (c Card) toValue() int {
	if '2' <= c.Label && c.Label <= '9' {
		return int(c.Label - '0')
	}
	switch c.Label {
	case 'T':
		return 10
	case 'Q':
		return 11
	case 'K':
		return 12
	case 'A':
		return 13
	case 'J':
		return 0
	default:
		return -1
	}
}

type HandType uint

const (
	HighCard HandType = iota
	OnePair
	TwoPair
	ThreeOfAKind
	FullHouse
	FourOfAKind
	FiveOfAKind
)

// Hand is a hand of cards.
type Hand [5]Card

func (h Hand) freqMap() map[byte]int {
	freq := make(map[byte]int, 5)
	for _, card := range h {
		freq[card.Label]++
	}
	return freq
}

func freqMapHand(freq map[byte]int) HandType {
	if len(freq) == 1 {
		return FiveOfAKind
	}
	if len(freq) == 2 {
		for _, v := range freq {
			if v == 4 {
				return FourOfAKind
			}
			if v == 3 {
				return FullHouse
			}
		}
	}
	if len(freq) == 3 {
		for _, v := range freq {
			if v == 3 {
				return ThreeOfAKind
			}
			if v == 2 {
				return TwoPair
			}
		}
	}
	if len(freq) == 4 {
		return OnePair
	}
	return HighCard
}

func (h Hand) Jokerify() Hand {
	freq := make(map[byte]int, 5)
	for _, card := range h {
		freq[card.Label]++
	}

	for n := freq['J']; n != 0; n-- {
		maxCard := freqMax(freq)
		log.Printf("%s: decided that %c is the max card", h, maxCard)
		jokerPos := slices.Index(h[:], Card{Label: 'J'})
		h[jokerPos] = Card{maxCard, true}
	}

	return h
}

func freqMax(freq map[byte]int) byte {
	var maxFreq int
	var maxCard byte
	for c, f := range freq {
		if c == 'J' {
			continue
		}
		if f < maxFreq {
			continue
		}
		if f == maxFreq {
			// If there's a tie, pick the higher card.
			v1 := Card{c, false}.toValue()
			v2 := Card{maxCard, false}.toValue()
			if v1 < v2 {
				continue
			}
		}
		maxFreq = f
		maxCard = c
	}
	return maxCard
}

// Unjokerify returns a hand with previously jokered cards unjokered.
func (h Hand) Unjokerify() Hand {
	for i, card := range h {
		if card.WasJoker {
			h[i] = Card{'J', false}
		}
	}
	return h
}

// Sorted returns a sorted copy of the hand.
func (h Hand) Sorted() Hand {
	cards := []Card(h[:])
	slices.SortFunc(cards, Card.Compare)
	return Hand(cards)
}

func (h Hand) Compare(other Hand) int {
	type1 := freqMapHand(h.freqMap())
	type2 := freqMapHand(other.freqMap())
	if type1 != type2 {
		return int(type1 - type2)
	}

	h1 := h.Unjokerify()
	h2 := other.Unjokerify()

	// Try to break tie by comparing cards.
	for i := 0; i < 5; i++ {
		c1 := h1[i]
		c2 := h2[i]
		cx := c1.Compare(c2)
		if cx != 0 {
			return cx
		}
	}

	return 0
}

type BiddingHand struct {
	Hand Hand
	Bid  int
}

func parseInput(input string) []BiddingHand {
	lines := aocutil.SplitLines(input)
	lines = aocutil.FilterEmptyStrings(lines)

	hands := make([]BiddingHand, len(lines))
	for i, line := range lines {
		k, v, _ := strings.Cut(line, " ")
		for j, c := range k {
			hands[i].Hand[j] = Card{byte(c), false}
		}
		hands[i].Bid = aocutil.Atoi[int](v)
	}
	return hands
}

func part1(input string) int {
	hands := parseInput(input)
	slices.SortFunc(hands, func(a, b BiddingHand) int {
		return a.Hand.Compare(b.Hand)
	})

	for i := range hands {
		rank := i + 1
		hands[i].Bid *= rank
	}

	var winnings int
	for _, hand := range hands {
		winnings += hand.Bid
	}

	return winnings
}

func part2(input string) int {
	hands := parseInput(input)
	log.Printf(" pre-jokerify: %v", hands)
	for i, bidding := range hands {
		hands[i].Hand = bidding.Hand.Jokerify()
	}
	log.Printf("post-jokerify: %v", hands)

	slices.SortFunc(hands, func(a, b BiddingHand) int {
		return a.Hand.Compare(b.Hand)
	})

	for i := len(hands) - 1; i >= 0; i-- {
		rank := i + 1
		hands[i].Bid *= rank
		log.Printf("rank %d: %v", rank, hands[i])
	}

	var winnings int
	for _, hand := range hands {
		winnings += hand.Bid
		if winnings < 0 {
			panic("winnings < 0")
		}
	}

	return winnings
}
