// This file is part of lutils, utility programs for plaintext accounting
//
// Copyright 2021 Ben Fiedler
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package fmt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/spf13/pflag"
)

var (
	commentStart = regexp.MustCompile(`^(\s*)(; .*)$`)
	hasAmount    = regexp.MustCompile(`^\s+(?P<account>(?:[^ ]+ ??)+)(?: {2,}(?P<amount>-?\d.*))?$`)

	accountPadding *int
	amountPadding  *int

	overwriteFiles *bool
)

func initFlags() {
	accountPadding = pflag.Int("accountPadding", 40, "Account length to pad for")
	amountPadding = pflag.Int("amountPadding", 13, "Amount length to pad for")
	overwriteFiles = pflag.BoolP("overwriteFiles", "w", false, "Overwrite input files")
}

func RunFmt() {
	initFlags()
	pflag.Parse()

	if pflag.NArg() == 0 {
		processFile("<standard input>", os.Stdin, os.Stdout)
	} else {
		for _, file := range pflag.Args() {
			processFile(file, nil, os.Stdout)
		}
	}
}

func processFile(filename string, in io.Reader, out io.Writer) error {
	var perm fs.FileMode = 0644
	if in == nil {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		fi, err := f.Stat()
		if err != nil {
			return err
		}
		in = f
		perm = fi.Mode().Perm()
	}

	content, err := io.ReadAll(in)
	if err != nil {
		return err
	}

	buf, err := format(content)
	if err != nil {
		return err
	}

	if !bytes.Equal(buf, content) {
		if *overwriteFiles {
			backup, err := backupFile(filename+".", content, perm)
			if err != nil {
				return err
			}

			err = os.WriteFile(filename, buf, perm)
			if err != nil {
				os.Rename(backup, filename)
				return err
			}

			err = os.Remove(backup)
			if err != nil {
				return err
			}
		}
	}

	if !*overwriteFiles {
		_, err = out.Write(buf)
	}

	return err
}

func format(in []byte) ([]byte, error) {
	s := bufio.NewScanner(bytes.NewBuffer(in))
	var b bytes.Buffer
	out := bufio.NewWriter(&b)

	for s.Scan() {
		line := s.Text()

		var err error
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
			return nil, err
		}
	}

	err := s.Err()
	if err != nil {
		return nil, err
	}

	err = out.Flush()
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func backupFile(name string, content []byte, perm fs.FileMode) (string, error) {
	f, err := os.CreateTemp(filepath.Dir(name), filepath.Base(name))
	if err != nil {
		return "", err
	}

	backup := f.Name()

	err = f.Chmod(perm)
	if err != nil {
		f.Close()
		os.Remove(backup)
		return backup, err
	}

	_, err = f.Write(content)
	if err1 := f.Close(); err == nil {
		err = err1
	}

	return backup, err
}
