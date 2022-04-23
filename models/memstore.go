package models

import (
	"gitee.com/tzxhy/web/middlewares"
)

type MemStore interface {
	GetKey(key string) interface{}
	SetKey(key string, val interface{})
}

func GetKey(k string) interface{} {
	s := middlewares.GetStore()
	V, ok := s[k]
	if ok {
		return V
	}
	return nil
}

func SetKey(k string, val interface{}) {
	store := middlewares.GetStore()
	store[k] = val
}

func ClearKey(k string) {
	store := middlewares.GetStore()
	delete(store, k)
}
