// Copyright (C) 2013-2014 Space Monkey, Inc.

package flagfile

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/SpaceMonkeyInc/flagfile/parser"
)

var (
	flagfile = flag.String("flagfile", "", "a file (or multiple files, "+
		"comma-separated) from which to load flags")

	alias_mtx sync.Mutex
	aliases   = make(map[string]string)
	loaded    = false
	set_flags = make(map[string]bool)
)

func Alias(new_flag_name, old_flag_name string) {
	alias_mtx.Lock()
	defer alias_mtx.Unlock()
	if loaded {
		panic(fmt.Errorf("flags already loaded"))
	}
	existing, exists := aliases[new_flag_name]
	if exists {
		panic(fmt.Errorf("flag %#v is already defined to be %#v",
			new_flag_name, existing))
	}
	aliases[new_flag_name] = old_flag_name
}

func ActivelySet(flag_name string) bool {
	alias_mtx.Lock()
	defer alias_mtx.Unlock()
	return set_flags[flag_name]
}

func Load() {
	alias_mtx.Lock()
	defer alias_mtx.Unlock()
	if loaded {
		panic(fmt.Errorf("flags already loaded"))
	}

	last_alias_count := 0
	for {
		current_alias_count := len(aliases)
		if current_alias_count == 0 || current_alias_count == last_alias_count {
			break
		}
		last_alias_count = current_alias_count

		for new_alias, flag_name := range aliases {
			existing_flag := flag.Lookup(flag_name)
			if existing_flag == nil {
				continue
			}
			flag.Var(existing_flag.Value, new_alias, existing_flag.Usage)
			delete(aliases, new_alias)
		}
	}
	for _, flag_name := range aliases {
		panic(fmt.Errorf("flag %#v doesn't exist", flag_name))
	}

	flag.Parse()
	loaded = true

	for alias, flag_name := range aliases {
		alias_flag := flag.Lookup(alias)
		real_flag := flag.Lookup(flag_name)
		if alias_flag == nil {
			panic(fmt.Errorf("flag %#v doesn't exist", alias))
		}
		if real_flag == nil {
			panic(fmt.Errorf("flag %#v doesn't exist", flag_name))
		}
		if alias_flag.Value.String() != real_flag.Value.String() {
			panic(fmt.Errorf("aliased flags %#v and %#v don't match",
				flag_name, alias))
		}
	}

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
			if cmdline_set_flags[name] {
				return
			}
			err := flag.Set(name, value)
			if err != nil {
				panic(fmt.Errorf("Unable to set flag '%s' to '%s': %s",
					name, value, err))
			}
			set_flags[name] = true
		})
		fh.Close()
		if err != nil {
			panic(fmt.Errorf("'%s': %s", file, err))
		}
	}

	for alias, flag_name := range aliases {
		if set_flags[flag_name] {
			set_flags[alias] = true
		} else if set_flags[alias] {
			set_flags[flag_name] = true
		}
	}
}
