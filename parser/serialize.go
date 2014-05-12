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

package parser

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Serialize is the inverse of Parse. It automatically sorts the given keys and
// makes sections.
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
		short_key := key[idx+1:]
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
		_, err := fmt.Fprintf(out, "%s = %s\n", short_key, fixed_vals[key])
		if err != nil {
			return err
		}
	}
	return nil
}
