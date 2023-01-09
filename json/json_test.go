package json

import (
	"encoding/json"
	"testing"

	it "github.com/wlMalk/iterator"
	"github.com/wlMalk/iterator/internal/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkIteratorEqual[T comparable](t *testing.T, iter it.Iterator[T], items []T) {
	utils.CheckIteratorEqual[T](t, iter, items)
}

func TestIteratorMarshalJSON(t *testing.T) {
	iter := Iterator[int]{it.Range(1, 5, 1)}
	b, err := json.Marshal(iter)
	require.NoError(t, err)
	assert.Equal(t, []byte("[1,2,3,4,5]"), b)
}

func TestIteratorUnmarshalJSON(t *testing.T) {
	var iter *Iterator[int]
	err := json.Unmarshal([]byte("[1,2,3,4,5]"), &iter)
	require.NoError(t, err)
	checkIteratorEqual[int](t, iter, []int{1, 2, 3, 4, 5})
}
