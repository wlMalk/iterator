package iterator

// Pipe applies the modifiers to the given iterator and returns the resulting
// iterator
func Pipe[T any](iter Iterator[T], mods ...Modifier[T, T]) Iterator[T] {
	for i := range mods {
		iter = mods[i](iter)
	}
	return iter
}
