## Concurrency

```go
package main

import (
    "time"
    "fmt"
    
	"github.com/golang-microservices/template/common/util/concurrency"
)

func main() {
    func1 := func() {
            for char := 'a'; char < 'a' + 3; char++ {
                fmt.Printf("%c ", char)
            }
    }
    
    func2 := func() {
            for number := 1; number < 4; number++ {
                fmt.Printf("%d ", number)
            }
    }
    
    concurrency.Parallelize(func1, func2)  // a 1 b 2 c 3
    
    concurrency.ParallelizeTimeout(time.Minute, func1, func2)  // a 1 b 2 c 3
}
```