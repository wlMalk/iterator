package iterator

import (
	"testing"
)

func TestGroup(t *testing.T) {
	iter := FromSlice([]int{1, 2, 3, 1, 2, 3, 1, 2, 3})
	grouped := Group(iter)
	flattened := Flatten(grouped)

	checkIteratorEqual(t, flattened, []int{1, 1, 1, 2, 2, 2, 3, 3, 3})
}

func TestGroupFunc(t *testing.T) {
	iter := Range(0, 9, 1)
	grouped := GroupFunc(func(_ uint, item int) (bool, error) { return item%2 == 0, nil })(iter)
	flattened := Flatten(grouped)

	checkIteratorEqual(t, flattened, []int{0, 2, 4, 6, 8, 1, 3, 5, 7, 9})
}

func TestDuplicates(t *testing.T) {
	iter := Duplicates(FromSlice([]int{1, 1, 2, 1, 3, 4, 3}))

	checkIteratorEqual(t, iter, []int{1, 3})
}

func TestUniques(t *testing.T) {
	iter := Uniques(FromSlice([]int{1, 1, 2, 1, 3, 4, 3}))

	checkIteratorEqual(t, iter, []int{2, 4})
}
