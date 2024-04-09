package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
)

type customWriter struct {
	buf bytes.Buffer
	w   io.Writer
}

func newCustomWriter(w io.Writer) *customWriter {
	return &customWriter{w: w}
}

func (w *customWriter) Write(p []byte) (n int, err error) {
	w.buf.Write(p)
	return w.w.Write(p)
}

func (w *customWriter) Bytes() []byte {
	return w.buf.Bytes()
}

func main() {

	command := exec.Command("ping", "google.com")

	var errStdout, errStderr error
	stdoutIn, _ := command.StdoutPipe()
	stderrIn, _ := command.StderrPipe()
	stdout := newCustomWriter(os.Stdout)
	stderr := newCustomWriter(os.Stderr)
	err := command.Start()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
		wg.Done()
	}()

	_, errStderr = io.Copy(stderr, stderrIn)
	wg.Wait()

	err = command.Wait()
	if err != nil {
		log.Fatal(err)
	}

	if errStdout != nil || errStderr != nil {
		log.Fatalf("errStdout: %s, errStderr: %s", errStdout, errStderr)
	}

	fmt.Printf("out:\n%s\nerr:\n%s\n", string(stdout.Bytes()), string(stderr.Bytes()))

}
