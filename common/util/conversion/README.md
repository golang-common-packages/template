## Conversion

Return a pretty JSON representation of any interface

```go
package main

import (
    "fmt"

	"github.com/golang-common-packages/template/common/util/conversion"
)

func main() {
    x := map[string]interface{}{"number": 1, "string": "cool", "bool": true, "float": 1.5}    
    fmt.Println(conversion.PrettyJson(x))
}
```

```json
{
	"bool": true,
	"float": 1.5,
	"number": 1,
	"string": "cool"
}
```

Convert any interface to a String

```go
package main

import (
    "fmt"

	"github.com/golang-common-packages/template/common/util/conversion"
)

func main() {
    x := map[string]interface{}{"number": 1, "string": "cool", "bool": true, "float": 1.5}    
    fmt.Println(conversion.Stringify(x))
}
```

```
{"bool":true,"float":1.5,"number":1,"string":"cool"} <nil>
```

Convert any string back to its original struct

```go
package main

import (
    "fmt"
    
	"github.com/golang-common-packages/template/common/util/conversion"
)

func main() {
	x := "{\"bool\":true,\"float\":1.5,\"number\":1,\"string\":\"cool\"}"
	var results map[string]interface{}
    fmt.Println(conversion.Structify(x, &results))
    fmt.Println(results)
}
```

```
<nil>
map[bool:true float:1.5 number:1 string:cool]
```
