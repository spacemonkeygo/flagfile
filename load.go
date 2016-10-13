// Copyright (C) 2014 Space Monkey, Inc.
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
	"sync"

	"github.com/spacemonkeygo/flagfile/parser"
)

var (
	flagfile = flag.String("flagfile", "", "a file (or multiple files, "+
		"comma-separated) from which to load flags")

	mtx       sync.Mutex
	loaded    = false
	set_flags = make(map[string]bool)
)

// IsActivelySet returns whether or not the user configured the given flag.
// The value is false if the flag was not set by commandline or flagfile.
func IsActivelySet(flag_name string) bool {
	mtx.Lock()
	defer mtx.Unlock()
	return set_flags[flag_name]
}

func mustSet(flag_name, flag_value string) {
	err := flag.Set(flag_name, flag_value)
	if err != nil {
		panic(fmt.Errorf("Unable to set flag %#v to %#v: %s",
			flag_name, flag_value, err))
	}
}

type Option struct {
	flagfilePath string
	skipArgs     bool
	short_usage  func()
	full_usage   func()
}

// Flagfile tells Load to find default values from the flagfile at path (which
// will be overridden by user-provided values from arguments, if provided).
func Flagfile(path string) Option { return Option{flagfilePath: path} }

// OptFlagfile is like Flagfile but does not fail if the path doesn't exist.
func OptFlagfile(path string) Option {
	if _, err := os.Stat(path); err == nil {
		return Option{flagfilePath: path}
	}
	return Option{}
}

// ShortUsageFunc specifies what method to run when '-h' or '--help' are found.
func ShortUsageFunc(fn func()) Option {
	return Option{short_usage: fn}
}

// FullUsageFunc specifies what method to run when '--help-all' is found.
func FullUsageFunc(fn func()) Option {
	return Option{short_usage: fn}
}

// SkipArgs will tell Load to not call flag.Parse and otherwise avoid looking
// at process arguments
func SkipArgs() Option { return Option{skipArgs: true} }

// Load is the flagfile equivalent/replacement for flag.Parse()
// Call once at program start.
func Load(opts ...Option) {
	defer flagOut()
	mtx.Lock()
	defer mtx.Unlock()
	if loaded {
		panic(fmt.Errorf("flags already loaded"))
	}

	defineAliases()

	var flagfiles []string
	var skipArgs bool
	short_usage := ShortUsage
	full_usage := FullUsage
	for _, opt := range opts {
		if opt.flagfilePath != "" {
			flagfiles = append(flagfiles, opt.flagfilePath)
		}
		if opt.skipArgs {
			skipArgs = true
		}
		if opt.short_usage != nil {
			short_usage = opt.short_usage
		}
		if opt.full_usage != nil {
			full_usage = opt.full_usage
		}
	}

	cmdline_set_flags := map[string]bool{}
	if !skipArgs {
		flag.CommandLine.Usage = short_usage
		args := os.Args[1:]
		for _, arg := range args {
			if arg == "--" {
				break
			}
			if arg == "--help-all" {
				full_usage()
				os.Exit(2)
			}
		}
		flag.CommandLine.Parse(os.Args[1:])
		flag.Visit(func(f *flag.Flag) {
			cmdline_set_flags[f.Name] = true
			set_flags[f.Name] = true
		})
	}
	loaded = true

	flagfiles = append(flagfiles, strings.Split(*flagfile, ",")...)

	for len(flagfiles) > 0 {
		file := flagfiles[0]
		flagfiles = flagfiles[1:]
		if len(file) == 0 {
			continue
		}
		fh, err := os.Open(file)
		if err != nil {
			panic(fmt.Errorf("unable to open flagfile '%s': %s", file, err))
		}
		err = parser.Parse(fh, func(name, value string) {
			if name == "flagfile" {
				// allow flagfile chaining
				flagfiles = append(flagfiles, value)
				return
			}
			// command line flags override file flags
			if cmdline_set_flags[name] {
				return
			}
			mustSet(name, value)
			set_flags[name] = true
		})
		fh.Close()
		if err != nil {
			panic(fmt.Errorf("'%s': %s", file, err))
		}
	}

	setAliases()
}
