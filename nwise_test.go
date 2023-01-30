package iterator

import (
	"testing"
)

func TestNwise(t *testing.T) {
	cases := []struct {
		iter     Iterator[int]
		size     int
		expected [][]int
	}{
		{Empty[int](), 3, [][]int{}},
		{Once(1), 3, [][]int{}},
		{Range(1, 6, 1), 3, [][]int{{1, 2, 3}, {2, 3, 4}, {3, 4, 5}, {4, 5, 6}}},
	}

	for i := range cases {
		checkIteratorSliceEqual(t, ToSlices(Nwise[int](cases[i].size)(cases[i].iter)), cases[i].expected)
	}
}
