package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)
	s := NewScanner(os.Stdin)

	for {
		eg, err := s.Next()
		if err == nil {
			fmt.Println(eg)
			continue
		}
		if err == EOF {
			return
		}
		log.Fatal(err)
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

// Next returns the next
func (s *Scanner) Next() (Example, error) {
	var got bytes.Buffer
	var want bytes.Buffer

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
			fmt.Fprintln(&got, line)
		case insideWant:
			// NOTE: This requires `go test -v`, so that it's printed out before each test:
			if runRE.MatchString(line) || line == "FAIL" {
				s.state = outside
				eg.Got = got.String()
				eg.Want = want.String()
				return eg, nil
			}
			fmt.Fprintln(&want, line)
		}
	}

	err := s.scanner.Err()
	if err != nil {
		return eg, err
	}

	if got.Len() > 0 || want.Len() > 0 {
		eg.Got = got.String()
		eg.Want = want.String()
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
