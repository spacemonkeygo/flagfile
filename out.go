// Copyright (C) 2014 Space Monkey, Inc.

package flagfile

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/SpaceMonkeyGo/flagfile/parser"
)

var (
	flagOutPath = flag.String("flagout", "",
		"a file in which to write all configured settings")
)

func flagOut() {
	if *flagOutPath != "" {
		err := DumpToPath(*flagOutPath)
		if err != nil {
			log.Printf("failed writing requested flagout file: %s", err)
		}
	}
}

// Dump will write all configured flags to the given io.Writer in the flagfile
// serialization format for later parsing.
func Dump(out io.Writer) error {
	mtx.Lock()
	defer mtx.Unlock()
	vals := make(map[string]string)
	flag.VisitAll(func(f *flag.Flag) {
		if !isAlias(f.Name) {
			vals[f.Name] = f.Value.String()
		}
	})
	return parser.Serialize(vals, out)
}

// DumpToPath simply calls Dump on a new filehandle (O_CREATE|O_TRUNC) for the
// given path
func DumpToPath(path string) error {
	fh, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0600)
	if err != nil {
		return err
	}
	defer fh.Close()
	return Dump(fh)
}
