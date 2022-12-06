package main

import (
	"bytes"
	"fmt"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	input := []byte(aocutil.Input())

	const packetStartLen = 4
	packetStartIx := packetStartLen
	var packetStart []byte

	aocutil.SlidingWindow(input, packetStartLen, func(window []byte) bool {
		// Ensure that the bytes in this window are all unique.
		if !isUniq(window) {
			packetStartIx++
			return false
		}

		packetStart = window
		return true
	})

	fmt.Println(packetStartIx, string(packetStart))

	const messageStartLen = 14
	messageStartIx := messageStartLen
	var messageStart []byte

	aocutil.SlidingWindow(input, messageStartLen, func(window []byte) bool {
		// Ensure that the bytes in this window are all unique.
		if !isUniq(window) {
			messageStartIx++
			return false
		}

		messageStart = window
		return true
	})

	fmt.Println(messageStartIx, string(messageStart))
}

func isUniq(b []byte) bool {
	for _, c := range b {
		if bytes.Count(b, []byte{c}) > 1 {
			return false
		}
	}
	return true
}
