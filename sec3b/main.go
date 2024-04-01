package main

import (
	"flag"
	"fmt"
)

func main() {

	username := flag.String("username", "", "username")
	password := flag.String("password", "", "password")

	var port int
	flag.IntVar(&port, "port", 0, "db port")

	flag.Parse()

	fmt.Printf("username: %s\npassword: %s\nport: %d\n",
		*username, *password, port)
}
