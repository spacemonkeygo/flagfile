//
// Copyright (C) 2013 Space Monkey, Inc.
//

package parser

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

func Serialize(values map[string]string, out io.Writer) error {
	keys := make([]string, 0, len(values))
	fixed_vals := make(map[string]string, len(values))
	for key, val := range values {
		if !strings.Contains(key, ".") {
			key = "main." + key
		}
		keys = append(keys, key)
		fixed_vals[key] = val
	}
	sort.Strings(keys)
	last_section := ""
	for _, key := range keys {
		idx := strings.LastIndex(key, ".")
		section := key[:idx]
		key = key[idx+1:]
		if section != last_section {
			if last_section != "" {
				_, err := fmt.Fprintf(out, "\n")
				if err != nil {
					return err
				}
			}
			_, err := fmt.Fprintf(out, "[%s]\n", section)
			if err != nil {
				return err
			}
			last_section = section
		}
		_, err := fmt.Fprintf(out, "%s = %s\n", key, fixed_vals[key])
		if err != nil {
			return err
		}
	}
	return nil
}
