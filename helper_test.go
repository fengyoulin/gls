package gls_test

import (
	"github.com/fengyoulin/gls"
	"sync"
	"testing"
)

func TestGo(t *testing.T) {
	const limit = 10
	gls.Set("message", "here, gls!")
	var wg sync.WaitGroup
	wg.Add(limit)
	var m sync.Map
	fn := func() {
		defer wg.Done()
		ls, ok := gls.All()
		if !ok {
			t.Error("local storage not found")
		}
		cu, ok := ls["current"]
		if !ok {
			t.Error(`"current" not found`)
		}
		if _, ok = m.LoadOrStore(cu, true); ok {
			t.Errorf(`"current" duplicate: %v`, cu)
		}
		_, ok = ls["message"]
		if !ok {
			t.Error(`"message" not found`)
		}
	}
	for i := 0; i < limit; i++ {
		gls.Set("current", i)
		gls.Go(fn, nil)
	}
	wg.Wait()
}

func TestGoWith(t *testing.T) {
	const limit = 10
	gls.Set("message", "here, gls!")
	var wg sync.WaitGroup
	wg.Add(limit)
	var m sync.Map
	fn := func() {
		defer wg.Done()
		ls, ok := gls.All()
		if !ok {
			t.Error("local storage not found")
		}
		cu, ok := ls["current"]
		if !ok {
			t.Error(`"current" not found`)
		}
		if _, ok = m.LoadOrStore(cu, true); ok {
			t.Errorf(`"current" duplicate: %v`, cu)
		}
		_, ok = ls["message"]
		if ok {
			t.Error(`got "message"`)
		}
	}
	for i := 0; i < limit; i++ {
		gls.GoWith(fn, nil, map[string]interface{}{"current": i})
	}
	wg.Wait()
}
