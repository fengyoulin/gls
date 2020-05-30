# gls #

A **Goroutine Local Storage** implementation base on the runtime unique ID of each goroutine.

```shell script
$ go get github.com/fengyoulin/gls
```

## Usage

```go
import "github.com/fengyoulin/gls"
```

#### func  All

```go
func All() (kvs map[string]interface{}, ok bool)
```
All returns all the key/values in current goroutine's local storage

#### func  Clr

```go
func Clr()
```
Clr clears the current goroutine's local storage

#### func  Del

```go
func Del(key string)
```
Del deletes a key in current goroutine's local storage

#### func  Get

```go
func Get(key string) (val interface{}, ok bool)
```
Get get the value of a key in current goroutine's local storage

#### func  Put

```go
func Put(kvs map[string]interface{})
```
Put puts all the key/values into current goroutine's local storage

#### func  Set

```go
func Set(key string, val interface{})
```
Set set the value to a key in current goroutine's local storage

#### func  Go

```go
func Go(fn func(), c Cache)
```
Go start a new goroutine, inherit current goroutine's local storage

#### func  GoWith

```go
func GoWith(fn func(), c Cache, ls map[string]interface{})
```
GoWith start a new goroutine, with the given local storage

#### type Cache

```go
type Cache interface {
	// All returns all the key/values in current goroutine's local storage
	All() (kvs map[string]interface{}, ok bool)
	// Put puts all the key/values into current goroutine's local storage
	Put(kvs map[string]interface{})
	// Clr clears the current goroutine's local storage
	Clr()
	// Del deletes a key in current goroutine's local storage
	Del(key string)
	// Get get the value of a key in current goroutine's local storage
	Get(key string) (val interface{}, ok bool)
	// Set set the value to a key in current goroutine's local storage
	Set(key string, val interface{})
}
```

Cache is the gls cache interface

#### func  New

```go
func New(shardingMode bool) Cache
```
New creates a goroutine local storage cache

## Example

```go
package main

import (
	"fmt"
	"github.com/fengyoulin/gls"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	gls.Set("gls", "v0.4.0")

	for i := 0; i < 10; i++ {
		wg.Add(1)
		i := i
		gls.Go(func() {
			defer wg.Done()
			if ls, ok := gls.All(); ok {
				fmt.Printf("%d, %v\n", i, ls)
			}
		}, nil)
	}
	wg.Wait()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		gls.GoWith(func() {
			defer wg.Done()
			if ls, ok := gls.All(); ok {
				fmt.Printf("%v\n", ls)
			}
		}, nil, map[string]interface{}{"i": i})
	}
	wg.Wait()
}
```
The example before uses the default gls cache, you can use the `gls.New()` function to allocate a new gls cache for your purpose.

## Attention

You should always try to use `context.Context` as google suggested. If you prefer to use this `gls` module, remember to call the `Clr()` at end of your goroutine, it will protect you from a memory leak.
