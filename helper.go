package gls

// Go start a new goroutine, inherit current goroutine's local storage
func Go(fn func(), c Cache) {
	if c == nil {
		c = defaultCache
	}
	ls, _ := c.All()
	newProc(fn, ls, c)
}

// GoWith start a new goroutine, with the given local storage
func GoWith(fn func(), c Cache, ls map[string]interface{}) {
	if c == nil {
		c = defaultCache
	}
	newProc(fn, ls, c)
}

func newProc(fn func(), ls map[string]interface{}, ca Cache) {
	go func() {
		defer ca.Clr()
		ca.Put(ls)
		fn()
	}()
}
