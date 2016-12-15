# csv
parse csv to gostruct

```
package main

import (
        "fmt"

        "github.com/radioinmyhead/csv"
)

type Base struct {
        Id int
}

type Animal struct {
        Base  `csv:"extends"`
        Name  string
        Order string
}

var data = `id,name,order
#uuid,#name,#xxxx
12,Platypus,Monotremata
34,Quoll,Dasyuromorphia`

func main() {
        var ret []Animal // or []*Animal
        err := csv.ReadString(data, &ret)
        fmt.Println(ret, err)
}
```
