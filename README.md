# gls #

A **Goroutine Local Storage** implementation base on the runtime unique ID of each goroutine.

**Example:**

```go
package main

import (
	"fmt"
	"github.com/fengyoulin/gls"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			gls.Set("index", idx)
			something()
			gls.Clr()
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func something() {
	if v, ok := gls.Get("index"); ok {
		if idx, ok := v.(int); ok {
			fmt.Printf("index: %d\n", idx)
		}
	}
}
```
The example before uses the default gls cache, you can use the `gls.New()` function to allocate a new gls cache for your purpose.

**Attention:**

You should always try to use `context.Context` as google suggested. If you prefer to use this `gls` module, remember to call the `Clr()` at end of your goroutine, it will protect you from a memory leak.
