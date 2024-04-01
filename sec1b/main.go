package main

import (
	"fmt"
)

func main() {
	greeting("John")
}

func greeting(name string) {
	fmt.Println("hello,", name)
}
