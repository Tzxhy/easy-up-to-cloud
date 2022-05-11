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

func Has[T comparable](slice []T, has T) bool {
	for _, item := range slice {

		if has == item {
			return true
		}
	}
	return false
}

func HasByFunc[T interface{}](slice []T, is func(T) bool) bool {
	for _, item := range slice {
		_has := is(item)
		if _has {
			return true
		}
	}
	return false
}

// 使切片中的元素都唯一
func Unique[T comparable](slice []T) *[]T {
	var newSlice []T
	var _map = make(map[T]int, 0)
	for _, item := range slice {
		_, has := _map[item]
		if !has {
			newSlice = append(newSlice, item)
			_map[item] = 1
		}
	}
	return &newSlice
}