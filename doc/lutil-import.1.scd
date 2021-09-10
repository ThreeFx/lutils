LUTIL-IMPORT(1) "sr.ht/~bfiedler/lutil" "General Commands Manual"

# NAME

lutil-import - converts csv files to ledger transactions

# SYNOPSIS

*lutil import* _type_ [_flags_]

# DESCRIPTION

*lutil import* converts csv files from stdin to ledger transactions and outputs
them on stdout. _type_ determins the semantics of the *lutil import* csv file.

# OPTIONS

*-l* _filename_, *--ledger* _filename_
	Use _filename_ to determine whether transactions have already been booked.
	Booked transactions are skipped during import.

# EXAMPLES

Read a Kontoauszug and append ledger transactions to _ledger.leg_ using the
_zkb_ importer:

	ledger import zkb < Kontoauszug\ XXXXXX.csv >> ledger.leg

Use an existing ledger file to confirm that no new transactions are pending:

	ledger import zkb -l myledger.leg < Kontoauszug\ XXXXXX.csv

# AUTHORS

Developed and maintained by Ben Fiedler <git@bfiedler.ch>
