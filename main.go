package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

func main() {
	// Pipe Stdin through to Stdout, but tap into the stream for our parser:
	in := io.TeeReader(os.Stdin, os.Stdout)
	// TODO: buffer Stdin so that we can put our fancy formatted diff output
	// right after the real test output.

	s := bufio.NewScanner(in)
	s.Split(ScanTestOutputs)

	for s.Scan() {
		fmt.Printf("\033[1;91m%v\033[0m\n", s.Text())
	}
	err := s.Err()
	if err != nil {
		log.Fatal(err)
	}
}

var runHeaderRE = regexp.MustCompile(`(?m)^=== RUN\b.*$`)

// ScanTestOutputs is a split function for a bufio.Scanner that returns the
// output related to a single test
func ScanTestOutputs(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	matches := runHeaderRE.FindAllIndex(data, 2)
	if len(matches) == 0 {
		// Ask for more data:
		return 0, nil, nil
	}
	headerLoc := matches[0]
	start := headerLoc[1] + 1 // one after the end of the header
	if len(matches) == 2 {
		// Found a header before and a header after, so return everything in
		// between:
		end := matches[1][0] // start of the next header (we'll capture the trailing newline)
		between := data[start:end]
		return end, between, nil
	}
	if len(matches) == 1 && atEOF {
		// Return everything after the header:
		return len(data), data[start : len(data)-1], nil
		// FIXME: Skip the last two lines of output, they show the status
	}
	return 0, nil, nil
}

// Example is the parsed the output of a failing example.
type Example struct {
	Name string
	Got  string
	Want string
}

func (e Example) String() string {
	return fmt.Sprintf(`Example{
	Name: %q,
	Got: %q,
	Want: %q,
}`, e.Name, e.Got, e.Want)
}
