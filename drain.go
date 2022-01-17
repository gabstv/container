package container

// Drain fetches all buffered items from a channel and returns them as a slice.
func Drain[T any](ch <-chan T) []T {
	var ret []T
	for {
		select {
		case item, ok := <-ch:
			if !ok {
				return ret
			}
			ret = append(ret, item)
		default:
			return ret
		}
	}
}
