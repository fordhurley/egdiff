package main

import "regexp"

var lineEndingRE = regexp.MustCompile(`\r?\n`)

// replaceLineEndings changes \r\n to \n. It's idempotent.
func replaceLineEndings(raw string) string {
	return lineEndingRE.ReplaceAllString(raw, "\n")
}
