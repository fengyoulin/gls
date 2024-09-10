package gls_test

import (
	"github.com/fengyoulin/gls"
	"strconv"
	"sync/atomic"
	"testing"
)

const counterStart = 123456789

var (
	singleCache   gls.Cache
	shardingCache gls.Cache
	counter       atomicInteger
)

func init() {
	singleCache = gls.New(false)
	shardingCache = gls.New(true)
}

func TestSingle(t *testing.T) {
	testCache(singleCache, t)
}

func TestSharding(t *testing.T) {
	testCache(shardingCache, t)
}

func testCache(c gls.Cache, t *testing.T) {
	counter.Set(counterStart)
	for i := 0; i < 1000; i++ {
		str := counter.Next().String()
		_, ok := c.Get(str)
		if ok {
			t.Errorf("unexpected key: %s\n", str)
		}
		c.Set(str, str)
		v, ok := c.Get(str)
		if !ok {
			t.Errorf("key: %s not found\n", str)
		}
		s, ok := v.(string)
		if !ok {
			t.Errorf("key: %s has unexpected type: %T\n", str, v)
		}
		if s != str {
			t.Errorf("key: %s has unexpected value: %s\n", str, s)
		}
		a, ok := c.All()
		if !ok {
			t.Errorf("get all key/values failed\n")
		}
		v, ok = a[str]
		if !ok {
			t.Errorf("key: %s not found\n", str)
		}
		s, ok = v.(string)
		if !ok {
			t.Errorf("key: %s has unexpected type: %T\n", str, v)
		}
		if s != str {
			t.Errorf("key: %s has unexpected value: %s\n", str, s)
		}
		c.Del(str)
		_, ok = c.Get(str)
		if ok {
			t.Errorf("key: %s should be deleted\n", str)
		}
		c.Put(a)
		v, ok = c.Get(str)
		if !ok {
			t.Errorf("key: %s not found\n", str)
		}
		s, ok = v.(string)
		if !ok {
			t.Errorf("key: %s has unexpected type: %T\n", str, v)
		}
		if s != str {
			t.Errorf("key: %s has unexpected value: %s\n", str, s)
		}
		c.Del(str)
		_, ok = c.Get(str)
		if ok {
			t.Errorf("key: %s should be deleted\n", str)
		}
	}
	c.Clr()
}

func BenchmarkSingle_Set(b *testing.B) {
	counter.Set(counterStart)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			str := counter.Next().String()
			singleCache.Set(str, str)
		}
	})
}

func BenchmarkSharding_Set(b *testing.B) {
	counter.Set(counterStart)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			str := counter.Next().String()
			shardingCache.Set(str, str)
		}
	})
}

func BenchmarkSingle_Get(b *testing.B) {
	counter.Set(counterStart)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			str := counter.Next().String()
			singleCache.Get(str)
		}
	})
}

func BenchmarkSharding_Get(b *testing.B) {
	counter.Set(counterStart)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			str := counter.Next().String()
			shardingCache.Get(str)
		}
	})
}

func BenchmarkSingle_Del(b *testing.B) {
	counter.Set(counterStart)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			str := counter.Next().String()
			singleCache.Del(str)
		}
	})
}

func BenchmarkSharding_Del(b *testing.B) {
	counter.Set(counterStart)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			str := counter.Next().String()
			shardingCache.Del(str)
		}
	})
}

func BenchmarkSingle_Clr(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			singleCache.Clr()
		}
	})
}

func BenchmarkSharding_Clr(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			shardingCache.Clr()
		}
	})
}

type atomicInteger int64

func (a *atomicInteger) Set(val int64) {
	atomic.StoreInt64((*int64)(a), val)
}

func (a *atomicInteger) Next() atomicInteger {
	return atomicInteger(atomic.AddInt64((*int64)(a), 1))
}

func (a atomicInteger) String() string {
	return strconv.FormatInt(int64(a), 10)
}
