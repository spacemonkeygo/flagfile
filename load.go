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

var flagfile = flag.String("flagfile", "", "a file from which to load flags")

func SetFlagfile(file *string) {
	flagfile = file
}

func Load() {
	flag.Parse()
	if len(*flagfile) == 0 {
		return
	}

	set_flags := map[string]bool{}
	flag.Visit(func(f *flag.Flag) {
		set_flags[f.Name] = true
	})

	fh, err := os.Open(*flagfile)
	if err != nil {
		panic(fmt.Sprintf("unable to open flagfile '%s': %s", *flagfile, err))
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
		panic(fmt.Sprintf("unable to read flagfile '%s': %s", *flagfile, err))
	}
}
