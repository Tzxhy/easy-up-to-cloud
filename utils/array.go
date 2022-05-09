package utils

func Some[T interface{}](slice []T, has func(T) bool) bool {
	for _, item := range slice {
		_has := has(item)
		if _has {
			return true
		}
	}
	return false
}
