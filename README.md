`egdiff` -- add formatted diffs for failing go examples


Turn this:


    $ go test -v ./...
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
    FAIL

Into this:

    $ go test -v ./... | egdiff
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
    --- Want
    +++ Got
    @@ -1,4 +1,4 @@
    "a\n\nb\n\nc"
    -"a\nb/nc"
    +"a\nb\nc"
    "a\nb\nc"
    "abc"
    FAIL


Install

    go get -u github.com/fordhurley/egdiff


Pipe *verbose* test output to it:

    go test -v . | egdiff

**TODO:** don't require verbose flag
