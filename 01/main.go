package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func main() {
	in, _ := os.ReadFile("input")
	scanner := bufio.NewScanner(bytes.NewReader(in))

	type elfCalorie struct {
		cals  int
		index int
	}

	elfCalories := make([]elfCalorie, 0, 100)
	var currentSum int
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			elfCalories = append(elfCalories, elfCalorie{
				cals:  currentSum,
				index: len(elfCalories),
			})
			currentSum = 0
		} else {
			v, _ := strconv.Atoi(text)
			currentSum += v
		}
	}

	sort.Slice(elfCalories, func(i, j int) bool {
		return elfCalories[i].cals > elfCalories[j].cals
	})

	fmt.Println(elfCalories[0])
	fmt.Println(elfCalories[1])
	fmt.Println(elfCalories[2])
	fmt.Println(elfCalories[0].cals + elfCalories[1].cals + elfCalories[2].cals)
}
