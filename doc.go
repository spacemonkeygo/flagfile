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

/*
Package flagfile provides disk serialization and parsing extensions to the
builtin flags library

If you use flagfile.Load() instead of flag.Parse() in your main method, not
only can you continue to use flags as you always have, but you can specify
a --flagfile option and store your flags in multiple config files on disk.

flagfile provides the following additional config options:

	--flagfile: a comma-separated list of paths to load
	--flagout: writes all configured flag settings to this path after load

Precedence

A flag set on the command line has the highest precedence.
A flag set in a flagfile has the next highest.
The default value has the lowest.

Flagfile format

Flagfiles are simple text files containing lines of the form `key = val` with
optional section headers. Lines may be prefixed with a `#` to comment them out.
Example:

	some.flag = 20
	# some.other.flag = 50
	flag3 = 10m
	flag4 = a string value

	[section1]
	flag1 = 30
	flag2 = 40

	[section2]
	flag1 = 50
	flag2 = true

If a section is specified in square braces, all of the following flags up until
the next section are effectively prefixed with the section name followed by a
period.

See github.com/SpaceMonkeyGo/flagfile/parser for more information on the file
format.
*/
package flagfile
