package aocutil

import (
	"bufio"
	"bytes"
	"container/heap"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"

	_ "github.com/davecgh/go-spew/spew"
	_ "gonum.org/v1/gonum"
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
		end := Min2(i+size, len(numbers))
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

// Maxs returns the maximum of many values.
func Maxs[T constraints.Ordered](vs ...T) T {
	var max T
	for _, v := range vs {
		if v > max {
			max = v
		}
	}
	return max
}

// Max2 returns the maximum of two values.
func Max2[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Mins returns the minimum of many values.
func Mins[T constraints.Ordered](vs ...T) T {
	var min T
	for _, v := range vs {
		if v < min {
			min = v
		}
	}
	return min
}

// Min2 returns the minimum of two values.
func Min2[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
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

// SplitN splits a string into n parts, and panics if the number of parts is not
// n.
func SplitN(s, sep string, n int) []string {
	parts := strings.SplitN(s, sep, n)
	Assertf(len(parts) == n, "SplitN(%q, %q, %d): got %d", s, sep, n, len(parts))
	return parts
}

// FieldsN splits a string into n fields, and panics if the number of fields is
// not n.
func FieldsN(s string, n int) []string {
	parts := strings.Fields(s)
	Assertf(len(parts) == n, "FieldsN(%q, %d): got %d", s, n, len(parts))
	return parts
}

// Heap is a heap of ordered values.
type Heap[T constraints.Ordered] internalHeap[T]

type HeapOpts[T constraints.Ordered] struct {
	// Less is a custom less function. If nil, min or max will be used. This
	// takes precedence over Max.
	Less func(a, b T) bool
	// Max is true if the heap is a max heap. Largest values will be put first.
	// Otherwise, the heap is a min heap.
	Max bool
	// Cap is the preallocated capacity of the heap.
	Cap int
	// Limit determines the behavior of Push when the heap is full. A zero-value
	// sets no limit.
	Limit int
}

func (h HeapOpts[T]) less(a, b T) bool {
	if h.Less != nil {
		return h.Less(a, b)
	}
	if h.Max {
		return a > b
	}
	return a < b
}

func (h HeapOpts[T]) more(a, b T) bool {
	if h.Less != nil {
		// https://pkg.go.dev/sort#Interface
		if h.Less(a, b) {
			// a < b
			return false
		}
		if h.Less(b, a) {
			// !(a < b) && b < a
			// = a >= b && b < a
			// = a > b
			return true
		}
		// !(a < b) && !(b < a)
		// = a >= b && b >= a
		// = a == b
		return false
	}
	if h.Max {
		return a < b
	}
	return a > b
}

// NewMinHeap returns a new min heap.
func NewMinHeap[T constraints.Ordered]() *Heap[T] {
	return NewHeapOpts(HeapOpts[T]{Max: false})
}

// NewMaxHeap returns a new max heap.
func NewMaxHeap[T constraints.Ordered]() *Heap[T] {
	return NewHeapOpts(HeapOpts[T]{Max: true})
}

// NewHeapOpts creates a new heap with options.
func NewHeapOpts[T constraints.Ordered](opts HeapOpts[T]) *Heap[T] {
	slice := make([]T, 0, opts.Cap)
	h := internalHeap[T]{
		heap: slice,
		opts: opts,
	}
	heap.Init((*internalHeap[T])(&h))
	return (*Heap[T])(&h)
}

// Push pushes a value onto the heap.
func (h *Heap[T]) Push(v T) {
	if h.opts.Limit == 0 || len(h.heap) < h.opts.Limit {
		heap.Push((*internalHeap[T])(h), v)
		return
	}

	min := h.heap[0]
	if h.opts.more(v, min) {
		x := h.Pop()
		heap.Push((*internalHeap[T])(h), v)
		log.Println("heap limit exceeded, popped", x, "for", v)
	}
}

// Pop pops a value from the heap.
func (h *Heap[T]) Pop() T {
	return heap.Pop((*internalHeap[T])(h)).(T)
}

// Fix fixes the heap after a value has been modified at the given index.
func (h *Heap[T]) Fix(i int) {
	heap.Fix((*internalHeap[T])(h), i)
}

// Len returns the number of items in the heap.
func (h *Heap[T]) Len() int {
	return len(h.heap)
}

// Sort sorts the heap.
func (h *Heap[T]) Sort() {
	sort.Sort((*internalHeap[T])(h))
}

// ToSlice returns the underlying heap slice.
func (h *Heap[T]) ToSlice() []T {
	return h.heap
}

// ToSortedSlice returns a sorted slice of the heap.
func (h *Heap[T]) ToSortedSlice() []T {
	h.Sort()
	return h.heap
}

type internalHeap[T constraints.Ordered] struct {
	heap []T
	opts HeapOpts[T]
}

var _ heap.Interface = (*internalHeap[int])(nil)
var _ sort.Interface = (*internalHeap[int])(nil)

func (h internalHeap[T]) Len() int {
	return len(h.heap)
}

func (h internalHeap[T]) Less(i, j int) bool {
	return h.opts.less(h.heap[i], h.heap[j])
}

func (h internalHeap[T]) Swap(i, j int) {
	h.heap[i], h.heap[j] = h.heap[j], h.heap[i]
}

// Push pushes the element x onto the heap.
func (h *internalHeap[T]) Push(x interface{}) {
	h.heap = append(h.heap, x.(T))
}

// Pop removes the minimum element (according to Less) from the heap and returns
// it.
func (h *internalHeap[T]) Pop() interface{} {
	old := h.heap
	n := len(old)
	x := old[n-1]
	h.heap = old[0 : n-1]
	return x
}

// Scanner is a scanner for a reader.
type Scanner struct {
	s *bufio.Scanner
}

// NewScanner returns a new scanner for a reader.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{s: bufio.NewScanner(r)}
}

// NewBytesScanner returns a new scanner for a byte slice.
func NewBytesScanner(b []byte) *Scanner {
	return NewScanner(bytes.NewReader(b))
}

// SetSplitter sets the split function for the underlying scanner.
func (s *Scanner) SetSplitter(scanner bufio.SplitFunc) {
	s.s.Split(scanner)
}

// Next returns the next token.
func (s *Scanner) Next() bool {
	return s.s.Scan()
}

// Token returns the current token.
func (s *Scanner) Token() string {
	return s.s.Text()
}
