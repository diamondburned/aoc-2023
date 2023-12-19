package main

import (
	"testing"
)

func TestCountStraightMoves(t *testing.T) {
	tests := []struct {
		in   pointPath
		want int
	}{
		{
			in:   pointPath{{0, 0}, {1, 0}, {2, 0}, {3, 0}},
			want: 3,
		},
		{
			in:   pointPath{{0, 1}, {0, 0}, {1, 0}, {2, 0}, {3, 0}},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := tt.in.countStraightMoves(); got != tt.want {
				t.Errorf("countStraightMoves(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}
