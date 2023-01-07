package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type iterator[T any] interface {
	Next() bool
	Get() (T, error)
	Close() error
	Err() error
}

func CheckIteratorEqual[T any](t *testing.T, iter iterator[T], items []T) {
	defer func() {
		iter.Close()
	}()

	i := 0
	for iter.Next() {
		require.Less(t, i, len(items))

		item, err := iter.Get()
		require.NoError(t, err)
		assert.Equal(t, items[i], item)
		i++
	}
	assert.Equal(t, len(items), i)
	require.NoError(t, iter.Err())
	require.NoError(t, iter.Close())
}
