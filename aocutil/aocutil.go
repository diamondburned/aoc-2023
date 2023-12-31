package aocutil

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"

	"github.com/aclements/go-moremath/fit"
	"github.com/mohae/deepcopy"
	"golang.org/x/exp/constraints"
	"gonum.org/v1/gonum/stat/combin"

	_ "github.com/davecgh/go-spew/spew"
	_ "github.com/sourcegraph/conc"
	_ "gonum.org/v1/gonum"
)

var silent atomic.Bool

var (
	part1Only = false
	part2Only = false
)

func init() {
	w := &logPrefixedWriter{
		lastTime: time.Now(),
		noColors: false,
	}

	log.SetFlags(0)
	log.SetOutput(w)

	if !testing.Testing() {
		silent_ := flag.Bool("s", false, "suppress output")
		flag.BoolVar(&part1Only, "1", false, "run only part 1")
		flag.BoolVar(&part2Only, "2", false, "run only part 2")
		flag.Parse()
		silent.Store(*silent_)
	}

	if silent.Load() {
		SilenceLogging()
	}
}

type logPrefixedWriter struct {
	lastTime time.Time
	noColors bool
}

func (w *logPrefixedWriter) Write(b []byte) (n int, err error) {
	now := time.Now()
	d := now.Sub(w.lastTime)
	w.lastTime = now

	const f1 = "+% 3d.%06ds ⎸ "
	const f2 = "\033[38;5;248m" + f1 + "\033[0m"
	pf := f2
	if w.noColors {
		pf = f1
	}
	p := fmt.Sprintf(pf, int(d.Seconds()), int(d.Microseconds())%1000000)
	return writePrefixed(os.Stderr, b, []byte(p+""+log.Prefix()))
}

// ScopeLogf enters a log scope.
func ScopeLogf(f string, v ...any) (unscope func()) {
	old := log.Prefix()
	new := old
	if f != "" {
		new += fmt.Sprintf(f, v...) + ": "
	} else {
		new += "."
	}
	log.SetPrefix(new)
	return func() { log.SetPrefix(old) }
}

// IsSilent returns true if the -s flag is given.
func IsSilent() bool { return silent.Load() }

// SilenceLogging disables logging.
func SilenceLogging() {
	silent.Store(true)
	log.SetOutput(io.Discard)
}

// Run runs the given functions with the stdin input.
func Run(p1, p2 func(string) int) {
	input := ReadStdin()
	switch {
	case part1Only:
		fmt.Println(p1(input))
	case part2Only:
		fmt.Println(p2(input))
	default:
		fmt.Println(p1(input))
		fmt.Println(p2(input))
	}
}

// ParseAndRun runs the given functions with the input after parsing it.
func ParseAndRun[T any](parse func(string) T, p1, p2 func(T) int) {
	input := ReadStdin()
	value := parse(input)
	switch {
	case part1Only:
		fmt.Println(p1(value))
	case part2Only:
		fmt.Println(p2(value))
	default:
		fmt.Println(p1(value))
		fmt.Println(p2(value))
	}
}

// ReadFile reads a file into a string, panicking if it fails.
func ReadFile(name string) string {
	v := E2(os.ReadFile(name))
	return string(v)
}

// ReadStdin reads stdin into a string, panicking if it fails.
func ReadStdin() string {
	v := E2(io.ReadAll(os.Stdin))
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

// SplitBlocks splits the given string using two new lines after trimming new
// lines.
func SplitBlocks(s string) []string {
	s = strings.Trim(s, "\n")
	return strings.Split(s, "\n\n")
}

// SplitLines splits a string into lines after trimming trailing new lines.
func SplitLines(s string) []string {
	s = strings.Trim(s, "\n")
	ss := strings.Split(s, "\n")
	ss = FilterEmptyStrings(ss)
	return ss
}

// SplitLineFields splits a string into lines and then fields.
func SplitLineFields(s string) [][]string {
	lines := SplitLines(s)
	fields := make([][]string, len(lines))
	for i, line := range lines {
		fields[i] = strings.Fields(line)
	}
	return fields
}

// FilterEmptyStrings returns strs with empty strings removed.
func FilterEmptyStrings(strs []string) []string {
	strs2 := strs[:0]
	for _, s := range strs {
		if strings.TrimSpace(s) != "" {
			strs2 = append(strs2, s)
		}
	}
	return strs2
}

// ReplaceStringRange replaces the range str[i:j] by the given substr.
func ReplaceStringRange[T ~string](str T, i, j int, substr T) T {
	return str[:i] + substr + str[j:]
}

// ReplaceStringIndex replaces the character at index i by the given substr.
func ReplaceStringIndex[T ~string](str T, i int, substr T) T {
	return ReplaceStringRange(str, i, i+1, substr)
}

// Sscanf is a wrapper around fmt.Sscanf that panics on error.
func Sscanf(s string, format string, args ...interface{}) {
	_, err := fmt.Sscanf(s, format, args...)
	Assertf(err == nil, "Sscanf(%q, %q, ...): %v", s, format, err)
}

// Atoi converts a string to an int, panicking if it fails.
func Atoi[T constraints.Signed](a string) T {
	v, err := strconv.ParseInt(a, 10, int(unsafe.Sizeof(T(0))*8))
	Assertf(err == nil, "failed to parse int: %v", err)
	return T(v)
}

// Atois converts a slice of strings to a slice of ints, panicking if it fails.
func Atois[T constraints.Signed](a []string) []T {
	return Map(a, Atoi[T])
}

// Itoa converts an int to a string.
func Itoa[T constraints.Signed](v T) string {
	s := strconv.FormatInt(int64(v), 10)
	return s
}

// Itoas converts a slice of ints to a slice of strings.
func Itoas[T constraints.Signed](a []T) []string {
	return Map(a, Itoa[T])
}

// Atou converts a string to an uint, panicking if it fails.
func Atou[T constraints.Unsigned](a string) uint {
	v, err := strconv.ParseUint(a, 10, int(unsafe.Sizeof(T(0))*8))
	Assertf(err == nil, "failed to parse uint: %v", err)
	return uint(v)
}

// Atous converts a slice of strings to a slice of uints, panicking if it fails.
func Atous[T constraints.Unsigned](a []string) []uint {
	return Map(a, Atou[T])
}

// Atof converts a string to a float, panicking if it fails.
func Atof[T constraints.Float](a string) T {
	return T(E2(strconv.ParseFloat(a, 64)))
}

// Atofs converts a slice of strings to a slice of floats, panicking if it
// fails.
func Atofs[T constraints.Float](a []string) []T {
	return Map(a, Atof[T])
}

// Trim removes the value v from the beginning and end of the slice.
func Trim[T comparable](slice []T, v T) []T {
	return trim(slice, v, trimBoth)
}

// TrimLeft removes the value v from the beginning of the slice.
func TrimLeft[T comparable](slice []T, v T) []T {
	return trim(slice, v, trimLeft)
}

// TrimRight removes the value v from the end of the slice.
func TrimRight[T comparable](slice []T, v T) []T {
	return trim(slice, v, trimRight)
}

type trimType uint8

const (
	_ trimType = iota
	trimLeft
	trimRight
	trimBoth
)

func (t trimType) String() string {
	switch t {
	case trimLeft:
		return "left"
	case trimRight:
		return "right"
	case trimBoth:
		return "both"
	default:
		return fmt.Sprintf("trimType(%d)", t)
	}
}

func trim[T comparable](slice []T, v T, trim trimType) []T {
	switch v := any(v).(type) {
	// Fast path.
	case byte:
		src := any(slice).([]byte)
		var dst []byte
		switch trim {
		case trimLeft:
			dst = bytes.TrimLeft(src, string(v))
		case trimRight:
			dst = bytes.TrimRight(src, string(v))
		case trimBoth:
			dst = bytes.Trim(src, string(v))
		}
		return any(dst).([]T)
	default:
		if trim == trimBoth || trim == trimLeft {
			for i := 0; i < len(slice); i++ {
				if slice[i] != v {
					slice = slice[i:]
					break
				}
			}
		}
		if trim == trimBoth || trim == trimRight {
			for i := len(slice) - 1; i >= 0; i-- {
				if slice[i] != v {
					slice = slice[:i+1]
					break
				}
			}
		}
		return slice
	}
}

//
// Have we gone too far?
//

func sliceAsBytes[T any](slice []T) []byte {
	if len(slice) == 0 {
		return nil
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(&slice[0])), len(slice)*int(unsafe.Sizeof(slice[0])))
}

func valueAsBytes[T any](value *T) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(value)), unsafe.Sizeof(value))
}

// Count returns the number of times v appears in the slice.
func Count[T comparable](slice []T, v T) int {
	if b, ok := any(v).(byte); ok {
		src := any(slice).([]byte)
		return bytes.Count(src, []byte{b})
	}

	if unsafe.Sizeof(v) < 32 {
		return bytes.Count(sliceAsBytes(slice), valueAsBytes(&v))
	}

	var count int
	for _, x := range slice {
		if x == v {
			count++
		}
	}
	return count
}

// MaybeAt returns the element at index i, or nil if i is out of bounds.
func MaybeAt[T comparable](slice []T, i int) *T {
	if i < 0 || i >= len(slice) {
		return nil
	}
	return &slice[i]
}

// Index returns the index of the first instance of v in slice, or -1 if v is
// not in slice.
func Index[T comparable](slice []T, v T) int {
	if b, ok := any(v).(byte); ok {
		src := any(slice).([]byte)
		return bytes.IndexByte(src, b)
	}

	if unsafe.Sizeof(v) < 32 {
		bslice := sliceAsBytes(slice)
		bvalue := valueAsBytes(&v)

		idx := bytes.Index(bslice, bvalue)
		if idx == -1 {
			return -1
		}

		return idx / len(bvalue)
	}

	for i, x := range slice {
		if x == v {
			return i
		}
	}
	return -1
}

// Contains returns true if the slice contains the value.
func Contains[T comparable](slice []T, v T) bool {
	return Index(slice, v) != -1
}

// Map transforms a slice of strings into another slice of strings.
func Map[T1 any, T2 any](a []T1, f func(T1) T2) []T2 {
	v := make([]T2, len(a))
	for i, s := range a {
		v[i] = f(s)
	}
	return v
}

// Filter filters a slice of strings into another slice of strings.
func Filter[T any](a []T, f func(T) bool) []T {
	v := make([]T, 0, Clamp(len(a), 0, 128))
	for _, s := range a {
		if f(s) {
			v = append(v, s)
		}
	}
	return v
}

// FilterInplace filters a slice of strings in place. The given slice is
// modified.
func FilterInplace[T any](a []T, f func(T) bool) []T {
	v := a[:0]
	for _, s := range a {
		if f(s) {
			v = append(v, s)
		}
	}
	return v
}

// FilterIxs filters a slice and returns the indices at which f(a[i]) returns
// true.
func FilterIxs[T any](a []T, f func(T) bool) []int {
	ixs := make([]int, 0, Clamp(len(a), 0, 128))
	for i, s := range a {
		if f(s) {
			ixs = append(ixs, i)
		}
	}
	return ixs
}

// SlidingWindow calls fn for each window of size n in slice. If fn returns
// true, then the function will break.
func SlidingWindow[T any](slice []T, size int, fn func([]T) bool) int {
	if size > len(slice) {
		log.Panicf(
			"SlidingWindow: size %d is larger than slice length %d",
			size, len(slice),
		)
	}

	for n := range slice {
		if n+size > len(slice) {
			break
		}
		if fn(slice[n : n+size]) {
			return n
		}
	}

	return -1
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

// Mul returns the multiplication of a slice of numbers.
func Mul[T constraints.Integer | constraints.Float](numbers []T) T {
	mul := T(1)
	for _, n := range numbers {
		mul *= n
	}
	return mul
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

// Minf returns the minimum of many values, with each one being compared
// using the given less function.
func Minf[T any](vs []T, less func(T, T) bool) T {
	if len(vs) == 0 {
		var z T
		return z
	}

	min := vs[0]
	for _, v := range vs[1:] {
		if less(v, min) {
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

// Abs returns the absolute value of n.
func Abs[T constraints.Integer | constraints.Float | constraints.Signed](n T) T {
	if n < 0 {
		return -n
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

// CompareOrdered compares two ordered values. It returns 0 if they are equal, 1
// if a > b, and -1 if a < b.
func CompareOrdered[T constraints.Ordered](a, b T) int {
	if a == b {
		return 0
	}
	if a < b {
		return -1
	}
	return 1
}

// ReverseCompare reverses the order of a compare function.
func ReverseCompare[T any](compare func(T, T) int) func(T, T) int {
	return func(a, b T) int {
		return -compare(a, b)
	}
}

// Sort sorts a slice of strings.
func Sort[T constraints.Ordered](slice []T) {
	slices.Sort(slice)
}

// SortReverse sorts a slice of strings in reverse order.
func SortReverse[T constraints.Ordered](slice []T) {
	slices.SortFunc(slice, ReverseCompare(CompareOrdered[T]))
}

// IsSorted returns true if the slice is sorted.
func IsSorted[T constraints.Ordered](slice []T) bool {
	return slices.IsSorted(slice)
}

// IsReverseSorted returns true if the slice is sorted in reverse order.
func IsReverseSorted[T constraints.Ordered](slice []T) bool {
	return slices.IsSortedFunc(slice, ReverseCompare(CompareOrdered[T]))
}

// Uniq returns a slice with all duplicate elements removed. It sorts the slice
// before doing so.
func Uniq[T constraints.Ordered](slice []T) []T {
	Sort(slice)

	cursor := 0
	for i := 1; i < len(slice); i++ {
		if slice[i] != slice[i-1] {
			cursor++
			slice[cursor] = slice[i]
		}
	}

	return slice[:cursor+1]
}

// IsUniq returns true if the slice has no duplicate elements.
func IsUniq[T constraints.Ordered](slice []T) bool {
	// See aocutil_test.go's BenchmarkIsUniq. We pick 100 because it's the sweet
	// spot. Beyond 100, allocating a map is faster.
	const bruteforceThreshold = 100
	var z T
	if len(slice)*int(unsafe.Sizeof(z)) < bruteforceThreshold {
		return isUniqBruteforce(slice)
	}
	return isUniqSet(slice)
}

func isUniqBruteforce[T comparable](slice []T) bool {
	for i, v := range slice {
		if Contains(slice[i+1:], v) {
			return false
		}
	}
	return true
}

func isUniqSet[T comparable](slice []T) bool {
	set := make(map[T]struct{}, len(slice))
	for _, v := range slice {
		_, ok := set[v]
		if ok {
			return false
		}
		set[v] = struct{}{}
	}
	return true
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
func NewBytesScanner[T []byte | string](b T) *Scanner {
	var r io.Reader
	switch b := any(b).(type) {
	case []byte:
		r = bytes.NewReader(b)
	case string:
		r = strings.NewReader(b)
	}
	return NewScanner(r)
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

// Pair is a pair of values.
type Pair[K comparable, V any] struct {
	K K
	V V
}

// MapPairs converts a map to a slice of pairs. The order of the pairs is
// undefined.
func MapPairs[K comparable, V any](m map[K]V) []Pair[K, V] {
	pairs := make([]Pair[K, V], 0, len(m))
	for k, v := range m {
		pairs = append(pairs, Pair[K, V]{K: k, V: v})
	}
	return pairs
}

type prefixedWriter struct {
	w io.Writer
	p []byte
}

// PrefixedWriter returns a writer that prefixes each line with the given
// prefix.
func PrefixedWriter(w io.Writer, prefix string) io.Writer {
	return &prefixedWriter{
		w: w,
		p: []byte(prefix),
	}
}

// PrefixedStdout returns a writer to stdout with the given prefix.
func PrefixedStdout(prefix string) io.Writer {
	return PrefixedWriter(os.Stdout, prefix)
}

func (w *prefixedWriter) Write(b []byte) (int, error) {
	return writePrefixed(w.w, b, w.p)
}

func writePrefixed(dst io.Writer, b, prefix []byte) (int, error) {
	var total int
	var prefixed bool
	for _, line := range bytes.SplitAfter(b, []byte("\n")) {
		if len(line) > 0 && !prefixed {
			_, err := dst.Write(prefix)
			if err != nil {
				return 0, err
			}
			prefixed = true
		}

		n, err := dst.Write(line)
		if err != nil {
			return total, err
		}
		total += n

		if bytes.HasSuffix(line, []byte("\n")) {
			prefixed = false
		}
	}

	return total, nil
}

// Clone deep-copies in and returns a newly-allocated value.
func Clone[T any](in T) T {
	return deepcopy.Copy(in).(T)
}

// GCD returns the greatest common divisor of a and b.
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// LCM returns the least common multiple of the given integers.
// If only one integer is given, it is returned.
func LCM(integers ...int) int {
	if len(integers) < 2 {
		return integers[0]
	}

	a := integers[0]
	b := integers[1]
	integers = integers[2:]

	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}

// Combinations returns all combinations of k elements from the given slice.
func Combinations[T any](slice []T, k int) Iter[[]T] {
	gen := combin.NewCombinationGenerator(len(slice), k)
	dst := make([]int, k)
	buf := make([]T, k)
	return func(yield func([]T) bool) {
		for gen.Next() {
			combinations := gen.Combination(dst)
			for i, j := range combinations {
				buf[i] = slice[j]
			}
			if !yield(buf[:len(combinations)]) {
				return
			}
		}
	}
}

// PairCombinations returns all combinations of pairs from the given slice.
func PairCombinations[T any](slice []T) Iter2[T, T] {
	return func(yield func(T, T) bool) {
		for pair := range Combinations(slice, 2) {
			if !yield(pair[0], pair[1]) {
				return
			}
		}
	}
}

// Permutations returns all permutations of k elements from the given slice.
func Permutations[T any](slice []T, k int) Iter[[]T] {
	gen := combin.NewPermutationGenerator(len(slice), k)
	dst := make([]int, k)
	buf := make([]T, k)
	return func(yield func([]T) bool) {
		for gen.Next() {
			combinations := gen.Permutation(dst)
			for i, j := range combinations {
				buf[i] = slice[j]
			}
			if !yield(buf[:len(combinations)]) {
				return
			}
		}
	}
}

// WaitForKeypress waits for a keypress.
func WaitForKeypress() byte {
	fmt.Print("Press any key to continue...")
	var b [1]byte
	E2(os.Stdin.Read(b[:]))
	fmt.Println()
	return b[0]
}

// PositiveMod returns the positive modulus of k and n.
// It works if k is both positive and negative.
func PositiveMod[T constraints.Integer](k, n T) T {
	return (k%n + n) % n
}

// Polyfit returns the coefficients of the polynomial of degree degree that
// fits the points (xs[i], ys[i]) for i = 0, ..., len(xs) - 1.
func Polyfit(xs, ys []float64, degree int) fit.PolynomialRegressionResult {
	return fit.PolynomialRegression(xs, ys, nil, degree)
}

// RoundedRegression rounds the coefficients of the given polynomial
// regression result.
func RoundedRegression(res fit.PolynomialRegressionResult) fit.PolynomialRegressionResult {
	res.Coefficients = Map(res.Coefficients, math.Round)
	return res
}

// CalculateRegression casts the coefficients of the given polynomial
// regression result to any signed type and calculates the value of the
// polynomial at x. This function is useful when the coefficients are
// calculated using floating-point arithmetic and the result needs to be
// calculated using integer arithmetic.
func CalculateRegression[T constraints.Signed](res fit.PolynomialRegressionResult, x T) T {
	coeffs := Map(res.Coefficients, func(f float64) T { return T(f) })
	y := coeffs[0]
	xp := x
	for _, c := range coeffs[1:] {
		y += xp * c
		xp *= x
	}
	return y
}
