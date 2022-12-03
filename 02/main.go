package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	input, _ := os.ReadFile("input")
	lines := strings.Split(string(input), "\n")

	var totalScore int
	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := aocutil.FieldsN(line, 2)

		move1, ok1 := ParseShape(parts[0])
		aocutil.Assertf(ok1, "invalid move: %q", parts[0])

		outcome, ok2 := ParseOutcome(parts[1])
		aocutil.Assertf(ok2, "invalid outcome: %q", parts[1])

		move2 := move1.NextMove(outcome)
		score := int(move2) + int(move2.Score(move1))

		log.Printf(
			"made move %[1]v (%[1]d), want outcome %[2]v (%[2]d), playing with move %[3]v (%[3]d)",
			move1, outcome, move2)
		log.Printf("  => score: %d + %d = %d", move2, move2.Score(move1), score)

		totalScore += score
	}

	fmt.Println(totalScore)
}

// OutcomeScore returns the score for the winning move.
type OutcomeScore int

const (
	Loss OutcomeScore = 0
	Draw OutcomeScore = 3
	Win  OutcomeScore = 6
)

// ParseOutcome parses a string into an OutcomeScore.
func ParseOutcome(input string) (OutcomeScore, bool) {
	switch input {
	case "X":
		return Loss, true
	case "Y":
		return Draw, true
	case "Z":
		return Win, true
	default:
		return 0, false
	}
}

func (o OutcomeScore) String() string {
	switch o {
	case Loss:
		return "Loss"
	case Draw:
		return "Draw"
	case Win:
		return "Win"
	default:
		return fmt.Sprintf("OutcomeScore(%d)", o)
	}
}

// Shape is a type of rock-paper-scissors move.
type Shape uint8

const (
	_ Shape = iota
	Rock
	Paper
	Scissors
)

// ParseShape parses a string into a RPSMove.
func ParseShape(input string) (Shape, bool) {
	switch input {
	case "A", "X":
		return Rock, true
	case "B", "Y":
		return Paper, true
	case "C", "Z":
		return Scissors, true
	default:
		return 0, false
	}
}

func (rps Shape) String() string {
	switch rps {
	case Rock:
		return "Rock"
	case Paper:
		return "Paper"
	case Scissors:
		return "Scissors"
	default:
		return fmt.Sprintf("Move(%d)", rps)
	}
}

// Score returns the score for the given move.
func (rps Shape) Score(other Shape) OutcomeScore {
	switch rps {
	case Rock:
		switch other {
		case Rock:
			return Draw
		case Paper:
			return Loss
		case Scissors:
			return Win
		}
	case Paper:
		switch other {
		case Rock:
			return Win
		case Paper:
			return Draw
		case Scissors:
			return Loss
		}
	case Scissors:
		switch other {
		case Rock:
			return Loss
		case Paper:
			return Win
		case Scissors:
			return Draw
		}
	}

	return Loss
}

// NextMove returns the next move for the given move and the wanted outcome.
func (rps Shape) NextMove(outcome OutcomeScore) Shape {
	switch outcome {
	case Loss:
		return rps.LosingMove()
	case Draw:
		return rps
	case Win:
		return rps.WinningMove()
	default:
		return 0
	}
}

// LosingMove returns the losing move for the given move.
func (rps Shape) LosingMove() Shape {
	switch rps {
	case Rock:
		return Scissors
	case Paper:
		return Rock
	case Scissors:
		return Paper
	default:
		return 0
	}
}

// WinningMove returns the winning move for the given move.
func (rps Shape) WinningMove() Shape {
	switch rps {
	case Rock:
		return Paper
	case Paper:
		return Scissors
	case Scissors:
		return Rock
	default:
		return 0
	}
}

func (rps Shape) ABC() string {
	switch rps {
	case Rock:
		return "A"
	case Paper:
		return "B"
	case Scissors:
		return "C"
	default:
		return ""
	}
}

func (rps Shape) XYZ() string {
	switch rps {
	case Rock:
		return "X"
	case Paper:
		return "Y"
	case Scissors:
		return "Z"
	default:
		return ""
	}
}
