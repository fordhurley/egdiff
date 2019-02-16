package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func sayHi(name string) {
	fmt.Printf("hi, %s!\n", name)
}

func Example_sayHi() {
	sayHi("winston")
	sayHi("sadie")

	// Output:
	// hi, winston!
	// hi, sadie!
}

func Test_sayHi(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	stdout := os.Stdout
	os.Stdout = w
	defer func() {
		os.Stdout = stdout
	}()

	sayHi("boy")
	w.Close()

	out, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}

	expected := "hi, boy!\n"
	if string(out) != expected {
		t.Errorf("expected: %q, got %q", expected, string(out))
	}
}
