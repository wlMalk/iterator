package iterator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLen(t *testing.T) {
	cases := []struct {
		iter     Iterator[int]
		expected uint
	}{
		{Empty[int](), 0},
		{Once(1), 1},
		{Limit[int](5)(Ascending(0, 1)), 5},
	}

	for i := range cases {
		length, err := Len(cases[i].iter)
		require.NoError(t, err)
		assert.Equal(t, cases[i].expected, length)
	}
}

func TestEqual(t *testing.T) {
	cases := []struct {
		iterA    Iterator[int]
		iterB    Iterator[int]
		expected bool
	}{
		{Empty[int](), Empty[int](), true},
		{Empty[int](), Once(1), false},
		{Once(1), Empty[int](), false},
	}

	for i := range cases {
		equal, err := Equal(cases[i].iterA, cases[i].iterB)
		require.NoError(t, err)
		assert.Equal(t, cases[i].expected, equal)
	}
}

func TestOne(t *testing.T) {
	cases := []struct {
		iter     Iterator[int]
		expected int
		err      error
	}{
		{FromSlice([]int{}), 0, ErrNoItems},
		{FromSlice([]int{1, 2}), 0, ErrMultiItems},
		{FromSlice([]int{1}), 1, nil},
	}

	for i := range cases {
		item, err := One(cases[i].iter)
		require.ErrorIs(t, err, cases[i].err)
		assert.Equal(t, cases[i].expected, item)
	}
}
