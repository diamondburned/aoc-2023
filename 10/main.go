package main

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	var program []Instruction

	input := aocutil.InputString()
	lines := aocutil.SplitLines(input)

	for _, line := range lines {
		parts := strings.Fields(line)

		instrFn, ok := instructions[parts[0]]
		aocutil.Assertf(ok, "unknown instruction %q", parts[0])

		program = append(program, instrFn(parts[1:]))
	}

	{
		x := 1
		var totalCycles int

		var strengths []int

		for _, instr := range program {
			for cycle := InstructionCycle(instr); cycle > 0; cycle-- {
				totalCycles++

				if totalCycles == 20 || (totalCycles > 20 && ((totalCycles+20)%40) == 0) {
					strengths = append(strengths, x*totalCycles)
					log.Println("cycle", totalCycles, "strength", x*totalCycles)
				}
			}

			switch instr := instr.(type) {
			case Addx:
				x += instr.N
			}
		}

		fmt.Println(strengths)
		fmt.Println(aocutil.Sum(strengths))
	}
	fmt.Println("=======")
	{
		x := 1
		var totalCycles int

		const crtW = 40
		const crtH = 6
		var crt [crtW * crtH]bool
		var crtX int
		var crtY int

		for _, instr := range program {
			for cycle := InstructionCycle(instr); cycle > 0; cycle-- {
				// Draw if x is +-1 within pos.
				if crtX-1 <= x && x <= crtX+1 {
					crt[crtY*crtW+crtX] = true
				}

				totalCycles++
				crtX++
				if crtX == crtW {
					crtX = 0
					crtY++
				}

			}

			switch instr := instr.(type) {
			case Addx:
				x += instr.N
			}
		}

		for y := 0; y < crtH; y++ {
			for x := 0; x < crtW; x++ {
				if crt[y*crtW+x] {
					fmt.Print("#")
				} else {
					fmt.Print(".")
				}
			}

			fmt.Println()
		}
	}
}

type Instruction interface {
	instruction()
}

func InstructionCycle(instr Instruction) int {
	return 1 + reflect.TypeOf(instr).NumField()
}

var instructions = map[string]func(args []string) Instruction{
	"noop": func(args []string) Instruction { return Noop{} },
	"addx": func(args []string) Instruction {
		return Addx{aocutil.Atoi[int](args[0])}
	},
}

type Noop struct{}
type Addx struct{ N int }

func (Noop) instruction() {}
func (Addx) instruction() {}
