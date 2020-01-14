## Compressions

```go
package main

import (
    "fmt"

	"github.com/golang-microservices/template/common/util/compression"
)

func main() {
	fmt.Println(compression.Compress([]byte("Backend Golang")))
}
```