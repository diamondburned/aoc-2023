package main

import "testing"

func TestHash(t *testing.T) {
	tests := []struct {
		in  string
		out uint8
	}{
		{"HASH", 52},
	}

	for _, test := range tests {
		if v := hash(test.in); v != test.out {
			t.Errorf("hash(%q) = %d, want %d", test.in, v, test.out)
		}
	}
}
