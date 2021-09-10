ifeq ($(PREFIX),)
    PREFIX := /usr/local
endif

.POSIX:
.SUFFIXES:
.SUFFIXES: .scd .1 .7

PREFIX?=/usr/local
_INSTDIR=$(DESTDIR)$(PREFIX)
BINDIR?=$(_INSTDIR)/bin
SHAREDIR?=$(_INSTDIR)/share/lutil
MANDIR?=$(_INSTDIR)/share/man

all: lutil doc

lutil:
	go build

%.1: %.1.scd
	scdoc < $< > $@

%.7: %.7.scd
	scdoc < $< > $@

DOCS := \
	doc/lutil-fmt.1 \
	doc/lutil-import.1 \
	doc/lutil-viseca.1 \
	doc/lutil-import.7

doc: $(DOCS)

clean:
	$(RM) lutil doc/*.1 doc/*.7

install: $(DOCS) lutil
	install -m755 lutil $(BINDIR)/lutil
	mkdir -p $(MANDIR)/man1
	mkdir -p $(MANDIR)/man7
	install -m644 doc/lutil-fmt.1 $(MANDIR)/man1/lutil-fmt.1
	install -m644 doc/lutil-import.1 $(MANDIR)/man1/lutil-import.1
	install -m644 doc/lutil-viseca.1 $(MANDIR)/man1/lutil-viseca.1
	install -m644 doc/lutil-import.7 $(MANDIR)/man7/lutil-import.7

uninstall:
	$(RM) $(BINDIR)/lutil
	$(RM) $(MANDIR)/man1/lutil-fmt.1
	$(RM) $(MANDIR)/man1/lutil-import.1
	$(RM) $(MANDIR)/man1/lutil-viseca.1
	$(RM) $(MANDIR)/man7/lutil-import.7
	$(RMDIR_IF_EMPTY) $(MANDIR)/man1
	$(RMDIR_IF_EMPTY) $(MANDIR)/man7

.PHONY: all doc clean install uninstall
