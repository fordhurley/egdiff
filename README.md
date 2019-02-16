Output a nice diff for failing golang examples.


Turn this:


    $ go test ./...
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
    FAIL

Into this:

    $ go test ./... | egdiff
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
    diff:
    2c2
    < "a\nb/nc"
    ---
    > "a\nb\nc"
    FAIL

(or something even better)


The example formatting code is here: https://golang.org/src/testing/example.go

Based on that, the line prefix `--- FAIL: Example` should be enough to identify
the beginning of a failing test. This means the tool could probably work well as
a simple pipeline:

    go test ./... | egdiff


Slightly fancier would be to run the tool directly, like:

    egdiff ./...

And it would run only the examples and format the output. That seems
diminishingly valuable, though. Formatting the output is the one thing this tool
does.
