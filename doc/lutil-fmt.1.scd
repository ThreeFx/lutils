LUTIL-FMT(1) "sr.ht/~bfiedler/lutil" "General Commands Manual"

# NAME

lutil-fmt - format ledger files

# SYNOPSIS

*lutil fmt* [_flags_] [--] [_filenames_]

# DESCRIPTION

*lutil fmt* formats ledger files, similar to *gofmt*. By default it reads the
files specified in _filenames_ and prints formatted contents of each file to
stdout. If no filename is given, stdin is read instead.

Flags and filenames can be mixed. *lutil fmt* will continue to interpret flags
until a flag separator (--) is encountered.

# OPTIONS

*--accountPadding* _n_
	Number of characters to pad an account with. If an account name is longer
	than _n_ characters, it will overflow the line. No truncation occurs.

*--amountPadding* _n_
	Number of characters to pad an account with. If an amount is longer than
	_n_ characters, it will overflow the line. No truncation occurs.

*-w*, *--overwriteFiles*
	Overwrite the files given on the command line instead of printing them to
	stdout.

# EXAMPLES

Read stdin and output a formatted version:

	ledger fmt

Read the contents of _file1.leg_ and _file2.leg_ and replace them with a formatted version:

	ledger fmt -w file1.leg file2.leg

# AUTHORS

Developed and maintained by Ben Fiedler <git@bfiedler.ch>
