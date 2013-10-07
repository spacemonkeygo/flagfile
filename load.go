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
)

var flagfile = flag.String("flagfile", "", "a file (or multiple files, "+
	"comma-separated) from which to load flags")

func Load() {
	flag.Parse()

	set_flags := map[string]bool{}
	flag.Visit(func(f *flag.Flag) {
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

			scanner := bufio.NewScanner(fh)
			for scanner.Scan() {
				option := strings.TrimSpace(scanner.Text())
				if len(option) == 0 || option[0] == '#' {
					continue
				}
				parts := strings.SplitN(option, "=", 2)
				if len(parts) != 2 {
					panic(fmt.Sprintf("Unable to parse flagfile line '%s'", option))
				}
				name := strings.TrimSpace(parts[0])
				if set_flags[name] {
					continue
				}
				value := strings.TrimSpace(parts[1])
				err := flag.Set(name, value)
				if err != nil {
					panic(fmt.Sprintf("Unable to set flag '%s' to '%s': %s",
						name, value, err))
				}
			}
			if err = scanner.Err(); err != nil {
				panic(fmt.Sprintf("unable to read flagfile '%s': %s", file, err))
			}
		}()
	}
}
