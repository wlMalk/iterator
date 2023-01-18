package iterator

import (
	"testing"
)

func TestMerge(t *testing.T) {
	a := Range(1, 3, 1)
	b := Range(4, 6, 1)
	c := Range(7, 9, 1)

	merged := Merge(a, b, c)
	checkIteratorEqualUnordered(t, merged, []int{1, 2, 3, 4, 5, 6, 7, 8, 9})
}
