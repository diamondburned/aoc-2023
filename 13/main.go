package main

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	input := aocutil.InputString()
	blocks := aocutil.SplitBlocks(input)

	var packetPairs [][2]Packet

	for _, block := range blocks {
		lines := aocutil.SplitLines(block)
		packets := [2]Packet{
			ParsePacket(lines[0]),
			ParsePacket(lines[1]),
		}
		packetPairs = append(packetPairs, packets)
	}

	part1(aocutil.Clone(packetPairs))
	part2(aocutil.Clone(packetPairs))
}

func part1(packetPairs [][2]Packet) {
	var indices []int
	for i, pair := range packetPairs {
		order := itemIsOrdered(pair[0], pair[1])

		log.Printf("%d: compare %v vs %v", i+1, pair[0], pair[1])
		log.Println("   ->", order)

		if order == ordered {
			indices = append(indices, i+1)
		}
	}

	fmt.Println("part 1:", aocutil.Sum(indices))
}

func part2(packetPairs [][2]Packet) {
	packets := make([]Packet, 0, len(packetPairs)*2)
	for _, pair := range packetPairs {
		packets = append(packets, pair[0], pair[1])
	}

	// Add divider packets.
	dividers := [2]Packet{
		{Packet{Data(2)}},
		{Packet{Data(6)}},
	}
	packets = append(packets, dividers[:]...)

	SortPackets(packets)

	var dividerIxs [2]int
	for i, packet := range packets {
		fmt.Println(i, packet)
		for j, divider := range dividers {
			if reflect.DeepEqual(divider, packet) {
				dividerIxs[j] = i + 1
			}
		}
	}

	fmt.Println(dividerIxs[0] * dividerIxs[1])
}

type item interface {
	item()
}

type Packet []item
type Data int

func (p Packet) item() {}
func (b Data) item()   {}

func ParsePacket(line string) Packet {
	if !strIsPacket(line) {
		panic("invalid packet")
	}

	return parsePacket(line).(Packet)
}

func parsePacket(line string) item {
	if !strIsPacket(line) {
		v := aocutil.Atoi[int](line)
		return Data(v)
	}

	line = line[1 : len(line)-1]

	var packet Packet
	var buffer string

	for i := 0; i < len(line); i++ {
		if line[i] == '[' {
			end := scanPair(line[i:])
			packet = append(packet, parsePacket(line[i:i+end+1]))

			i += end
			continue
		}

		if line[i] == ',' {
			if buffer != "" {
				packet = append(packet, parsePacket(buffer))
			}
			buffer = ""
			continue
		}

		buffer += string(line[i])
		continue
	}

	if buffer != "" {
		packet = append(packet, parsePacket(buffer))
	}

	return packet
}

// scanPair scans line until it finds an ending ] that matches the starting [.
func scanPair(line string) int {
	var depth int
	for i, c := range line {
		if c == '[' {
			depth++
		} else if c == ']' {
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func strIsPacket(line string) bool {
	return strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]")
}

func (p Packet) String() string {
	var b strings.Builder
	b.WriteByte('[')
	for i, item := range p {
		fmt.Fprint(&b, item)
		if i != len(p)-1 {
			b.WriteByte(',')
		}
	}
	b.WriteByte(']')
	return b.String()
}

type Order int8

const (
	undefinedOrder Order = iota - 1
	ordered
	unordered
)

func (o Order) String() string {
	switch o {
	case ordered:
		return "ordered"
	case unordered:
		return "unordered"
	default:
		return "undefined"
	}
}

func PacketsAreOrdered(p1, p2 Packet) Order {
	return itemIsOrdered(p1, p2)
}

func SortPackets(ps []Packet) {
	sort.Slice(ps, func(i, j int) bool {
		return PacketsAreOrdered(ps[i], ps[j]) == ordered
	})
}

func itemIsOrdered(v1, v2 item) Order {
	d1, isData1 := v1.(Data)
	d2, isData2 := v2.(Data)

	if isData1 && isData2 {
		log.Println("compare data", d1, d2)
		// Lower one comes first.
		if d1 < d2 {
			return ordered
		}
		if d1 > d2 {
			return unordered
		}
		return undefinedOrder
	}

	var p1 Packet
	var p2 Packet

	if isData1 {
		p1 = Packet{d1}
	} else {
		p1 = v1.(Packet)
	}

	if isData2 {
		p2 = Packet{d2}
	} else {
		p2 = v2.(Packet)
	}

	log.Println("compare packet", p1, p2)

	plen := aocutil.Min2(len(p1), len(p2))
	for i := 0; i < plen; i++ {
		order := itemIsOrdered(p1[i], p2[i])
		if order != undefinedOrder {
			return order
		}
	}

	if len(p1) < len(p2) {
		return ordered
	}

	if len(p1) > len(p2) {
		return unordered
	}

	return undefinedOrder
}
