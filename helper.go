package gls

import "unsafe"

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

type arg struct {
	sz int32
	fp unsafe.Pointer
	fn func()
	ls map[string]interface{}
	ca Cache
}

func newProc(fn func(), ls map[string]interface{}, ca Cache) {
	fp := proc
	sz := uint32(unsafe.Sizeof(arg{}))
	ag := &arg{
		sz: int32(uintptr(sz) - 2*unsafe.Sizeof(fn)),
		fp: *(*unsafe.Pointer)(unsafe.Pointer(&fp)),
		fn: fn,
		ls: ls,
		ca: ca,
	}
	np := newproc
	reflectcall(nil, *(*unsafe.Pointer)(unsafe.Pointer(&np)), unsafe.Pointer(ag), sz, sz)
}

func proc(fn func(), ls map[string]interface{}, ca Cache) {
	defer ca.Clr()
	ca.Put(ls)
	fn()
}

//go:linkname newproc runtime.newproc
func newproc(siz int32, fn unsafe.Pointer)

//go:linkname reflectcall runtime.reflectcall
func reflectcall(argtype unsafe.Pointer, fn, arg unsafe.Pointer, argsize uint32, retoffset uint32)
