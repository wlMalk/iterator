package csv

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wlMalk/iterator"
	"github.com/wlMalk/iterator/internal/utils"
)

func checkIteratorEqual[T any](t *testing.T, iter iterator.Iterator[T], items []T) {
	utils.CheckIteratorEqual[T](t, iter, items)
}

func TestRead(t *testing.T) {
	r := bytes.NewReader([]byte("a,b\n1,1\n2,2\n3,3\n"))
	reader := Read[string](r)
	reader.ExpectHeader(true)

	header, err := reader.Header()
	require.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, header)

	checkIteratorEqual[[]string](t, reader, [][]string{{"1", "1"}, {"2", "2"}, {"3", "3"}})
}

func TestWrite(t *testing.T) {
	iter := iterator.Map(
		func(_ int, item string) ([]string, error) {
			return []string{item, item}, nil
		},
	)(
		iterator.Strings(
			iterator.Range(1, 3, 1),
		),
	)

	var buf bytes.Buffer
	err := Write(&buf, iter).Write()
	b := buf.Bytes()

	require.NoError(t, err)
	assert.Equal(t, []byte("1,1\n2,2\n3,3\n"), b)
}

// func TestIteratorUnmarshalJSON(t *testing.T) {
// 	var it Iterator[int]
// 	err := json.Unmarshal([]byte("[1,2,3,4,5]"), &it)
// 	require.NoError(t, err)
// 	checkIteratorEqual(t, it, []int{1, 2, 3, 4, 5})
// }
