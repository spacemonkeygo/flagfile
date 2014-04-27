flagfile
========

flagfile provides a simple api for flags on top of the flag package. It allows
you to specify a file with flags to set and also allows easy loading of flags
into structured types.

Precedence
----------

1. A flag set on the command line has the highest precedence.
2. A flag set in the flagfile has the next highest.
3. The default value has the lowest.

Usage
-----

See the examples in [godoc](http://godoc.org/github.com/SpaceMonkeyGo/flagfile)


Flagfile format
---------------

Flagfiles are simple text files containing lines of the form `key = val`. Lines
may be prefixed with a `#` to comment them out. Example:

	some.flag = 20
	# some.other.flag = 50
	flag3 = 10m
	flag4 = a string value
