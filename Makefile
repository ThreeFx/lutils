ifeq ($(PREFIX),)
    PREFIX := /usr/local
endif

.PHONY: all doc clean

all: lutil doc

lutil:
	go build

doc:
	scdoc < doc/lutil-fmt.1.scd > doc/lutil-fmt.1

clean:
	rm -f lutil doc/*.1

install:
	install -m755 lutil $(PREFIX)/bin/lutil
	mkdir -p $(PREFIX)/share/man/man1
	install -m644 doc/lutil-fmt.1 $(PREFIX)/share/man/man1/lutil-fmt.1
