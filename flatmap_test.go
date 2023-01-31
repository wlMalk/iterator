package iterator

import (
	"strings"
	"testing"
)

func TestFlatMap(t *testing.T) {
	cases := []struct {
		iter     Iterator[string]
		expected []string
	}{
		{Pipe(
			FromSlice([]string{"Hello", "World"}),
			FlatMap(func(_ int, _ int, str string) (Iterator[string], error) {
				return FromSlice(strings.Split(str, "")), nil
			}),
		), []string{"H", "e", "l", "l", "o", "W", "o", "r", "l", "d"}},
	}

	for i := range cases {
		checkIteratorEqual(t, cases[i].iter, cases[i].expected)
	}
}
