package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
)

func main() {
	var prioritySum int
	var groupIx int

	const groupSize = 3
	groupRucksacks := make([][groupSize]string, 1, 32)

	input, _ := os.ReadFile("input")
	scanner := bufio.NewScanner(bytes.NewReader(input))
	for scanner.Scan() {
		str := scanner.Text()
		if len(str) == 0 {
			continue
		}

		if groupIx < groupSize {
			groupRucksacks[len(groupRucksacks)-1][groupIx] = str
			groupIx++
			continue
		}

		groupIx = 1
		groupRucksacks = append(groupRucksacks, [3]string{str})
	}

	for groupIx, group := range groupRucksacks {
		// rucksackSet := make(map[byte]struct{}, 10)
		// items := [3]byte{}

		// for rucksackIx, rucksack := range group {
		// 	log.Printf("group %d, rucksack %d: %s", groupIx, rucksackIx, rucksack)
		// 	part1 := rucksack[:len(rucksack)/2]
		// 	part2 := rucksack[len(rucksack)/2:]
		// 	for _, b := range part1 {
		// 		rucksackSet[b] = struct{}{}
		// 	}
		// 	for _, b := range part2 {
		// 		_, ok := rucksackSet[b]
		// 		if ok {
		// 			items[rucksackIx] = b
		// 			break
		// 		}
		// 	}
		// }

		// var currentPriority int
		// var currentCount int
		// var currentItem byte

		// for b, count := range rucksackSet {
		// 	priority := priority(b)
		// 	if count > groupSize && count > currentCount && priority > currentPriority {
		// 		currentPriority = priority
		// 		currentCount = count
		// 		currentItem = b
		// 	}
		// }

		// prioritySum += currentPriority
		// log.Printf("for group %d, item %c, count %d, priority is %d", groupIx, currentItem, currentCount, currentPriority)
		// continue

		rucksackSet := make(map[byte]int, 10)
		// items := [3]byte{}

		for _, rucksack := range group {
			miniset := make(map[byte]struct{}, 10)
			for _, b := range []byte(rucksack) {
				_, ok := miniset[b]
				if !ok {
					rucksackSet[b]++
					miniset[b] = struct{}{}
				}
			}
		}

		var currentCount int
		var currentItem byte

		for b, count := range rucksackSet {
			if count == groupSize {
				currentCount = count
				currentItem = b
				break
			}
		}

		priority := priority(currentItem)
		log.Printf("for group %d, item %c, count %d, priority is %d", groupIx, currentItem, currentCount, priority)
		prioritySum += priority

		// var currentPriority int
		// var currentCount int
		// var currentItem byte

		// for b, count := range rucksackSet {
		// 	priority := priority(b)
		// 	if count > groupSize && count > currentCount && priority > currentPriority {
		// 		currentPriority = priority
		// 		currentCount = count
		// 		currentItem = b
		// 	}
		// }

		// prioritySum += currentPriority
		// continue
	}
	// part1Raw := bytes[:len(bytes)/2]
	// part2Raw := bytes[len(bytes)/2:]

	// backpack := make(map[byte]struct{}, len(part1Raw))
	// for _, b := range part1Raw {
	// 	backpack[b] = struct{}{}
	// }
	// for _, b := range part2Raw {
	// 	_, ok := backpack[b]
	// 	if ok {
	// 		prioritySum += priority(b)
	// 		log.Printf("for %c, priority is %d", b, priority(b))
	// 		break
	// 	}
	// }

	fmt.Println(prioritySum)
}

func priority(b byte) int {
	if 'a' <= b && b <= 'z' {
		return int(b-'a') + 1
	}
	if 'A' <= b && b <= 'Z' {
		return int(b-'A') + 27
	}
	return 0
}
