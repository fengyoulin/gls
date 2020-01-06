package gls

import (
	"github.com/fengyoulin/goid"
	"runtime"
	"sync"
)

// Cache is the gls cache interface
type Cache interface {
	// All returns all the key/values in current goroutine's local storage
	All() (kvs map[string]interface{}, ok bool)
	// Clr clears the current goroutine's local storage
	Clr()
	// Del deletes a key in current goroutine's local storage
	Del(key string)
	// Get get the value of a key in current goroutine's local storage
	Get(key string) (val interface{}, ok bool)
	// Set set the value to a key in current goroutine's local storage
	Set(key string, val interface{})
}

type single struct {
	lock sync.RWMutex
	data map[int64]map[string]interface{}
}

func (s *single) All() (kvs map[string]interface{}, ok bool) {
	var m map[string]interface{}
	s.lock.RLock()
	m, ok = s.data[goid.ID()]
	s.lock.RUnlock()
	if !ok {
		return
	}
	if len(m) == 0 {
		return
	}
	kvs = make(map[string]interface{})
	for key, val := range m {
		kvs[key] = val
	}
	return
}

func (s *single) Clr() {
	s.lock.Lock()
	delete(s.data, goid.ID())
	s.lock.Unlock()
}

func (s *single) Del(key string) {
	s.lock.RLock()
	m, ok := s.data[goid.ID()]
	s.lock.RUnlock()
	if ok {
		delete(m, key)
	}
}

func (s *single) Get(key string) (val interface{}, ok bool) {
	s.lock.RLock()
	m, Ok := s.data[goid.ID()]
	s.lock.RUnlock()
	if Ok {
		val, ok = m[key]
	}
	return
}

func (s *single) Set(key string, val interface{}) {
	id := goid.ID()
	s.lock.RLock()
	m, ok := s.data[id]
	s.lock.RUnlock()
	if !ok {
		m = make(map[string]interface{})
		s.lock.Lock()
		s.data[id] = m
		s.lock.Unlock()
	}
	m[key] = val
}

type sharding []*single

func (s *sharding) All() (kvs map[string]interface{}, ok bool) {
	return s.shard().All()
}

func (s *sharding) Clr() {
	s.shard().Clr()
}

func (s *sharding) Del(key string) {
	s.shard().Del(key)
}

func (s *sharding) Get(key string) (val interface{}, ok bool) {
	return s.shard().Get(key)
}

func (s *sharding) Set(key string, val interface{}) {
	s.shard().Set(key, val)
}

func (s *sharding) shard() *single {
	return (*s)[int(goid.ID())%len(*s)]
}

// New creates a goroutine local storage cache
func New(shardingMode bool) Cache {
	if !shardingMode {
		return &single{
			data: make(map[int64]map[string]interface{}),
		}
	}
	n := runtime.NumCPU()
	c := sharding(make([]*single, n))
	for i := 0; i < n; i++ {
		c[i] = &single{
			data: make(map[int64]map[string]interface{}),
		}
	}
	return &c
}
