package iterator

import (
	"testing"
)

func TestSplit(t *testing.T) {
	cases := []struct {
		iter     Iterator[[]int]
		expected [][]int
	}{
		{ToSlices[int]()(Split(0)(Empty[int]())), [][]int{}},
		{ToSlices[int]()(Split(0)(FromSlice([]int{1, 2, 3, 1, 2, 3}))), [][]int{{1, 2, 3, 1, 2, 3}}},
		{
			ToSlices[int]()(
				Split(3)(FromSlice([]int{1, 2, 3, 1, 2, 3})),
			),
			[][]int{{1, 2}, {1, 2}},
		},
	}

	for i := range cases {
		checkIteratorSliceEqual(t, cases[i].iter, cases[i].expected)
	}
}

func TestSplitLeading(t *testing.T) {
	cases := []struct {
		iter     Iterator[[]int]
		expected [][]int
	}{
		{ToSlices[int]()(SplitLeading(0)(Empty[int]())), [][]int{}},
		{
			ToSlices[int]()(
				SplitLeading(3)(FromSlice([]int{1, 2, 3, 1, 2, 3})),
			),
			[][]int{{1, 2, 3}, {1, 2, 3}},
		},
	}

	for i := range cases {
		checkIteratorSliceEqual(t, cases[i].iter, cases[i].expected)
	}
}

func TestSplitTrailing(t *testing.T) {
	cases := []struct {
		iter     Iterator[[]int]
		expected [][]int
	}{
		{ToSlices[int]()(SplitLeading(0)(Empty[int]())), [][]int{}},
		{
			ToSlices[int]()(
				SplitTrailing(3)(FromSlice([]int{1, 2, 3, 1, 2, 3})),
			),
			[][]int{{1, 2}, {3, 1, 2}, {3}},
		},
	}

	for i := range cases {
		checkIteratorSliceEqual(t, cases[i].iter, cases[i].expected)
	}
}

func TestChunk(t *testing.T) {
	cases := []struct {
		iter     Iterator[[]int]
		expected [][]int
	}{
		{ToSlices[int]()(Chunk[int](0)(Empty[int]())), [][]int{}},
		{ToSlices[int]()(Chunk[int](0)(FromSlice([]int{1, 2, 3, 1, 2, 3}))), [][]int{{1, 2, 3, 1, 2, 3}}},
		{
			ToSlices[int]()(
				Chunk[int](3)(FromSlice([]int{1, 2, 3, 1, 2, 3})),
			),
			[][]int{{1, 2, 3}, {1, 2, 3}},
		},
	}

	for i := range cases {
		checkIteratorSliceEqual(t, cases[i].iter, cases[i].expected)
	}
}
