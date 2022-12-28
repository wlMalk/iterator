package iterator

import (
	"testing"
)

func TestUnfold(t *testing.T) {
	a := Unfold([]uint{0, 1}, func(index uint, s []uint) (uint, []uint, bool, error) {
		return s[0], []uint{s[1], s[0] + s[1]}, true, nil
	})
	a = Limit[uint](10)(a)
	checkIteratorEqual(t, a, []uint{0, 1, 1, 2, 3, 5, 8, 13, 21, 34})
}
