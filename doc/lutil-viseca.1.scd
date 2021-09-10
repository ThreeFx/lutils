LUTIL-VISECA(1) "sr.ht/~bfiedler/lutil" "General Commands Manual"

# NAME

lutil-viseca - save transactions from Viseca one as csv

# SYNOPSIS

*lutil viseca* [_flags_] --cardID 0000ABCDEFGHXXXX --cookie asdfasdfasdf

# DESCRIPTION

*lutil viseca* downloads credit card transactions from Viseca one and outputs
them to stdout as csv. The output format can be understood by
*lutil-import(1)*'s _viseca_ type.

# OPTIONS

*--cardID* _id_
	Required. Card ID to download transactions for. You can find this by
	naviating to the desired transaction list in your browser and copying it
	from the URL. Note: This is *not* your credit card number, but it still
	contains the last 4 digits of your credit card number, so be careful where
	you put this.

*--cookie* _cookie_
	Required. Cookie to use for the Viseca one API. Extract this from your
	browser's store after logging in to Viseca one. At time of writing it has
	the name "AL_SESS-S".

*--cookieName* _name_
	Cookie name to use for the Viseca one API. Defaults to "AL_SESS-S".

*--fromDate* _date_
	Earliest date to consider for transactions. Leave blank if you want to go
	back as far as possible. The date format is YYYY-MM-DD.

*--toDate* _date_
	Latest date to consider for transactions. Leave blank if you want to go
	back as far as possible. The date format is YYYY-MM-DD.

# EXAMPLES

Download all available transactions for card 0000ABCDEFGHXXXX.

	ledger viseca --cardID 0000ABCDEFGHXXXX --cookie asdfasdfasdf

# NOTES

This client sets the page size to 1000 internally. If you have more than 1000
transactions then not all might be returned. Patches are welcome.

*lutil-viseca(1)* uses internal Viseca APIs that may change at any point
without notice.

# AUTHORS

Intially developed by Tobias Nehrlich<tobi@anothertobi.io>, adapted into lutils
and maintained by Ben Fiedler <git@bfiedler.ch>.
