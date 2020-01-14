## Conditions

```go
package main

import (
    "fmt"

	"github.com/golang-microservices/template/common/util/condition"
)

func main() {
	fmt.Println(condition.IfThen(1 == 1, "Yes")) // "Yes"
	fmt.Println(condition.IfThen(1 != 1, "Woo")) // nil
	fmt.Println(condition.IfThen(1 < 2, "Less")) // "Less"

	fmt.Println(condition.IfThenElse(1 == 1, "Yes", false)) // "Yes"
	fmt.Println(condition.IfThenElse(1 != 1, nil, 1))       // 1
	fmt.Println(condition.IfThenElse(1 < 2, nil, "No"))     // nil

	fmt.Println(condition.DefaultIfNil(nil, nil))  // nil
	fmt.Println(condition.DefaultIfNil(nil, ""))   // ""
	fmt.Println(condition.DefaultIfNil("A", "B"))  // "A"
	fmt.Println(condition.DefaultIfNil(true, "B")) // true
	fmt.Println(condition.DefaultIfNil(1, false))  // 1

	fmt.Println(condition.FirstNonNil(nil, nil))                // nil
	fmt.Println(condition.FirstNonNil(nil, ""))                 // ""
	fmt.Println(condition.FirstNonNil("A", "B"))                // "A"
	fmt.Println(condition.FirstNonNil(true, "B"))               // true
	fmt.Println(condition.FirstNonNil(1, false))                // 1
	fmt.Println(condition.FirstNonNil(nil, nil, nil, 10))       // 10
	fmt.Println(condition.FirstNonNil(nil, nil, nil, nil, nil)) // nil
	fmt.Println(condition.FirstNonNil())                        // nil
}
```