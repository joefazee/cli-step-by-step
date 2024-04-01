package main

import (
	"bytes"
	"testing"
)

func Test_greeting(t *testing.T) {

	buff := new(bytes.Buffer)
	name := "john"
	greeting(buff, name)
	expected := "hello, john\n"

	got := buff.String()
	if expected != got {
		t.Errorf("expected %v; got %v", expected, got)
	}
}
