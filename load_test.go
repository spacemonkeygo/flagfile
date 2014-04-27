// Copyright (C) 2013 Space Monkey, Inc.

package flagfile

import (
	"testing"
	"time"
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
	Setup("prefix", &Config)

	// this loads the flags into config. flagfile automatically adds a flag
	// named flagfile, which when specified will load the key/value pairs in
	// that file as if they were specified on the command line.
	Load()
}
