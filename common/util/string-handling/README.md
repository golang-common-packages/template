## String handling

```go
package main

import (
    "fmt"
    
	"github.com/golang-common-packages/template/common/util/string-handling"
)

func main() {
	fmt.Println(strhandling.IsEmpty(""))     // true
	fmt.Println(strhandling.IsEmpty("text")) // false
	fmt.Println(strhandling.IsEmpty("	"))  // false

	fmt.Println(strhandling.IsNotEmpty(""))     // false
	fmt.Println(strhandling.IsNotEmpty("text")) // true
	fmt.Println(strhandling.IsNotEmpty("	")) // true

	fmt.Println(strhandling.IsBlank(""))     // true
	fmt.Println(strhandling.IsBlank("	"))  // true
	fmt.Println(strhandling.IsBlank("text")) // false

	fmt.Println(strhandling.IsNotBlank(""))     // false
	fmt.Println(strhandling.IsNotBlank("	")) // false
	fmt.Println(strhandling.IsNotBlank("text")) // true

	fmt.Println(strhandling.Left("", 5))            // "     "
	fmt.Println(strhandling.Left("X", 5))           // "X    "
	fmt.Println(strhandling.Left("😎⚽", 4))        // "😎⚽  "
	fmt.Println(strhandling.Left("ab\u0301cde", 8)) // "ab́cde   "

	fmt.Println(strhandling.Right("", 5))            // "     "
	fmt.Println(strhandling.Right("X", 5))           // "    X"
	fmt.Println(strhandling.Right("😎⚽", 4))        // "  😎⚽"
	fmt.Println(strhandling.Right("ab\u0301cde", 8)) // "   ab́cde"

	fmt.Println(strhandling.Center("", 5))            // "     "
	fmt.Println(strhandling.Center("X", 5))           // "  X  "
	fmt.Println(strhandling.Center("😎⚽", 4))        // " 😎⚽ "
	fmt.Println(strhandling.Center("ab\u0301cde", 8)) // "  ab́cde "

	fmt.Println(strhandling.Length(""))                                          // 0
	fmt.Println(strhandling.Length("X"))                                         // 1
	fmt.Println(strhandling.Length("b\u0301"))                                   // 1
	fmt.Println(strhandling.Length("😎⚽"))                                      // 2
	fmt.Println(strhandling.Length("Les Mise\u0301rables"))                      // 14
	fmt.Println(strhandling.Length("ab\u0301cde"))                               // 5
	fmt.Println(strhandling.Length("This `\xc5` is an invalid UTF8 character"))  // 37
	fmt.Println(strhandling.Length("The quick bròwn 狐 jumped over the lazy 犬")) // 40

	fmt.Println(strhandling.Reverse(""))                                            // ""
	fmt.Println(strhandling.Reverse("X"))                                           // "X"
	fmt.Println(strhandling.Reverse("😎⚽"))                                        // "⚽😎"
	fmt.Println(strhandling.Reverse("Les Mise\u0301rables"))                        // "selbare\u0301siM seL"
	fmt.Println(strhandling.Reverse("This `\xc5` is an invalid UTF8 character"))    // "retcarahc 8FTU dilavni na si `�` sihT"
	fmt.Println(strhandling.Reverse("The quick bròwn 狐 jumped over the lazy 犬"))  // "犬 yzal eht revo depmuj 狐 nwòrb kciuq ehT"

	c := [100]byte{'a', 'b', 'c'}
    fmt.Println("C: ", len(c), c[:4])
    g := CToGoString(c[:])			// C:  100 [97 98 99 0]
    fmt.Println("Go:", len(g), g)	// Go: 3 abc
}
```