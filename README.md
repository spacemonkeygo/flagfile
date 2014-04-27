flagfile
========

If you use flagfile.Load() instead of flag.Parse() in your main method, not
only can you continue to use flags as you always have, but you can specify
a --flagfile option and store your flags in multiple config files on disk.

flagfile provides the following additional config options:

	--flagfile: a comma-separated list of paths to load
	--flagout: writes all configured flag settings to this path after load

Precedence
----------

1. A flag set on the command line has the highest precedence.
2. A flag set in a flagfile has the next highest.
3. The default value has the lowest.

Flagfile format
---------------

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
