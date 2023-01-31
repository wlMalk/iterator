package iterator

import (
	"testing"
)

func TestUnfold(t *testing.T) {
	a := Unfold([]int{0, 1}, func(index int, s []int) (int, []int, bool, error) {
		return s[0], []int{s[1], s[0] + s[1]}, true, nil
	})
	a = Limit[int](10)(a)
	checkIteratorEqual(t, a, []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34})
}
