package main

import (
	"fmt"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	input := aocutil.InputBytes()

	const packetStartLen = 4
	packetStartIx := aocutil.SlidingWindow(input, packetStartLen, aocutil.IsUniq[byte])
	fmt.Println(packetStartIx+packetStartLen, string(input[packetStartIx:packetStartIx+packetStartLen]))

	const messageStartLen = 14
	msgStartIx := aocutil.SlidingWindow(input, messageStartLen, aocutil.IsUniq[byte])
	fmt.Println(msgStartIx+messageStartLen, string(input[msgStartIx:msgStartIx+messageStartLen]))
}
