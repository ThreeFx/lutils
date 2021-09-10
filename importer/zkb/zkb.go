// This file is part of lutils, utility programs for plaintext accounting
//
// Copyright 2021 Ben Fiedler
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package zkb

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"git.sr.ht/~bfiedler/lutil/importer/types"
)

const (
	Datum        int = 0
	Buchungstext     = 1
	ZKBReferenz      = 2
	Belastung        = 4
	Gutschrift       = 5
)

const formatString = `%s
    %-50s%s
    %-50s%s
`

type Importer struct{}

type Transaction struct {
	date         time.Time // in format: dd.mm.yyyy
	buchungstext string
	zkbReferenz  string
	belastung    string // in CHF, format 1000.00
	gutschrift   string // in CHF, format 1000.00
}

func (i Importer) Name() string { return "ZKB" }

func (i Importer) Parse(inp io.Reader) ([]types.Transaction, error) {
	in := csv.NewReader(inp)
	in.Comma = ';'
	in.LazyQuotes = true

	_, err := in.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading csv header: %w", err)
	}

	var ts []types.Transaction

	r, err := in.Read()
	for err != io.EOF {
		if err != nil {
			return nil, err
		}

		t := transactionFromRecord(r)
		if t != nil {
			log.Printf("read transaction %s", t.buchungstext)
			ts = append(ts, t)
		}

		r, err = in.Read()
	}

	return ts, nil
}

// may return nil
func transactionFromRecord(r []string) *Transaction {
	date, err := time.Parse("02.01.2006", r[Datum])
	if err != nil {
		log.Fatalf("error parsing time: %v", err)
	}

	return &Transaction{
		date:         date,
		buchungstext: r[Buchungstext],
		zkbReferenz:  r[ZKBReferenz],
		belastung:    r[Belastung],
		gutschrift:   r[Gutschrift],
	}
}

func (t *Transaction) ID() string {
	return t.zkbReferenz
}

func (t *Transaction) WriteTo(w io.Writer, accountName string) (n int64, err error) {
	date := t.date.Format("2006/01/02")
	m, err := fmt.Fprintf(w, "%s *\n", date)
	n += int64(m)
	if err != nil {
		return n, fmt.Errorf("error writing transaction: %v", err)
	}

	m, err = fmt.Fprintf(w, "    ; Buchungstext \"%s\"\n", t.buchungstext)
	n += int64(m)
	if err != nil {
		return n, fmt.Errorf("error writing transaction: %v", err)
	}

	if t.gutschrift == "" && t.belastung == "" ||
		t.gutschrift != "" && t.belastung != "" {
		return n, errors.New("exactly one of gutschift and belastung must be set")
	}

	if t.gutschrift != "" {
		m, err = fmt.Fprintf(w, "    %-50s  %8s CHF\n", accountName, t.gutschrift)
		n += int64(m)
		if err != nil {
			return n, fmt.Errorf("error writing transaction: %v", err)
		}

		m, err = fmt.Fprintf(w, "    %-50s  %8s CHF\n", "TODO", "-"+t.gutschrift)
		n += int64(m)
		if err != nil {
			return n, fmt.Errorf("error writing transaction: %v", err)
		}
	}

	if t.belastung != "" {
		m, err = fmt.Fprintf(w, "    %-50s  %8s CHF\n", accountName, "-"+t.belastung)
		n += int64(m)
		if err != nil {
			return n, fmt.Errorf("error writing transaction: %v", err)
		}

		m, err = fmt.Fprintf(w, "    %-50s  %8s CHF\n", "TODO", t.belastung)
		n += int64(m)
		if err != nil {
			return n, fmt.Errorf("error writing transaction: %v", err)
		}
	}

	return n, nil
}
