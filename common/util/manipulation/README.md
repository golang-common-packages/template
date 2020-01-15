## Manipulation

```go
package main

import (

"fmt"
"math/rand"
"time"
    
	"github.com/golang-common-packages/template/common/util/manipulation"
)

func main() {
	source := rand.NewSource(time.Now().UnixNano())

	array := []interface{}{"a", "b", "c"}
	manipulation.Shuffle(array, source)

	fmt.Println(array) // [c b a]
}
```