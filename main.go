package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/pmezard/go-difflib/difflib"
)

func main() {
	s := bufio.NewScanner(os.Stdin)
	s.Split(ScanTestOutputs)

	lastChunk := ""
	foundExample := false

	for s.Scan() {
		output := s.Text()
		eg, ok := parseFailingExample(output)
		fmt.Print(output)
		if ok {
			foundExample = true
			fmt.Printf("\033[1;91m%v\033[0m", eg.Diff())
		}
		lastChunk = output
	}
	err := s.Err()
	if err != nil {
		log.Fatal(err)
	}

	if !foundExample {
		fmt.Fprintln(os.Stderr, "egdiff: no failing example found, did you pass `-v` to go test?")
	}

	lastNewline := NthFromLastIndex([]byte(lastChunk), []byte("\n"), 1)
	if lastChunk[lastNewline+1:lastNewline+5] == "FAIL" {
		os.Exit(1)
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
		if atEOF {
			// Return whatever we got:
			return len(data), data, nil
		}
		// Ask for more data:
		return 0, nil, nil
	}

	headerLoc := matches[0]
	start := headerLoc[0] // including the header

	if len(matches) == 2 {
		// Found a header before and a header after, so return everything up to
		// the next header:
		nextHeaderLoc := matches[1]
		end := nextHeaderLoc[0] // start of the next header (we'll capture the trailing newline)
		return end, data[start:end], nil
	}

	if len(matches) == 1 && atEOF {
		// Return everything after the header, but throw out the last two lines
		// of output because they show the status:
		index := NthFromLastIndex(data[start:], []byte("\n"), 2)
		// Include that last newline:
		index++
		return index, data[start:index], nil
	}

	// Ask for more data:
	return 0, nil, nil
}

// NthFromLastIndex is like bytes.LastIndex, but can skip back N from the end.
// Counting start from 0, so:
//
//    NthFromLastIndex(data, sep, 0) == bytes.LastIndex(data, sep)
//
func NthFromLastIndex(data []byte, sep []byte, n int) int {
	index := bytes.LastIndex(data, sep)
	if index < 0 {
		return index
	}

	for i := 0; index > 0 && i < n; i++ {
		lastIndex := bytes.LastIndex(data[0:index], sep)
		if lastIndex < 0 {
			index = 0
			break
		}
		index = lastIndex
	}
	return index
}

// Example is the parsed output of a failing example.
type Example struct {
	Name string
	Got  string
	Want string
}

func parseFailingExample(s string) (Example, bool) {
	// Trim RUN header:
	splits := strings.SplitAfterN(s, "\n", 2)
	if len(splits) != 2 {
		return Example{}, false
	}
	s = splits[1]

	// Extract name from FAIL header:
	matches := failExampleHeaderRE.FindAllStringSubmatch(s, 1)
	if len(matches) != 1 {
		return Example{}, false
	}
	eg := Example{
		Name: matches[0][1],
	}

	// Grab everything after the "got:" line:
	splits = strings.SplitAfterN(s, "\ngot:\n", 2)
	if len(splits) != 2 {
		return Example{}, false
	}
	s = splits[1]

	// Split on the "want:" line:
	splits = strings.SplitN(s, "\nwant:\n", 2)
	if len(splits) != 2 {
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
