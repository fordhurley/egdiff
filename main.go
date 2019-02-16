package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
)

var runRE = regexp.MustCompile(`^=== RUN\b`)
var failRE = regexp.MustCompile(`^--- FAIL: Example.+\)$`)
var gotRE = regexp.MustCompile(`^got:$`)
var wantRE = regexp.MustCompile(`^want:$`)

type state int

const (
	outside state = iota
	inside
	insideGot
	insideWant
)

func main() {
	state := outside

	var got bytes.Buffer
	var want bytes.Buffer

	reset := func() {
		if got.Len() > 0 || want.Len() > 0 {
			fmt.Println("CAPTURED GOT INPUT:")
			fmt.Println(got.String())
			fmt.Println()
			fmt.Println("CAPTURED WANT INPUT:")
			fmt.Println(want.String())
			fmt.Println()
		}
		got.Reset()
		want.Reset()
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		switch state {
		case outside:
			if failRE.MatchString(line) {
				reset()
				state = inside
			}
		case inside:
			if gotRE.MatchString(line) {
				state = insideGot
			} else {
				// The very next line _has_ to be "got:", or something must be wrong...
				log.Fatalf("line should be got!!!!\n%s", line)
			}
		case insideGot:
			if wantRE.MatchString(line) {
				state = insideWant
				continue
			}
			fmt.Fprintln(&got, line)
		case insideWant:
			// NOTE: This requires `go test -v`, so that it's printed out before each test:
			if runRE.MatchString(line) {
				state = outside
				continue
			}
			fmt.Fprintln(&want, line)
		}
	}
	err := scanner.Err()
	if err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(1)
	}

	reset()
}
