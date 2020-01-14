## Encodings

```go
package main

import (
    "fmt"

	"github.com/golang-microservices/template/common/util/encoding"
)

func main() {
	fmt.Println(xencodings.Base32Encode([]byte("Backend Golang")))
	fmt.Println(xencodings.Base64Encode([]byte("Backend Golang")))
}
```