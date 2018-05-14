# gopowerschool
powerschool in go


usage:

```go
package main

import (
        "github.com/reteps/gopowerschool"
        "fmt"
)

func main() {
        client := gopowerschool.Client("https://example.com")
        student, err := client.GetStudent("username", "password")
        if err != nil {
                panic(err)
        }   
        fmt.Println(student)
}
```
