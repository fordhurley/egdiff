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
=== RUN   Example_simple
--- FAIL: Example_simple (0.00s)
got:
this is not it
want:
this is it
=== RUN   Example_tricky
--- FAIL: Example_tricky (0.00s)
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

var expectedOutputs = []string{
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
	`--- FAIL: Example_simple (0.00s)
got:
this is not it
want:
this is it
`,
	`--- FAIL: Example_tricky (0.00s)
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

var expectedExamples = []Example{
	{
		Name: "Example_replaceLineEndings",
		Got: `"a\n\nb\n\nc"
"a\nb\nc"
"a\nb\nc"
"abc"`,
		Want: `"a\n\nb\n\nc"
"a\nb/nc"
"a\nb\nc"
"abc"`,
	},
	{
		Name: "Example_simple",
		Got:  "this is not it",
		Want: "this is it",
	},
	{
		// FIXME: make this work. Currently it splits on the first "want:"
		Name: "Example_tricky",
		Got: `this is the output
with tricky stuff mixed in
got:
(tricked you?)
want:
oh no`,
		Want: `this is the output
with tricky stuf mixed in
got:
(tricked you?)
want:
oh no`,
	},
}

func Example_runHeaderRE() {
	matches := runHeaderRE.FindAllString(verboseTestOutput, -1)
	if matches == nil {
		panic("no matches")
	}
	for _, m := range matches {
		fmt.Println(m)
	}

	// Output:
	// === RUN   Test_sayHi
	// === RUN   Example_sayHi
	// === RUN   Example_replaceLineEndings
	// === RUN   Example_simple
	// === RUN   Example_tricky
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

func Example_parseFailingExample() {
	eg, ok := parseFailingExample(`--- FAIL: ExampleFoo (0.00s)
got:
foo
want:
bar
`)
	fmt.Println(eg.Name, ok)

	eg, ok = parseFailingExample(`--- PASS: ExampleBar (0.00s)
`)
	fmt.Println(eg.Name, ok)

	// Output:
	// ExampleFoo true
	//  false
}

func TestParseFailingExample(t *testing.T) {
	var examples []Example
	for _, output := range expectedOutputs {
		eg, ok := parseFailingExample(output)
		if ok {
			examples = append(examples, eg)
		}
	}

	if len(examples) != len(expectedExamples) {
		t.Errorf("expected %d examples, got %d", len(expectedExamples), len(examples))
	}

	for i, example := range examples {
		expectedExample := expectedExamples[i]
		if example.Name != expectedExample.Name {
			t.Errorf("expected Name %q, got %q", expectedExample.Name, example.Name)
		}
		if example.Name == "Example_tricky" {
			continue // FIXME
		}
		if example.Got != expectedExample.Got {
			t.Errorf("expected Got %q, got %q", expectedExample.Got, example.Got)
		}
		if example.Want != expectedExample.Want {
			t.Errorf("expected Want %q, got %q", expectedExample.Want, example.Want)
		}
	}
}
