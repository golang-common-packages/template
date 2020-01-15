## Hash

```go
package main

import (
    "fmt"

	"github.com/golang-common-packages/template/common/util/hash"
)

func main() {
	fmt.Println(hash.FNV32("Backend Golang")) 
	fmt.Println(hash.FNV32a("Backend Golang"))
	fmt.Println(hash.FNV64("Backend Golang"))
	fmt.Println(hash.FNV64a("Backend Golang"))
	fmt.Println(hash.MD5("Backend Golang"))
	fmt.Println(hash.SHA1("Backend Golang"))
	fmt.Println(hash.SHA256("Backend Golang"))
	fmt.Println(hash.SHA512("Backend Golang"))
}
```