## Manipulation

```go
package main

import (
    "math/rand"
	"time"
    "fmt"
    
	"github.com/golang-microservices/template/common/util/manipulation"
)

func main() {
	source := rand.NewSource(time.Now().UnixNano())

	array := []interface{}{"a", "b", "c"}
	manipulation.Shuffle(array, source)

	fmt.Println(array) // [c b a]
}
```