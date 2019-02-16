package main

import (
	"fmt"
	"regexp"
)

var lineEndingRE = regexp.MustCompile(`\r?\n`)

// replaceLineEndings changes \r\n to \n. It's idempotent.
func replaceLineEndings(raw string) string {
	return lineEndingRE.ReplaceAllString(raw, "\n")
}
func Example_replaceLineEndings() {
	fmt.Printf("%#v\n", replaceLineEndings("a\r\n\r\nb\r\n\r\nc"))
	fmt.Printf("%#v\n", replaceLineEndings("a\r\nb\r\nc"))
	fmt.Printf("%#v\n", replaceLineEndings("a\nb\nc"))
	fmt.Printf("%#v\n", replaceLineEndings("abc"))

	// Output:
	// "a\n\nb\n\nc"
	// "a\nb/nc"
	// "a\nb\nc"
	// "abc"
}
