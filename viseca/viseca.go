// This file is part of lutils, utility programs for plaintext accounting
//
// Copyright 2021 Tobias Nehrlich
// Copyright 2021 Ben Fiedler
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package viseca

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/spf13/pflag"
)

const (
	urlPre  = "https://api.one.viseca.ch/v1/card/"
	urlPost = "/transactions?stateType=unknown&offset=0&pagesize=1000"
)

var (
	cardID     *string
	cookie     *string
	cookieName *string
	fromDate   *string
	toDate     *string
)

type Transaction struct {
	ID           string  `json:"transactionId"`
	Date         string  `json:"date"`
	Amount       float64 `json:"amount"`
	MerchantName string  `json:"merchantName"`
	Details      string  `json:"details"`
}

type Transactions struct {
	Transactions []Transaction `json:"list"`
}

func initFlags() {
	cardID = pflag.String("cardID", "", "Card ID to use")
	cookie = pflag.String("cookie", "", "Viseca One cookie (AL_SESSeS=...), excluding the AL_SESS-S part")
	cookieName = pflag.String("cookieName", "AL_SESS-S", "cookie name for Viseca one")

	fromDate = pflag.String("fromDate", "", "earliest date to look for transactions")
	toDate = pflag.String("toDate", "", "latest date to look for transactions")
}

func RunViseca() {
	initFlags()
	pflag.Parse()

	if *cardID == "" || *cookie == "" {
		log.Fatalf("--cardID and --cookie are required")
	}

	client := &http.Client{}
	url := urlPre + *cardID + urlPost
	if *fromDate != "" {
		url += "&dateFrom=" + *fromDate
	}
	if *toDate != "" {
		url += "&dateTo=" + *toDate
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Cookie", *cookieName+"="+*cookie)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var ts Transactions
	err = json.Unmarshal(data, &ts)
	if err != nil {
		log.Fatal(err)
	}
	reverse(ts.Transactions)

	fmt.Println("\"ID\",\"Date\",\"Payee\",\"Details\",\"Amount\"")
	for _, v := range ts.Transactions {
		fmt.Printf("\"%s\",\"%s\",\"%s\",\"%s\",\"%.02f\"\n", v.ID, v.Date, v.MerchantName, v.Details, v.Amount)
	}
}

func reverse(a []Transaction) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}
