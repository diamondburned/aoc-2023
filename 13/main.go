package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"reflect"
	"sort"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	log.SetOutput(io.Discard)

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

		if order == Ordered {
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
		log.Printf("%02d: %v", i, packet)
		for j, divider := range dividers {
			if reflect.DeepEqual(divider, packet) {
				dividerIxs[j] = i + 1
			}
		}
	}

	fmt.Println("part 2:", dividerIxs[0]*dividerIxs[1])
}

type item interface {
	item()
}

type Packet []item
type Data int

func (p Packet) item() {}
func (b Data) item()   {}

func unmarshalItemJSON(b []byte) (item, error) {
	var raws []json.RawMessage
	if err := json.Unmarshal(b, &raws); err != nil {
		var v int
		if err := json.Unmarshal(b, &v); err != nil {
			return nil, err
		}
		return Data(v), nil
	}

	items := make([]item, len(raws))
	for i, raw := range raws {
		item, err := unmarshalItemJSON(raw)
		if err != nil {
			return nil, err
		}
		items[i] = item
	}

	return Packet(items), nil
}

func ParsePacket(line string) Packet {
	p := aocutil.E2(unmarshalItemJSON([]byte(line)))
	return p.(Packet)
}

func (p Packet) String() string {
	b := aocutil.E2(json.Marshal(p))
	return string(b)
}

type Order int8

const (
	UndefinedOrder Order = iota - 1
	Ordered
	Unordered
)

func (o Order) String() string {
	switch o {
	case Ordered:
		return "ordered"
	case Unordered:
		return "unordered"
	default:
		return "undefined"
	}
}

// SortPackets sorts the packets in-place.
func SortPackets(ps []Packet) {
	sort.Slice(ps, func(i, j int) bool {
		return PacketsAreOrdered(ps[i], ps[j]) == Ordered
	})
}

// PacketsAreOrdered returns whether p1 is ordered before p2.
func PacketsAreOrdered(p1, p2 Packet) Order {
	return itemIsOrdered(p1, p2)
}

func itemIsOrdered(v1, v2 item) Order {
	d1, isData1 := v1.(Data)
	d2, isData2 := v2.(Data)

	if isData1 && isData2 {
		log.Println("compare data", d1, d2)
		// Lower one comes first.
		if d1 < d2 {
			return Ordered
		}
		if d1 > d2 {
			return Unordered
		}
		return UndefinedOrder
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
		if order != UndefinedOrder {
			return order
		}
	}

	if len(p1) < len(p2) {
		return Ordered
	}

	if len(p1) > len(p2) {
		return Unordered
	}

	return UndefinedOrder
}
