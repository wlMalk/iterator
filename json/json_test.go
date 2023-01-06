package json

import (
	"encoding/json"
	"testing"

	"github.com/wlMalk/iterator"
	"github.com/wlMalk/iterator/internal/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkIteratorEqual[T comparable](t *testing.T, iter Iterator[T], items []T) {
	utils.CheckIteratorEqual[T](t, iter, items)
}

func TestIteratorMarshalJSON(t *testing.T) {
	it := Iterator[int]{iterator.Range(1, 5, 1)}
	b, err := json.Marshal(it)
	require.NoError(t, err)
	assert.Equal(t, []byte("[1,2,3,4,5]"), b)
}

func TestIteratorUnmarshalJSON(t *testing.T) {
	var it Iterator[int]
	err := json.Unmarshal([]byte("[1,2,3,4,5]"), &it)
	require.NoError(t, err)
	checkIteratorEqual(t, it, []int{1, 2, 3, 4, 5})
}
