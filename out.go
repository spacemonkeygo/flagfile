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
	"io"
	"log"
	"os"

	"github.com/spacemonkeygo/flagfile/parser"
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
// DumpMyFlag will write all flags contained in the flags array to the given io.Writer in the flagfile
// serialization format for later parsing.
func DumpMyFlags(myflags []string,out io.Writer) error{
	mtx.Lock()
	defer mtx.Unlock()
	vals := make(map[string]string)

	flag.VisitAll(func (f *flag.Flag){
		for _,flg := range myflags {
			if f.Name == flg {
				vals[f.Name] = f.Value.String()
			}
		}

	})
	return parser.Serialize(vals,out)
}

func DumpMyFlagsToPath(myflags []string,path string) error{
	fh, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0600)
	if err != nil {
		return err
	}
	defer fh.Close()

	return DumpMyFlags(myflags,fh)
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