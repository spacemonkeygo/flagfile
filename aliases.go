// Copyright (C) 2013-2014 Space Monkey, Inc.

package flagfile

import (
	"flag"
	"fmt"
)

var (
	aliases = make(map[string]string)
)

func Alias(new_flag_name, old_flag_name string) {
	mtx.Lock()
	defer mtx.Unlock()
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

func mustLookup(flag_name string) *flag.Flag {
	val := flag.Lookup(flag_name)
	if val == nil {
		panic(fmt.Errorf("flag %#v doesn't exist", flag_name))
	}
	return val
}

func defineAliases() {
	// TODO: if we set up a flag alias dependency tree, we wouldn't have to loop
	//  around looking for a fixed point
	alias_copy := make(map[string]string, len(aliases))
	for alias, flag_name := range aliases {
		alias_copy[alias] = flag_name
	}
	last_alias_count := 0
	for {
		current_alias_count := len(alias_copy)
		if current_alias_count == 0 || current_alias_count == last_alias_count {
			break
		}
		last_alias_count = current_alias_count

		for new_alias, flag_name := range alias_copy {
			existing_flag := flag.Lookup(flag_name)
			if existing_flag == nil {
				continue
			}
			flag.Var(existing_flag.Value, new_alias, existing_flag.Usage)
			delete(alias_copy, new_alias)
		}
	}
	for _, flag_name := range alias_copy {
		panic(fmt.Errorf("flag %#v doesn't exist", flag_name))
	}
}

func setAliases() {
	// TODO: if we set up a flag alias dependency tree, we wouldn't have to loop
	//  around looking for a fixed point
	alias_copy := make(map[string]string, len(aliases))
	for alias, flag_name := range aliases {
		alias_copy[alias] = flag_name
	}
	last_alias_count := 0
	for {
		current_alias_count := len(alias_copy)
		if current_alias_count == 0 || current_alias_count == last_alias_count {
			break
		}
		last_alias_count = current_alias_count

		for alias, flag_name := range alias_copy {
			if set_flags[flag_name] {
				if set_flags[alias] {
					// both set
					if mustLookup(flag_name).Value.String() !=
						mustLookup(alias).Value.String() {
						panic(fmt.Errorf("aliased flags %#v and %#v both set and mismatch",
							flag_name, alias))
					}
					delete(alias_copy, alias)
					continue
				} else {
					// flag set but alias unset
					mustSet(alias, mustLookup(flag_name).Value.String())
					set_flags[alias] = true
					delete(alias_copy, alias)
					continue
				}
			} else {
				if set_flags[alias] {
					// alias set but flag unset
					mustSet(flag_name, mustLookup(alias).Value.String())
					set_flags[flag_name] = true
					delete(alias_copy, alias)
					continue
				} else {
					// neither set
					continue
				}
			}
		}
	}
	for alias, flag_name := range alias_copy {
		if mustLookup(flag_name).Value.String() !=
			mustLookup(alias).Value.String() {
			panic(fmt.Errorf("aliased flags %#v and %#v defaults mismatch?",
				flag_name, alias))
		}
	}
}
