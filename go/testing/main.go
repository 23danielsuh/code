package main

import "fmt"

func main() {
    test := func(name, t string) {
        fmt.Println(name, t)
    }
    test("hi", "hello")
}
