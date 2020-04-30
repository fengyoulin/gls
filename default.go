package gls

import "runtime"

var defaultCache Cache

func init() {
	defaultCache = New(runtime.NumCPU() > 1)
}

// All returns all the key/values in current goroutine's local storage
func All() (kvs map[string]interface{}, ok bool) {
	return defaultCache.All()
}

// Put puts all the key/values into current goroutine's local storage
func Put(kvs map[string]interface{}) {
	defaultCache.Put(kvs)
}

// Clr clears the current goroutine's local storage
func Clr() {
	defaultCache.Clr()
}

// Del deletes a key in current goroutine's local storage
func Del(key string) {
	defaultCache.Del(key)
}

// Get get the value of a key in current goroutine's local storage
func Get(key string) (val interface{}, ok bool) {
	return defaultCache.Get(key)
}

// Set set the value to a key in current goroutine's local storage
func Set(key string, val interface{}) {
	defaultCache.Set(key, val)
}
