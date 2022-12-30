package iterator

import (
	"testing"
)

func TestKeys(t *testing.T) {
	a := Keys(FromMap(map[string]int{
		"first":  1,
		"second": 2,
		"third":  3,
		"fourth": 4,
		"fifth":  5,
	}))

	checkIteratorEqualUnordered(t, a, []string{
		"first",
		"second",
		"third",
		"fourth",
		"fifth",
	})
}

func TestValues(t *testing.T) {
	a := Values(FromMap(map[string]int{
		"first":  1,
		"second": 2,
		"third":  3,
		"fourth": 4,
		"fifth":  5,
	}))

	checkIteratorEqualUnordered(t, a, []int{
		1,
		2,
		3,
		4,
		5,
	})
}

func TestAssocKeys(t *testing.T) {
	keys := FromSlice([]string{
		"first",
		"second",
		"third",
		"fourth",
		"fifth",
	})

	values := FromSlice([]int{
		1,
		2,
		3,
		4,
		5,
	})

	checkIteratorMapEqual(t, AssocKeys[string, int](keys)(values), map[string]int{
		"first":  1,
		"second": 2,
		"third":  3,
		"fourth": 4,
		"fifth":  5,
	})
}

func TestAssocValues(t *testing.T) {
	keys := FromSlice([]string{
		"first",
		"second",
		"third",
		"fourth",
		"fifth",
	})

	values := FromSlice([]int{
		1,
		2,
		3,
		4,
		5,
	})

	checkIteratorMapEqual(t, AssocValues[string](values)(keys), map[string]int{
		"first":  1,
		"second": 2,
		"third":  3,
		"fourth": 4,
		"fifth":  5,
	})
}
