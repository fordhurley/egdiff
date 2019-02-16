package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
)

func main() {
	s := NewScanner(os.Stdin)

	for {
		eg, err := s.NextExample()
		if err == EOF {
			return
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(eg)
	}
}

type Example struct {
	Name string
	Got  string
	Want string
}

type Scanner struct {
	scanner *bufio.Scanner
	state   state
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		scanner: bufio.NewScanner(r),
	}
}

type state int

const (
	outside state = iota
	inside
	insideGot
	insideWant
)

var runRE = regexp.MustCompile(`^=== RUN\b`)
var failRE = regexp.MustCompile(`^--- FAIL: (Example\S*)\b.*\)$`)
var gotRE = regexp.MustCompile(`^got:$`)
var wantRE = regexp.MustCompile(`^want:$`)

// NextExample returns the output from the next example.
func (s *Scanner) NextExample() (Example, error) {
	s.state = outside

	eg := Example{}

	for s.scanner.Scan() {
		line := s.scanner.Text()

		switch s.state {
		case outside:
			matches := failRE.FindStringSubmatch(line)
			if matches != nil {
				eg.Name = matches[1]
				s.state = inside
			}
		case inside:
			if gotRE.MatchString(line) {
				s.state = insideGot
			} else {
				// The very next line _has_ to be "got:", or something must be wrong...
				return eg, fmt.Errorf("line should be got!!!!\n%s", line)
			}
		case insideGot:
			if wantRE.MatchString(line) {
				s.state = insideWant
				continue
			}
			eg.Got += line + "\n"
		case insideWant:
			// NOTE: This requires `go test -v`, so that it's printed out before each test:
			if runRE.MatchString(line) || line == "FAIL" {
				return eg, nil
			}
			eg.Want += line + "\n"
		}
	}

	err := s.scanner.Err()
	if err != nil {
		return eg, err
	}

	if len(eg.Got) > 0 || len(eg.Want) > 0 {
		return eg, nil
	}

	return eg, EOF
}

var EOF = errors.New("EOF")

func (e Example) String() string {
	return fmt.Sprintf(`Example{
	Name: %q,
	Got: %q,
	Want: %q,
}`, e.Name, e.Got, e.Want)
}
