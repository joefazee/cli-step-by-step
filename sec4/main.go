package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	wordCount := 0

	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		wordCount += len(words)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading input: %s", err.Error())
		os.Exit(1)
	}

	fmt.Printf("Total words: %d\n", wordCount)

}
