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

package utils_test

import (
	"testing"
	"time"

	"github.com/spacemonkeygo/flagfile"
	"github.com/spacemonkeygo/flagfile/utils"
)

func ExampleSetup(t *testing.T) {
	// make a config var to hold some flags
	var Config struct {
		Value int           `default:"20" usage:"this is some value"`
		Time  time.Duration `default:"10s" usage:"some duration"`
		What  string        `default:"such a string" usage:"quite the string"`
	}

	// this creates flags with names
	//     prefix.value
	//     prefix.time
	//     prefix.what
	utils.Setup("prefix", &Config)

	// this loads the flags into config. flagfile automatically adds a flag
	// named flagfile, which when specified will load the key/value pairs in
	// that file as if they were specified on the command line.
	flagfile.Load()
}
