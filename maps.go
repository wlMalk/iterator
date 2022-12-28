package iterator

// Keys is a modifier that takes a map iterator and returns an iterator for the keys
func Keys[K comparable, V any](iter Iterator[KV[K, V]]) Iterator[K] {
	return Map(func(_ uint, item KV[K, V]) (K, error) {
		return item.Key, nil
	})(iter)
}

// Values is a modifier that takes a map iterator and returns an iterator for the values
func Values[K comparable, V any](iter Iterator[KV[K, V]]) Iterator[V] {
	return Map(func(_ uint, item KV[K, V]) (V, error) {
		return item.Val, nil
	})(iter)
}

// Assoc the given keys to the given values in a new map iterator
func Assoc[K comparable, V any](keys Iterator[K], values Iterator[V]) Iterator[KV[K, V]] {
	var hasMore bool
	var item KV[K, V]
	var err error

	return &iterator[KV[K, V]]{
		next: func() bool {
			hasMore = keys.Next() && values.Next()
			if !hasMore {
				return hasMore
			}

			var key K
			key, err = keys.Get()
			if err != nil {
				return false
			}

			var value V
			value, err = values.Get()
			if err != nil {
				return false
			}

			item = KV[K, V]{Key: key, Val: value}

			return true
		},
		get: func() (KV[K, V], error) {
			if !hasMore || err != nil {
				return *new(KV[K, V]), err
			}

			return item, nil
		},
		close: func() error {
			if err := keys.Close(); err != nil {
				values.Close()
				return err
			}
			return values.Close()
		},
		err: func() error {
			if err != nil {
				return err
			}
			if err := keys.Err(); err != nil {
				return err
			}
			return values.Err()
		},
	}
}

// AssocKeys returns a modifier that associates the iterator to the given values
func AssocKeys[K comparable, V any](values Iterator[V]) Modifier[K, KV[K, V]] {
	return func(keys Iterator[K]) Iterator[KV[K, V]] {
		return Assoc(keys, values)
	}
}

// AssocValues returns a modifier that associates the iterator to the given keys
func AssocValues[K comparable, V any](keys Iterator[K]) Modifier[V, KV[K, V]] {
	return func(values Iterator[V]) Iterator[KV[K, V]] {
		return Assoc(keys, values)
	}
}
