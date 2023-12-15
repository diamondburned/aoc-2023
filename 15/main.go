package main

import (
	"log"
	"slices"
	"strings"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	aocutil.Run(part1, part2)
}

func hash(str string) uint8 {
	var h uint
	for _, b := range []byte(str) {
		h += uint(b)
		h *= 17
		h %= 256
	}
	return uint8(h)
}

func part1(input string) int {
	var sum int
	for _, str := range strings.Split(input, ",") {
		v := hash(str)
		log.Printf("%s: %d", str, v)
		sum += int(v)
	}
	return sum
}

type Sequence []Instruction

type Instruction struct {
	Operation byte
	Lens
}

type Lens struct {
	Label string
	Focal int
	Box   uint8
}

const (
	DashOperation  byte = '-'
	EqualOperation byte = '='
)

func parseBoxes(input string) Sequence {
	input = strings.ReplaceAll(input, "\n", "")
	input = strings.ReplaceAll(input, " ", "")

	parts := strings.Split(input, ",")
	sequence := make(Sequence, len(parts))

	for i, part := range parts {
		instruction := Instruction{Operation: EqualOperation}
		k, v, ok := strings.Cut(part, string(instruction.Operation))
		if ok {
			instruction.Label = k
			instruction.Focal = aocutil.Atoi[int](v)
		} else {
			instruction.Operation = DashOperation
			k, _, ok = strings.Cut(part, string(instruction.Operation))
			if !ok {
				panic("invalid input")
			}
			instruction.Label = k
			instruction.Focal = 0
		}
		instruction.Box = hash(instruction.Label)
		sequence[i] = instruction
	}

	return sequence
}

type Boxes map[uint8][]Lens

func evaluate(seq Sequence) Boxes {
	boxes := make(map[uint8][]Lens)

	for _, instruction := range seq {
	eval:
		switch instruction.Operation {
		case '=':
			for i, lens := range boxes[instruction.Box] {
				if lens.Label == instruction.Label {
					boxes[instruction.Box][i] = instruction.Lens
					break eval
				}
			}

			boxes[instruction.Box] = append(boxes[instruction.Box], instruction.Lens)

		case '-':
			for boxID, lenses := range boxes {
				i := slices.IndexFunc(lenses, func(lens Lens) bool { return lens.Label == instruction.Label })
				if i != -1 {
					lenses = slices.Delete(lenses, i, i+1)
					boxes[boxID] = lenses
				}
			}
		}
	}

	return boxes
}

func focusingPower(boxes Boxes) int {
	var total int
	for _, lenses := range boxes {
		for i, lens := range lenses {
			power := (1 + int(lens.Box)) * (i + 1) * lens.Focal
			total += power
		}
	}
	return total
}

func part2(input string) int {
	sequence := parseBoxes(input)
	boxes := evaluate(sequence)
	log.Printf("%v", boxes)
	return focusingPower(boxes)
}
