package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"

	"github.com/spf13/pflag"
)

var (
	commentStart = regexp.MustCompile(`^(\s*)(; .*)$`)
	hasAmount    = regexp.MustCompile(`^\s+(?P<account>(?:[^ ]+ ??)+)(?: {2,}(?P<amount>-?\d.*))?$`)

	accountPadding = pflag.Int("accountPadding", 40, "Account length to pad for")
	amountPadding  = pflag.Int("amountPadding", 13, "Amount length to pad for")

	overwriteFiles = pflag.BoolP("overwriteFiles", "w", false, "Overwrite input files")
)

func runFmt() {
	pflag.Parse()

	if pflag.NArg() == 0 {
		format(os.Stdin, os.Stdout)
	} else {
		if *overwriteFiles {
			log.Fatalf("not implemented yet")
		} else {
			for _, file := range pflag.Args() {
				r, err := os.Open(file)
				if err != nil {
					log.Fatalf("error opening file %s: %v", file, err)
				}

				format(r, os.Stdout)
			}
		}
	}
}

func format(in io.Reader, out io.Writer) error {
	s := bufio.NewScanner(in)

	var err error
	for s.Scan() {
		line := s.Text()

		if match := commentStart.FindStringSubmatch(line); match != nil {
			_, err = fmt.Fprintln(out, line)
		} else if match := hasAmount.FindStringSubmatch(line); match != nil {
			account := match[1]
			amount := match[2]

			if amount != "" {
				_, err = fmt.Fprintf(out, "    %-*s%*s\n", *accountPadding, account, *amountPadding, amount)
			} else {
				_, err = fmt.Fprintf(out, "    %s\n", account)
			}
		} else {
			_, err = fmt.Fprintln(out, line)
		}

		if err != nil {
			return err
		}
	}

	return s.Err()
}
