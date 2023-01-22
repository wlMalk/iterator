package iterator

import (
	"testing"
)

func TestSplit(t *testing.T) {
	cases := []struct {
		iter     Iterator[[]int]
		expected [][]int
	}{
		{ToSlices(Split(0)(Empty[int]())), [][]int{}},
		{ToSlices(Split(0)(FromSlice([]int{1, 2, 3, 1, 2, 3}))), [][]int{{1, 2, 3, 1, 2, 3}}},
		{
			ToSlices(
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
		{ToSlices(SplitLeading(0)(Empty[int]())), [][]int{}},
		{
			ToSlices(
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
		{ToSlices(SplitLeading(0)(Empty[int]())), [][]int{}},
		{
			ToSlices(
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
		{ToSlices(Chunk[int](0)(Empty[int]())), [][]int{}},
		{ToSlices(Chunk[int](0)(FromSlice([]int{1, 2, 3, 1, 2, 3}))), [][]int{{1, 2, 3, 1, 2, 3}}},
		{
			ToSlices(
				Chunk[int](3)(FromSlice([]int{1, 2, 3, 1, 2, 3})),
			),
			[][]int{{1, 2, 3}, {1, 2, 3}},
		},
	}

	for i := range cases {
		checkIteratorSliceEqual(t, cases[i].iter, cases[i].expected)
	}
}
