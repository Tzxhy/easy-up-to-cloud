package models

type Store = map[string]interface{}

var store = make(Store)

func GetStore() Store {
	return store
}

type MemStore interface {
	GetKey(key string) interface{}
	SetKey(key string, val interface{})
}

func GetKey(k string) interface{} {
	s := GetStore()
	V, ok := s[k]
	if ok {
		return V
	}
	return nil
}

func SetKey(k string, val interface{}) {
	store := GetStore()
	store[k] = val
}

func ClearKey(k string) {
	store := GetStore()
	delete(store, k)
}
