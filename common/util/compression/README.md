## Compressions

```go
package main

import (
    "fmt"

	"github.com/golang-common-packages/template/common/util/compression"
)

func main() {
	fmt.Println(compression.Compress([]byte("Backend Golang")))
}
```