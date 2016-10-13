// Copyright (C) 2016 Space Monkey, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package flagfile

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func formatFlag(f *flag.Flag) string {
	// Two spaces before -; see next two comments.
	s := fmt.Sprintf("  -%s", f.Name)
	name, usage := flag.UnquoteUsage(f)
	if len(name) > 0 {
		s += " " + name
	}
	// Boolean flags of one ASCII letter are so common we
	// treat them specially, putting their usage on the same line.
	if len(s) <= 4 { // space, space, '-', 'x'.
		s += "\t"
	} else {
		// Four spaces before the tab triggers good alignment
		// for both 4- and 8-space tab stops.
		s += "\n    \t"
	}
	s += usage
	switch f.DefValue {
	case "0", "false", "":
	default:
		s += fmt.Sprintf(" (default %#v)", f.DefValue)
	}
	return s
}

func header() { fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0]) }

func isWithheld(name string) bool {
	switch name {
	case "flagfile", "flagout":
		return true
	default:
		return false
	}
}

// ShortUsage only outputs to stderr flags without "." and withholds some
// system flags.
func ShortUsage() {
	header()
	flag.VisitAll(func(f *flag.Flag) {
		if strings.Contains(f.Name, ".") || isWithheld(f.Name) {
			return
		}
		fmt.Fprint(os.Stderr, formatFlag(f), "\n")
	})
}

// FullUsage outputs full usage information to stderr. All flags.
func FullUsage() {
	header()

	var withheld []string

	flag.VisitAll(func(f *flag.Flag) {
		if strings.Contains(f.Name, ".") {
			return
		}
		if isWithheld(f.Name) {
			withheld = append(withheld, formatFlag(f))
			return
		}
		fmt.Fprint(os.Stderr, formatFlag(f), "\n")
	})

	if len(withheld) > 0 {
		fmt.Fprintln(os.Stderr)
	}
	for _, msg := range withheld {
		fmt.Fprint(os.Stderr, msg, "\n")
	}

	current_section := ""
	flag.VisitAll(func(f *flag.Flag) {
		pos := strings.LastIndex(f.Name, ".")
		if pos == -1 {
			return
		}
		section := f.Name[:pos]
		if section != current_section {
			current_section = section
			fmt.Fprintln(os.Stderr)
		}
		fmt.Fprint(os.Stderr, formatFlag(f), "\n")
	})
}
