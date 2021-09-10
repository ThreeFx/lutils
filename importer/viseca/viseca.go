// This file is part of lutils, utility programs for plaintext accounting
//
// Copyright 2021 Ben Fiedler
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package viseca

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"git.sr.ht/~bfiedler/lutil/importer/types"
)

const (
	urlPre  = "https://api.one.viseca.ch/v1/card/"
	urlPost = "/transactions?stateType=unknown&offset=0&pagesize=1000&dateFrom=2021-01-01"

	ZKBKreditkarte = "ZKB Kreditkarte"

	ID           = 0
	Datum        = 1
	MerchantName = 2
	Details      = 3
	Amount       = 4
)

type Transaction struct {
	id           string
	date         time.Time
	merchantName string
	details      string
	amount       float64
}

type Transactions struct {
	Transactions []Transaction `json:"list"`
}

type Importer struct{}

func (i Importer) Name() string { return "Viseca" }

func (i Importer) Parse(inp io.Reader) ([]types.Transaction, error) {
	in := csv.NewReader(inp)
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
			log.Printf("read transaction \"%s\"", t.merchantName)
			ts = append(ts, t)
		}

		r, err = in.Read()
	}

	return ts, nil
}

// may return nil
func transactionFromRecord(r []string) *Transaction {
	log.Printf("%#v\n", r)
	date, err := time.Parse("2006-01-02T15:04:05", r[Datum])
	if err != nil {
		log.Fatalf("error parsing time: %v", err)
	}

	amt, err := strconv.ParseFloat(r[Amount], 64)
	if err != nil {
		log.Fatalf("error parsing amount: %v", err)
	}

	if r[Details] == "Ihre Zahlung - Danke" {
		// ignore balance payments, these are tracked from the ZKB side
		return nil
	}

	return &Transaction{
		id:           r[ID],
		date:         date,
		merchantName: r[MerchantName],
		details:      r[Details],
		amount:       amt,
	}
}

func (t *Transaction) ID() string {
	return t.id
}

func (t *Transaction) WriteTo(w io.Writer) (n int64, err error) {
	date := t.date.Format("2006/01/02")
	m, err := fmt.Fprintf(w, "%s *\n", date)
	n += int64(m)
	if err != nil {
		return n, fmt.Errorf("error writing transaction: %v", err)
	}

	m, err = fmt.Fprintf(w, "    ; %s\n", t.merchantName)
	n += int64(m)
	if err != nil {
		return n, fmt.Errorf("error writing transaction: %v", err)
	}

	//m, err = fmt.Fprintf(w, "    ; details \"%s\"\n", t.details)
	//n += int64(m)
	//if err != nil {
	//	return n, fmt.Errorf("error writing transaction: %v", err)
	//}

	m, err = fmt.Fprintf(w, "    %-50s  %5.2f CHF\n", ZKBKreditkarte, -t.amount)
	n += int64(m)
	if err != nil {
		return n, fmt.Errorf("error writing transaction: %v", err)
	}

	m, err = fmt.Fprintf(w, "    %-50s  %5.2f CHF\n", "TODO", t.amount)
	n += int64(m)
	if err != nil {
		return n, fmt.Errorf("error writing transaction: %v", err)
	}

	return n, nil
}
