package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

const EXIT = "exit"

func start(prompt string, in io.Reader, out io.Writer) {

	scanner := bufio.NewScanner(in)
	for {
		fmt.Print(prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		if line == EXIT || line == "" {
			fmt.Fprint(out, "program closed")
			os.Exit(0)
		}

		parts := strings.Split(line, " ")
		if len(parts) == 0 {
			fmt.Fprint(out, "no input")
			os.Exit(1)
		}

		command := parts[0]
		args := parts[1:]

		runner := exec.Command(command, args...)
		result, err := runner.CombinedOutput()

		if err != nil {
			fmt.Fprint(out, err.Error())
			os.Exit(1)
		}

		out.Write(result)

	}

}

func main() {
	start(">>> ", os.Stdin, os.Stdout)
}
