package aocutil

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

func TestSlidingWindow(t *testing.T) {
	var num = []int{0, 1, 2, 3, 4}
	var res = [][]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}}
	var turn int

	SlidingWindow(num, 2, func(w []int) bool {
		if !reflect.DeepEqual(res[turn], w) {
			t.Errorf("mismatch at %d: %v", turn, w)
		}
		turn++
		return false
	})
}

func TestChunk(t *testing.T) {
	var num = []int{0, 1, 2, 3, 4}
	var res = [][]int{{0, 1}, {2, 3}, {4}}
	var turn int

	Chunk(num, 2, func(w []int) {
		if !reflect.DeepEqual(res[turn], w) {
			t.Errorf("mismatch at %d: %v", turn, w)
		}
		turn++
	})
}

func TestHeap(t *testing.T) {
	type test struct {
		name string
		opts HeapOpts[int]
		in   []int
		out  []int
	}

	tests := []test{
		{
			name: "min_preorder",
			opts: HeapOpts[int]{},
			in:   []int{1, 2, 3, 4, 5},
			out:  []int{1, 2, 3, 4, 5},
		},
		{
			name: "min_reverse",
			opts: HeapOpts[int]{},
			in:   []int{5, 4, 3, 2, 1},
			out:  []int{1, 2, 3, 4, 5},
		},
		{
			name: "min_random",
			opts: HeapOpts[int]{},
			in:   []int{3, 5, 1, 4, 2},
			out:  []int{1, 2, 3, 4, 5},
		},
		{
			name: "max_preorder",
			opts: HeapOpts[int]{Max: true},
			in:   []int{5, 4, 3, 2, 1},
			out:  []int{5, 4, 3, 2, 1},
		},
		{
			name: "max_reverse",
			opts: HeapOpts[int]{Max: true},
			in:   []int{1, 2, 3, 4, 5},
			out:  []int{5, 4, 3, 2, 1},
		},
		{
			name: "max_random",
			opts: HeapOpts[int]{Max: true},
			in:   []int{3, 5, 1, 4, 2},
			out:  []int{5, 4, 3, 2, 1},
		},
		{
			name: "min_limit",
			opts: HeapOpts[int]{Limit: 3},
			in:   []int{3, 5, 1, 4, 2},
			out:  []int{3, 2, 1},
		},
		{
			name: "max_limit",
			opts: HeapOpts[int]{Limit: 3, Max: true},
			in:   []int{3, 5, 1, 4, 2},
			out:  []int{3, 4, 5},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := NewHeapOpts(test.opts)
			for _, v := range test.in {
				h.Push(v)
			}

			h.Sort()
			t.Logf("heap: %v", h.ToSlice())

			for _, v := range test.out {
				if h.Pop() != v {
					t.Errorf("mismatch: %v", h)
				}
			}
		})
	}
}

func TestSort(t *testing.T) {
	type test struct {
		in  []int
		out []int
		rev bool
	}

	tests := []test{
		{[]int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}, false},
		{[]int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}, false},
		{[]int{3, 5, 1, 4, 2}, []int{1, 2, 3, 4, 5}, false},
		{[]int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}, true},
		{[]int{5, 4, 3, 2, 1}, []int{5, 4, 3, 2, 1}, true},
		{[]int{3, 5, 1, 4, 2}, []int{5, 4, 3, 2, 1}, true},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
			input := append([]int(nil), test.in...)
			if test.rev {
				SortReverse(input)
			} else {
				Sort(input)
			}
			if !reflect.DeepEqual(input, test.out) {
				t.Errorf("unexpected Sort(%#v):\n"+
					"got    %#v\n"+
					"expect %#v",
					test.in, input, test.out)
			}
		})
	}
}

func TestUniq(t *testing.T) {
	type test struct {
		in  []int
		out []int
	}

	tests := []test{
		{[]int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{[]int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
		{[]int{1, 1, 2, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{[]int{5, 4, 3, 2, 1, 1, 1}, []int{1, 2, 3, 4, 5}},
		{[]int{1, 2, 3, 4, 5, 5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
			input := append([]int(nil), test.in...)
			input = Uniq(input)
			if !reflect.DeepEqual(input, test.out) {
				t.Errorf("unexpected Uniq(%#v):\n"+
					"got    %#v\n"+
					"expect %#v",
					test.in, input, test.out)
			}
		})
	}
}

func TestTrim(t *testing.T) {
	testTrim(t, trimLeft, []byte{0, 0, 1, 2, 3}, []byte{1, 2, 3}, 0)
	testTrim(t, trimRight, []byte{1, 2, 3, 0, 0}, []byte{1, 2, 3}, 0)
	testTrim(t, trimBoth, []byte{0, 1, 2, 3, 0}, []byte{1, 2, 3}, 0)
	testTrim(t, trimBoth, []byte{1, 2, 3}, []byte{1, 2, 3}, 0)

	testTrim(t, trimLeft, []int{0, 0, 1, 2, 3}, []int{1, 2, 3}, 0)
	testTrim(t, trimRight, []int{1, 2, 3, 0, 0}, []int{1, 2, 3}, 0)
	testTrim(t, trimBoth, []int{0, 1, 2, 3, 0}, []int{1, 2, 3}, 0)
	testTrim(t, trimBoth, []int{1, 2, 3}, []int{1, 2, 3}, 0)
}

func testTrim[T comparable](t *testing.T, trimType trimType, in, expect []T, val T) {
	t.Helper()
	t.Run(fmt.Sprintf("%T-%s", val, trimType), func(t *testing.T) {
		got := trim(in, val, trimType)
		if !reflect.DeepEqual(got, expect) {
			t.Fatalf(
				"unexpected trim(%v, %v, %v):\n"+
					"got    %#v\n"+
					"expect %#v",
				in, val, trimType,
				got, expect,
			)
		}
	})
}

func TestCount(t *testing.T) {
	type test struct {
		arr []int
		val int
		out int
	}

	tests := []test{
		{[]int{1, 2, 3, 4, 5}, 5, 1},
		{[]int{5, 4, 3, 2, 1}, 5, 1},
		{[]int{1, 2, 3, 4, 5}, 1, 1},
		{[]int{5, 4, 3, 2, 1}, 1, 1},
		{[]int{1, 2, 3, 4, 5}, 0, 0},
		{[]int{5, 4, 3, 2, 1}, 0, 0},
		{[]int{1, 2, 3, 4, 5, 5}, 5, 2},
		{[]int{5, 4, 3, 2, 1, 1}, 1, 2},
		{[]int{1, 2, 3, 4, 5, 5, 5}, 5, 3},
		{[]int{5, 4, 3, 2, 1, 1, 1}, 1, 3},
		{[]int{1, 2, 3, 4, 5, 5, 5}, 1, 1},
		{[]int{5, 4, 3, 2, 1, 1, 1}, 5, 1},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
			got := Count(test.arr, test.val)
			if got != test.out {
				t.Errorf("unexpected Count(%#v, %v):\n"+
					"got    %#v\n"+
					"expect %#v",
					test.arr, test.val, got, test.out)
			}
		})
	}
}

func TestIsUniq(t *testing.T) {
	type test struct {
		in  []int
		out bool
	}

	tests := []test{
		{[]int{1, 2, 3, 4, 5}, true},
		{[]int{5, 4, 3, 2, 1}, true},
		{[]int{1, 1, 2, 2, 3, 4, 5}, false},
		{[]int{5, 4, 3, 2, 1, 1, 1}, false},
	}

	type fn struct {
		name string
		fn   func([]int) bool
	}

	fns := []fn{
		{"auto", IsUniq[int]},
		{"bruteforce", isUniqBruteforce[int]},
		{"set", isUniqSet[int]},
	}

	for _, fn := range fns {
		for i, test := range tests {
			t.Run(fmt.Sprintf("%s_%d", fn.name, i), func(t *testing.T) {
				if got := fn.fn(test.in); got != test.out {
					t.Errorf("unexpected isUniq(%#v):\n"+
						"got    %#v\n"+
						"expect %#v",
						test.in, got, test.out)
				}
			})
		}
	}
}

func BenchmarkIsUniq(b *testing.B) {
	sizes := []int{50, 100, 200, 300, 400, 500, 1_000, 5_000}
	mrand := rand.New(rand.NewSource(0))

	samples := make([]int, sizes[len(sizes)-1])
	for i := range samples {
		samples[i] = mrand.Int()
	}

	benchmark := func(name string, f func([]int) bool) {
		for _, size := range sizes {
			b.Run(fmt.Sprintf("%s-%d", name, size), func(b *testing.B) {
				for j := 0; j < b.N; j++ {
					f(samples[:size])
				}
			})
		}
	}

	// BenchmarkIsUniq/bruteforce-50-16     652293       1078.00 ns/op
	// BenchmarkIsUniq/bruteforce-100-16    198391       3188.00 ns/op
	// BenchmarkIsUniq/bruteforce-200-16     49496      12972.00 ns/op
	// BenchmarkIsUniq/bruteforce-300-16     24886      28768.00 ns/op
	// BenchmarkIsUniq/bruteforce-400-16     11942      53370.00 ns/op
	// BenchmarkIsUniq/bruteforce-500-16      6681      94309.00 ns/op
	// BenchmarkIsUniq/bruteforce-1000-16     1509     476896.00 ns/op
	// BenchmarkIsUniq/bruteforce-5000-16       54   10414225.00 ns/op
	// BenchmarkIsUniq/set-50-16            267532       2249.00 ns/op
	// BenchmarkIsUniq/set-100-16           147910       4118.00 ns/op
	// BenchmarkIsUniq/set-200-16            75312       8706.00 ns/op
	// BenchmarkIsUniq/set-300-16            52758      11870.00 ns/op
	// BenchmarkIsUniq/set-400-16            35253      16641.00 ns/op
	// BenchmarkIsUniq/set-500-16            32888      17982.00 ns/op
	// BenchmarkIsUniq/set-1000-16           14428      36395.00 ns/op
	// BenchmarkIsUniq/set-5000-16            3145     196633.00 ns/op
	benchmark("auto", IsUniq[int])
	benchmark("bruteforce", isUniqBruteforce[int])
	benchmark("set", isUniqSet[int])
}
