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

// Load is the flagfile equivalent/replacement for flag.Parse()
// Call once at program start.
func Load() {
	defer flagOut()
	mtx.Lock()
	defer mtx.Unlock()
	if loaded {
		panic(fmt.Errorf("flags already loaded"))
	}

	defineAliases()

	flag.Parse()
	loaded = true

	cmdline_set_flags := map[string]bool{}
	flag.Visit(func(f *flag.Flag) {
		cmdline_set_flags[f.Name] = true
		set_flags[f.Name] = true
	})

	for _, file := range strings.Split(*flagfile, ",") {
		if len(file) == 0 {
			continue
		}
		fh, err := os.Open(file)
		if err != nil {
			panic(fmt.Errorf("unable to open flagfile '%s': %s", file, err))
		}
		err = parser.Parse(fh, func(name, value string) {
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
