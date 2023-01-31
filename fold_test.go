package iterator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReduce(t *testing.T) {
	cases := []struct {
		iter     Iterator[int]
		expected int
		err      error
	}{
		{Empty[int](), 0, ErrNoItems},
		{Once(1), 1, nil},
		{Limit[int](5)(Ascending(0, 1)), 10, nil},
	}

	for i := range cases {
		sum, err := Reduce(cases[i].iter, func(_ int, item, sum int) (int, error) {
			return sum + item, nil
		})
		require.ErrorIs(t, err, cases[i].err)
		assert.Equal(t, cases[i].expected, sum)
	}
}

func TestSum(t *testing.T) {
	cases := []struct {
		iter     Iterator[int]
		expected int
	}{
		{Empty[int](), 0},
		{Once(1), 1},
		{Limit[int](5)(Ascending(0, 1)), 10},
	}

	for i := range cases {
		sum, err := Sum(cases[i].iter)
		require.NoError(t, err)
		assert.Equal(t, cases[i].expected, sum)
	}
}

func TestString(t *testing.T) {
	cases := []struct {
		iter     Iterator[int]
		expected string
	}{
		{Empty[int](), ""},
		{Once(1), "1"},
		{Limit[int](5)(Ascending(0, 1)), "01234"},
	}

	for i := range cases {
		str, err := String(cases[i].iter)
		require.NoError(t, err)
		assert.Equal(t, cases[i].expected, str)
	}
}
