package iterator

import (
	"testing"
)

func TestFilter(t *testing.T) {
	a := Filter(func(_ uint, item int) (bool, error) {
		return item%2 == 0, nil
	})(Range(0, 10, 1))
	checkIteratorEqual(t, a, []int{0, 2, 4, 6, 8, 10})
}

func TestRemove(t *testing.T) {
	a := Remove(1)(FromSlice([]int{1, 1, 2, 2}))
	checkIteratorEqual(t, a, []int{2, 2})
}

func TestRemoveFunc(t *testing.T) {
	a := RemoveFunc(func(_ uint, item int) (bool, error) {
		return item%2 == 0, nil
	})(Range(0, 10, 1))
	checkIteratorEqual(t, a, []int{1, 3, 5, 7, 9})
}
