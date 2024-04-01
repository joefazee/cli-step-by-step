package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("pass all the args")
		os.Exit(1)
	}

	fmt.Println(os.Args[1])
	fmt.Println(os.Args[2])

}
