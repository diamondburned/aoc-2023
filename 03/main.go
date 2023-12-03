package main

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"strconv"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	input := aocutil.ReadStdin()

	// part1(input)
	part2(input)
}

type Schematic [][]byte

func parseSchematic(input string) Schematic {
	lines := aocutil.SplitLines(input)
	schematic := make(Schematic, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			continue
		}
		schematic = append(schematic, []byte(line))
	}
	return schematic
}

const Gear = '*'

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func isSymbol(b byte) bool {
	return b != 0 && !isDigit(b) && b != '.'
}

func (s Schematic) At(x, y int) byte {
	if x < 0 || y < 0 || y >= len(s) || x >= len(s[y]) {
		return 0
	}
	return s[y][x]
}

func part1(input string) {
	schematic := parseSchematic(input)
	var sum int
	for y, row := range schematic {
		// log.Printf("row: %q", row)
		for x := 0; x < len(row); x++ {
			start := x
			digitEnd := bytes.IndexFunc(row[start:], func(r rune) bool {
				return !isDigit(byte(r))
			})
			if digitEnd == -1 {
				digitEnd = len(row) - start
			}

			// log.Printf("y: %d, start: %d, digitEnd: %d", y, start, digitEnd)

			digits := row[start : start+digitEnd]
			if len(digits) == 0 {
				continue
			}

			// Check immediate top and bottom cells.
			var nearSymbol bool
			for dy := -1; dy <= 1 && !nearSymbol; dy++ {
				for dx := -1; dx <= len(digits) && !nearSymbol; dx++ {
					if dx == 0 && dy == 0 {
						continue
					}
					if isSymbol(schematic.At(x+dx, y+dy)) {
						nearSymbol = true
						// log.Printf("near symbol: %q at (%d, %d)", digits, x+dx, y+dy)
					}
				}
			}

			if nearSymbol {
				// log.Printf("found %q at (%d, %d)", digits, x, y)

				n, err := strconv.Atoi(string(digits))
				if err != nil {
					log.Panicf("failed to parse %q: %v", digits, err)
				}
				sum += n
			}

			x += digitEnd
		}
	}

	fmt.Println(sum)
}

func part2(input string) {
	schematic := parseSchematic(input)

	// gears tracks the gears at each point and its adjacent numbers.
	gears := make(map[image.Point][]int)

	for y, row := range schematic {
		// log.Printf("row: %q", row)
		for x := 0; x < len(row); x++ {
			start := x
			digitEnd := bytes.IndexFunc(row[start:], func(r rune) bool {
				return !isDigit(byte(r))
			})
			if digitEnd == -1 {
				digitEnd = len(row) - start
			}

			// log.Printf("y: %d, start: %d, digitEnd: %d", y, start, digitEnd)

			digits := row[start : start+digitEnd]
			if len(digits) == 0 {
				continue
			}

			n, err := strconv.Atoi(string(digits))
			if err != nil {
				log.Panicf("failed to parse %q: %v", digits, err)
			}

			// Check immediate top and bottom cells.
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= len(digits); dx++ {
					if dx == 0 && dy == 0 {
						continue
					}
					cell := schematic.At(x+dx, y+dy)
					if cell == Gear {
						log.Printf("found gear at (%d, %d)", x+dx, y+dy)
						pt := image.Pt(x+dx, y+dy)
						gears[pt] = append(gears[pt], n)
					}
					// log.Printf("near symbol: %q at (%d, %d)", digits, x+dx, y+dy)
				}
			}

			x += digitEnd
		}
	}

	var sum int

	for _, nums := range gears {
		if len(nums) != 2 {
			continue
		}

		gearRatio := 1
		for _, n := range nums {
			gearRatio *= n
		}
		sum += gearRatio
	}

	fmt.Println(sum)
}
