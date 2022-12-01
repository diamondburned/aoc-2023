package aocutil

import (
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

// ReadFile reads a file into a string, panicking if it fails.
func ReadFile(name string) string {
	v := E2(os.ReadFile(name))
	return string(v)
}

// SplitFile splits a file into lines, trimming whitespace.
func SplitFile(name, split string) []string {
	f := ReadFile(name)
	return strings.Split(f, split)
}

// SplitFileN splits a file into lines, trimming whitespace, and panics if the
// number of lines is not n.
func SplitFileN(name, split string, n int) []string {
	f := ReadFile(name)
	parts := strings.SplitN(f, split, n)
	Assertf(len(parts) == n, "expected %d parts, got %d", n, len(parts))
	return parts
}

// Atoi converts a string to an int, panicking if it fails.
func Atoi(a string) int {
	return E2(strconv.Atoi(a))
}

// Atof converts a string to a float, panicking if it fails.
func Atof(a string) float64 {
	return E2(strconv.ParseFloat(a, 64))
}

// SlidingWindow calls fn for each window of size n in slice.
func SlidingWindow[T any](slice []T, size int, fn func([]T)) {
	if size > len(slice) {
		return
	}

	for n := range slice {
		if n+size > len(slice) {
			break
		}
		fn(slice[n : n+size])
	}
}

// Chunk splits a slice into chunks of size n, calling fn for each chunk.
func Chunk[T any](numbers []T, size int, fn func([]T)) {
	for i := 0; i < len(numbers); i += size {
		end := Min(i+size, len(numbers))
		fn(numbers[i:end])
	}
}

// Sum returns the sum of a slice of numbers.
func Sum[T constraints.Ordered](numbers []T) T {
	var sum T
	for _, n := range numbers {
		sum += n
	}
	return sum
}

// Avg returns the average of a slice of numbers.
func Avg[T constraints.Integer | constraints.Float](numbers []T) T {
	return Sum(numbers) / T(len(numbers))
}

// MinMaxes returns the min and max of a slice of numbers.
func MinMaxes[T constraints.Ordered](numbers []T) (min, max T) {
	if len(numbers) == 0 {
		return
	}

	min = numbers[0]
	max = numbers[0]

	for _, n := range numbers {
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}

	return
}

// Max returns the maximum of two numbers.
func Max[T constraints.Ordered](vs ...T) T {
	var max T
	for _, v := range vs {
		if v > max {
			max = v
		}
	}
	return max
}

// Min returns the minimum of two numbers.
func Min[T constraints.Ordered](vs ...T) T {
	var min T
	for _, v := range vs {
		if v < min {
			min = v
		}
	}
	return min
}

// Clamp returns the value clamped to the range [min, max].
func Clamp[T constraints.Ordered](n, min, max T) T {
	if n > max {
		return max
	}
	if n < min {
		return min
	}
	return n
}

// E1 asserts that err is nil, and panics if it is not.
func E1(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

// E2 asserts that err is nil, and panics if it is not. The first value is
// returned.
func E2[T any](v T, err error) T {
	if err != nil {
		log.Panicln(err)
	}
	return v
}

// E3 asserts that err is nil, and panics if it is not. The first two values are
// returned.
func E3[T1 any, T2 any](v1 T1, v2 T2, err error) (T1, T2) {
	if err != nil {
		log.Panicln(err)
	}
	return v1, v2
}

// Assert asserts that cond is true, and panics if it is not.
func Assert(cond bool, msg ...any) {
	if !cond {
		log.Panicln(msg...)
	}
}

// Assertf asserts that cond is true, and panics if it is not.
func Assertf(cond bool, f string, v ...any) {
	if !cond {
		log.Panicf(f, v...)
	}
}
