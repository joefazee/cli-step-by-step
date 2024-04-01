package main

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func Test_greeting(t *testing.T) {

	originalStdout := os.Stdout

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	os.Stdout = w

	name := "john"
	greeting(name)
	_ = w.Close()
	os.Stdout = originalStdout
	expected := "hello, john\n"

	bs, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}

	got := string(bs)
	if expected != got {
		t.Errorf("expected %v; got %v", expected, got)
	}

	fmt.Println("done")
}
