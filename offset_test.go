package iterator

import (
	"testing"
)

func TestOffset(t *testing.T) {
	cases := []struct {
		iter     Iterator[int]
		expected []int
	}{
		{Offset[int](5)(Empty[int]()), []int{}},
		{Offset[int](1)(Once(1)), []int{}},
		{Offset[int](5)(Once(1)), []int{}},
		{Limit[int](5)(Offset[int](5)(Ascending(0, 1))), []int{5, 6, 7, 8, 9}},
	}

	for i := range cases {
		checkIteratorEqual(t, cases[i].iter, cases[i].expected)
	}
}
