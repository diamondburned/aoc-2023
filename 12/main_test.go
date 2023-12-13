package main

import (
	"fmt"
	"log"
	"testing"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func TestTotalValid(t *testing.T) {
	log.SetFlags(0)

	tests := []struct {
		in1  SpringsConditions
		in2  []int
		want int
	}{
		{`??`, []int{1}, 2},
		{`???.###`, []int{1, 1, 3}, 1},
		{`.??..??...?##.`, []int{1, 1, 3}, 4},
		{`?#?#?#?#?#?#?#?`, []int{1, 3, 1, 6}, 1},
		{`????.#...#...`, []int{4, 1, 1}, 1},
		{`????.######..#####.`, []int{1, 6, 5}, 4},
		{`?###????????`, []int{3, 2, 1}, 10},
	}
	var fail bool
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if fail {
				t.Skip()
			}
			in := SpringsRecord{tt.in1, tt.in2}
			if got := countValid(in); got != tt.want {
				t.Errorf("TotalValid(%v) = %v, want %v", in, got, tt.want)
				fail = true
			}
		})
	}
}

func BenchmarkTotalValid(b *testing.B) {
	aocutil.SilenceLogging()

	const awfulInput = `????????????????# 8,2,1,1`
	input := parseInput(awfulInput)[0]
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		countValid(input)
	}
}
