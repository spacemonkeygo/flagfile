//
// Copyright (C) 2013 Space Monkey, Inc.
//

package flagfile

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
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

	files := strings.Split(*flagfile, ",")
	for _, file := range files {
		if len(file) == 0 {
			continue
		}
		func() {
			fh, err := os.Open(file)
			if err != nil {
				panic(fmt.Sprintf("unable to open flagfile '%s': %s", file, err))
			}
			defer fh.Close()

			var section string
			scanner := bufio.NewScanner(fh)
			for scanner.Scan() {
				option := strings.TrimSpace(scanner.Text())
				if len(option) == 0 || option[0] == '#' {
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
					panic(fmt.Sprintf("Unable to parse flagfile line '%s'", option))
				}
				name := strings.TrimSpace(parts[0])
				if len(section) != 0 {
					name = section + name
				}
				if cmdline_set_flags[name] {
					continue
				}
				value := strings.TrimSpace(parts[1])
				err := flag.Set(name, value)
				if err != nil {
					panic(fmt.Sprintf("Unable to set flag '%s' to '%s': %s",
						name, value, err))
				}
				set_flags[name] = true
			}
			if err = scanner.Err(); err != nil {
				panic(fmt.Sprintf("unable to read flagfile '%s': %s", file, err))
			}
		}()
	}

	for alias, flag_name := range aliases {
		if set_flags[flag_name] {
			set_flags[alias] = true
		} else if set_flags[alias] {
			set_flags[flag_name] = true
		}
	}
}
