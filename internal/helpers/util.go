package helpers

func If[T any, K any](condition bool, conditionTrue T, conditionFalse K) (T || K) {
	if condition {
		return conditionTrue
	}
	return conditionFalse
}
