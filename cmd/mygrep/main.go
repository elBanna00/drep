package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	DIGITS       = "\\d"
	ALPHANUMERIC = "\\w"
)

// Ensures gofmt doesn't remove the "bytes" import above (feel free to remove this!)
var _ = bytes.ContainsAny

// Usage: echo <input_text> | your_program.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	ok, err := matchLine(line, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}

	// default exit code is 0 which means success
}

func matchLine(line []byte, pattern string) (bool, error) {
	if utf8.RuneCountInString(pattern) == 0 {
		return false, fmt.Errorf("unsupported pattern: %q", pattern)
	}

	isGroup := len(pattern) > 3 && pattern[0] == '[' && pattern[len(pattern)-1] == ']'
	var ok bool
	if pattern == DIGITS {
		ok = bytes.ContainsAny(line, "0123456789")
	} else if pattern == ALPHANUMERIC {
		ok = regexp.MustCompile("^[a-zA-Z0-9_]*$").MatchString(string(line))
	} else if isGroup {
		var start int
		isNegativeGroup := pattern[1] == '^'
		if isNegativeGroup {
			start = 2
		} else {
			start = 1
		}

		accept := pattern[start : len(pattern)-1]

		if strings.IndexAny(string(line), accept) != -1 && !isNegativeGroup {
			os.Exit(0)
		} else if strings.IndexAny(string(line), accept) == -1 && isNegativeGroup {
			os.Exit(0)
		}
	} else {
		// Uncomment this to pass the first stage
		ok = bytes.ContainsAny(line, pattern)
	}

	return ok, nil
}
