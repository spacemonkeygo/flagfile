// Copyright (C) 2013-2014 Space Monkey, Inc.

package flagfile

import (
	"flag"
	"fmt"
)

var (
	all_aliases = make(map[string][]string)
	alias_set   = make(map[string]bool)
)

func Alias(new_flag_name, old_flag_name string) {
	mtx.Lock()
	defer mtx.Unlock()
	if loaded {
		panic(fmt.Errorf("flags already loaded"))
	}
	alias_set[new_flag_name] = true
	all_aliases[old_flag_name] = append(
		all_aliases[old_flag_name], new_flag_name)
}

func isAlias(flag_name string) bool {
	return alias_set[flag_name]
}

func IsAlias(flag_name string) bool {
	mtx.Lock()
	defer mtx.Unlock()
	return isAlias(flag_name)
}

func mustLookup(flag_name string) *flag.Flag {
	val := flag.Lookup(flag_name)
	if val == nil {
		panic(fmt.Errorf("flag %#v doesn't exist", flag_name))
	}
	return val
}

func defineAliases() {
	for flag_name, flag_aliases := range all_aliases {
		existing_flag := flag.Lookup(flag_name)
		if existing_flag == nil {
			panic(fmt.Errorf("alias defined pointing to a non-existent flag %#v",
				flag_name))
		}
		for _, alias := range flag_aliases {
			flag.Var(existing_flag.Value, alias, existing_flag.Usage)
		}
	}
}

func setAliases() {
	for flag_name, flag_aliases := range all_aliases {
		flags := append([]string{flag_name}, flag_aliases...)
		var set_alias *string
		for _, flag := range flags {
			if !set_flags[flag] {
				continue
			}
			if set_alias != nil {
				panic(fmt.Errorf("multiple aliases of flag %#v set", flag_name))
			}
			set_alias = &flag
		}
		if set_alias == nil {
			continue
		}
		set_alias_val := mustLookup(*set_alias).Value.String()
		for _, flag := range flags {
			if set_flags[flag] {
				continue
			}
			mustSet(flag, set_alias_val)
		}
	}
}
