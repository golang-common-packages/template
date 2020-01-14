## Divide

* Given two integers representing the numerator and denominator of a fraction, return the fraction in string format with the repeating part enclosed in parentheses.

```go
package main

import (
    "fmt"

	"github.com/golang-microservices/template/common/util/divide"
)

func main() {
    fmt.Println(divide.Divide(0, 0))     // "ERROR"
    fmt.Println(divide.Divide(1, 2))     // "0.5(0)"
    fmt.Println(divide.Divide(0, 3))     // "0.(0)"
    fmt.Println(divide.Divide(10, 3))    // "3.(3)"
    fmt.Println(divide.Divide(22, 7))    // "3.(142857)"
    fmt.Println(divide.Divide(100, 145)) // "0.(6896551724137931034482758620)"
}
```