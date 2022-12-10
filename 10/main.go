package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	input := aocutil.InputString()
	lines := aocutil.SplitLines(input)
	program := make(Program, 0, len(lines))

	for _, line := range lines {
		parts := strings.Fields(line)

		parseInstruction, ok := instructionParsers[parts[0]]
		aocutil.Assertf(ok, "unknown instruction %q", parts[0])

		program = append(program, parseInstruction(parts[1:]))
	}

	part1(program)

	fmt.Println()
	fmt.Println("=======")
	fmt.Println()

	part2(program)
}

func part1(program Program) {
	var strengths []int
	x := 1

	program.Run(func(instr Instruction, cycles int) {
		if cycles == 20 || ((cycles+20)%40) == 0 {
			strengths = append(strengths, x*cycles)
		}

		switch instr := instr.(type) {
		case Addx:
			x += instr.N
		}
	})

	fmt.Println("part 1:")
	fmt.Println("  strengths:", strengths)
	fmt.Println("  strengths sum:", aocutil.Sum(strengths))
}

func part2(program Program) {
	var crtX, crtY int
	const crtW = 40
	const crtH = 6
	crt := NewCRT(crtW, crtH)
	x := 1

	program.Run(func(instr Instruction, cycles int) {
		// Draw if x is +-1 within pos.
		if crtX-1 <= x && x <= crtX+1 {
			crt.Set(crtX, crtY, true)
		}

		if crtX++; crtX == crtW {
			crtX = 0
			crtY++
		}

		switch instr := instr.(type) {
		case Addx:
			x += instr.N
		}
	})

	fmt.Println("part 2:")
	fmt.Println()
	crt.Display(os.Stdout)
}

// Program is a program.
type Program []Instruction

// Run runs the CPU with the given program. The total number of CPU cycles is
// returned.
func (p Program) Run(f func(instr Instruction, cycles int)) int {
	var cycles int
	for _, instruction := range p {
		for i, j := 0, InstructionCycles(instruction)-1; i < j; i++ {
			cycles++
			f(nil, cycles)
		}
		cycles++
		f(instruction, cycles)
	}
	return cycles
}

// Instruction is a CPU instruction.
type Instruction interface {
	instruction()
}

// InstructionCycles returns the number of cycles that the instruction takes.
func InstructionCycles(instr Instruction) int {
	return 1 + reflect.TypeOf(instr).NumField()
}

var instructionParsers = map[string]func(args []string) Instruction{
	"noop": func(args []string) Instruction { return Noop{} },
	"addx": func(args []string) Instruction {
		return Addx{aocutil.Atoi[int](args[0])}
	},
}

type Noop struct{}
type Addx struct{ N int }

func (Noop) instruction() {}
func (Addx) instruction() {}

// CRT describes a black-and-white CRT screen.
type CRT struct {
	Pix    []bool
	Stride int
}

// NewCRT creates a new CRT with the given width and height.
func NewCRT(w, h int) CRT {
	return CRT{
		Pix:    make([]bool, w*h),
		Stride: w,
	}
}

// Set sets the pixel at the given position.
func (crt CRT) Set(x, y int, val bool) {
	i := y*crt.Stride + x
	if i < 0 || i >= len(crt.Pix) {
		return
	}

	crt.Pix[i] = val
}

// Display displays the CRT on the console.
func (crt CRT) Display(out io.Writer) {
	buf := make([]rune, len(crt.Pix))
	for i := range buf {
		if crt.Pix[i] {
			buf[i] = '⬜'
		} else {
			buf[i] = 'ㅤ'
		}
	}

	for i := 0; i < len(crt.Pix); i += crt.Stride {
		io.WriteString(out, string(buf[i:i+crt.Stride]))
		io.WriteString(out, "\n")
	}
}
