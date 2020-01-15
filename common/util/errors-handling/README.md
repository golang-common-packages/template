## Errors handling

```go
package main

import (
    "fmt"

	"github.com/golang-common-packages/template/common/util/errhandling"
)

func main() {
	fmt.Println(errhandling.DefaultErrorIfNil(nil, "Cool"))                // "Cool"
	fmt.Println(errhandling.DefaultErrorIfNil(errors.New("Oops"), "Cool")) // "Oops"
}
```