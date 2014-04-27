// Copyright (C) 2014 Space Monkey, Inc.

package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Parse takes an io.Reader and calls the given calback for each key and
// unparsed value. This format is used by flagfile for loading and serializing
// flags.
//
// The file format understands comments (provided the line starts with
// a comment character ('#' or ';')), and essentially has macros for prefixes.
// If a line that starts with '[' and ends with ']' is found, that is used as
// the key prefix for the remaining keys until the next prefix macro. Prefixes
// are automatically joined with a '.'. The special prefix "main" stops
// prefix handling.
//
// An example file:
//
//    some.flag = 20
//    # some.other.flag = 50
//    flag3 = 10m
//    flag4 = a string value
//
//    [section1]
//    flag1 = 30
//    flag2 = 40
//
//    [section2]
//    flag1 = 50
//    flag2 = true
//
func Parse(in io.Reader, cb func(key, value string)) error {
	section := ""
	scanner := bufio.NewScanner(in)
	lineno := 0
	for scanner.Scan() {
		lineno += 1
		option := strings.TrimSpace(scanner.Text())
		if len(option) == 0 || option[0] == '#' || option[0] == ';' {
			continue
		}
		if option[0] == '[' && option[len(option)-1] == ']' {
			section = option[1:len(option)-1] + "."
			if section == "main." { // main means no section
				section = ""
			}
			continue
		}
		parts := strings.SplitN(option, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("unable to parse flagfile line %d", lineno)
		}
		name := strings.TrimSpace(parts[0])
		if section != "" {
			name = section + name
		}
		cb(name, strings.TrimSpace(parts[1]))
	}
	err := scanner.Err()
	if err != nil {
		return fmt.Errorf("unable to parse flagfile: %s", err)
	}
	return nil
}
