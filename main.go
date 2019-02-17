package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/pmezard/go-difflib/difflib"
)

func main() {
	// Pipe Stdin through to Stdout, but tap into the stream for our parser:
	in := io.TeeReader(os.Stdin, os.Stdout)
	// TODO: buffer Stdin so that we can put our fancy formatted diff output
	// right after the real test output.

	s := bufio.NewScanner(in)
	s.Split(ScanTestOutputs)

	for s.Scan() {
		eg, ok := parseFailingExample(s.Text())
		if ok {
			fmt.Printf("\033[1;91m%v\033[0m\n", eg.Diff())
		}
	}
	err := s.Err()
	if err != nil {
		log.Fatal(err)
	}
}

var (
	runHeaderRE         = regexp.MustCompile(`(?m)^=== RUN\b.*$`)
	failExampleHeaderRE = regexp.MustCompile(`^--- FAIL: (Example\S*)\s+\(.+\)`)
)

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
		// Return everything after the header, but throw out the last two lines
		// of output because they show the status:
		advance, token := trimLastNLines(data[start:], 2)
		return advance, token, nil
	}
	// Ask for more data:
	return 0, nil, nil
}

func trimLastNLines(data []byte, n int) (advance int, token []byte) {
	advance = len(data)
	for i := 0; advance > 0 && i < n; i++ {
		lastIndex := bytes.LastIndex(data[0:advance-1], []byte("\n"))
		if lastIndex < 0 {
			advance = 0
			break
		}
		advance = lastIndex + 1
	}
	return advance, data[:advance]
}

// Example is the parsed the output of a failing example.
type Example struct {
	Name string
	Got  string
	Want string
}

func parseFailingExample(s string) (Example, bool) {
	matches := failExampleHeaderRE.FindAllStringSubmatch(s, 1)
	if len(matches) != 1 {
		return Example{}, false
	}
	eg := Example{
		Name: matches[0][1],
	}

	splits := strings.SplitAfterN(s, "\ngot:\n", 2)
	if len(splits) != 2 {
		log.Printf("got: %#v", splits)
		return Example{}, false
	}
	s = splits[1]

	splits = strings.SplitN(s, "\nwant:\n", 2)
	if len(splits) != 2 {
		log.Printf("want: %#v", splits)
		return Example{}, false
	}
	eg.Got = splits[0]
	eg.Want = strings.TrimSuffix(splits[1], "\n")

	return eg, true
}

// Diff formats the difference between Want and Got.
func (eg Example) Diff() string {
	diff, _ := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:        difflib.SplitLines(eg.Want),
		B:        difflib.SplitLines(eg.Got),
		FromFile: "Want",
		ToFile:   "Got",
		Context:  2,
	})
	return diff
}
