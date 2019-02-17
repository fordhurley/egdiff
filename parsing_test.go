package main

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
)

const verboseTestOutput = `=== RUN   Test_sayHi
--- PASS: Test_sayHi (0.00s)
=== RUN   Example_sayHi
--- PASS: Example_sayHi (0.00s)
=== RUN   Example_replaceLineEndings
--- FAIL: Example_replaceLineEndings (0.00s)
got:
"a\n\nb\n\nc"
"a\nb\nc"
"a\nb\nc"
"abc"
want:
"a\n\nb\n\nc"
"a\nb/nc"
"a\nb\nc"
"abc"
=== RUN   Example_meta
--- FAIL: Example_meta (0.00s)
got:
this is the output
with tricky stuff mixed in
got:
(tricked you?)
want:
oh no
want:
this is the output
with tricky stuf mixed in
got:
(tricked you?)
want:
oh no
FAIL
FAIL	_/Users/ford/src/github.com/fordhurley/egdiff/testdata	0.005s
`

func Example_runHeaderRE() {
	matches := runHeaderRE.FindAllString(verboseTestOutput, -1)
	if matches == nil {
		panic("no matches")
	}
	fmt.Printf("%d matches:\n", len(matches))
	for _, m := range matches {
		fmt.Println(m)
	}

	// Output:
	// 4 matches:
	// === RUN   Test_sayHi
	// === RUN   Example_sayHi
	// === RUN   Example_replaceLineEndings
	// === RUN   Example_meta
}

func TestScanTestOutputs(t *testing.T) {
	var outputs []string

	s := bufio.NewScanner(strings.NewReader(verboseTestOutput))
	s.Split(ScanTestOutputs)
	for s.Scan() {
		outputs = append(outputs, s.Text())
	}
	err := s.Err()
	if err != nil {
		t.Fatal(err)
	}

	expectedOutputs := []string{
		"--- PASS: Test_sayHi (0.00s)\n",
		"--- PASS: Example_sayHi (0.00s)\n",
		`--- FAIL: Example_replaceLineEndings (0.00s)
got:
"a\n\nb\n\nc"
"a\nb\nc"
"a\nb\nc"
"abc"
want:
"a\n\nb\n\nc"
"a\nb/nc"
"a\nb\nc"
"abc"
`,
		`--- FAIL: Example_meta (0.00s)
got:
this is the output
with tricky stuff mixed in
got:
(tricked you?)
want:
oh no
want:
this is the output
with tricky stuf mixed in
got:
(tricked you?)
want:
oh no
`,
	}

	if len(outputs) != len(expectedOutputs) {
		t.Errorf("expected %d outputs, got %d", len(expectedOutputs), len(outputs))
	}

	for i, output := range outputs {
		expectedOutput := expectedOutputs[i]
		if output != expectedOutput {
			t.Errorf("expected:\n%q\ngot:\n%q", expectedOutput, output)
		}
	}
}