package iterator

import (
	"math/rand"
	"testing"
)

func TestRandom(t *testing.T) {
	cases := []struct {
		max      int
		expected []int
	}{
		{10, []int{6, 9, 6, 4, 4, 6, 0, 1, 0, 3}},
	}

	rand.Seed(1)

	for i := range cases {
		checkIteratorEqual(t, Limit[int](10)(Random(cases[i].max)), cases[i].expected)
	}
}

func TestRandomBetween(t *testing.T) {
	cases := []struct {
		min, max int
		expected []int
	}{
		{1, 100, []int{60, 94, 66, 44, 43, 68, 7, 16, 10, 30}},
	}

	rand.Seed(1)

	for i := range cases {
		checkIteratorEqual(t, Limit[int](10)(RandomBetween(cases[i].min, cases[i].max)), cases[i].expected)
	}
}

func TestSamples(t *testing.T) {
	cases := []struct {
		iter               Iterator[int]
		population, sample uint
		expected           []int
	}{
		{Range(1, 100, 1), 100, 10, []int{7, 9, 32, 36, 38, 54, 56, 71, 89, 98}},
	}

	rand.Seed(1)

	for i := range cases {
		checkIteratorEqual(t, Samples[int](cases[i].population, cases[i].sample)(cases[i].iter), cases[i].expected)
	}
}

func TestSamplesFunc(t *testing.T) {
	cases := []struct {
		iter               Iterator[int]
		population, sample uint
		fn                 func() float64
		expected           []int
	}{
		{Range(1, 100, 1), 100, 10, func() float64 { return 0.5 }, []int{82, 84, 86, 88, 90, 92, 94, 96, 98, 100}},
	}

	rand.Seed(1)

	for i := range cases {
		checkIteratorEqual(t, SamplesFunc[int](cases[i].population, cases[i].sample, cases[i].fn)(cases[i].iter), cases[i].expected)
	}
}
