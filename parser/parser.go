// Copyright (C) 2014 Space Monkey, Inc.

package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

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
