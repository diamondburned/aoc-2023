package aocutil

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestMap2D_Transpose(t *testing.T) {
	m := NewMap2DFromData([][]byte{
		{1, 2, 3},
		{4, 5, 6},
	})

	assert.Equal(t, NewMap2DFromData([][]byte{
		{1, 4},
		{2, 5},
		{3, 6},
	}), m.Transpose())
}
