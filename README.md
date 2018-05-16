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

using other methods in xml_parser.go (example: get picture):
```go
client := gopowerschool.Client("https://example.com")
session, userID, err := client.CreateUserSessionAndStudent("username", "password")
if err != nil {
        panic(err)
}
arguments := gopowerschool.GetStudentPhoto{UserSessionVO: session, StudentID: userID}
response, err := client.GetStudentPhoto(&arguments)
if err != nil {
        panic(err)
}
fmt.Println(string(response))
```
