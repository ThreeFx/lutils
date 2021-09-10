// This file is part of lutils, utility programs for plaintext accounting
//
// Copyright 2021 Ben Fiedler
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package zkb

import (
	"strings"
	"testing"
)

var (
	header     = `"Datum";"Buchungstext";"ZKB-Referenz";"Referenznummer";"Belastung CHF";"Gutschrift CHF";"Valuta";"Saldo CHF"`
	belastung  = `"31.12.2020";"Belastung";"X000000000000000-1";"";"10.00";"";"31.12.2020";"1234.56"`
	gutschrift = `"01.01.2021";"Gutschift";"X000000000000000-2";"";"";"10.00";"01.01.2021";"9876.54"`
)

var zkb = &Importer{}

func TestParse(t *testing.T) {
	tests := map[string]struct {
		content            string
		expectedNum        int
		expectedIDs        []string
		expectedBelastung  []string
		expectedGutschrift []string
	}{
		"header only":    {header, 0, nil, nil, nil},
		"eine buchung":   {strings.Join([]string{header, belastung}, "\n"), 1, []string{"X000000000000000-1"}, []string{"10.00"}, []string{""}},
		"zwei buchungen": {strings.Join([]string{header, belastung, gutschrift}, "\n"), 2, []string{"X000000000000000-1", "X000000000000000-2"}, []string{"10.00", ""}, []string{"", "10.00"}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ts, err := zkb.Parse(strings.NewReader(tc.content))
			if err != nil {
				t.Error(err)
			}

			if tc.expectedNum != len(ts) {
				t.Errorf("%s: wrong number of transactions, expected: %v, got: %v", name, tc.expectedNum, len(ts))
			}

			for i := range ts {
				tr := ts[i].(*Transaction)
				if tc.expectedIDs[i] != tr.ID() {
					t.Errorf("%s: wrong id, expected: %v, got: %v", name, tc.expectedIDs[i], tr.ID())
				}

				if tc.expectedBelastung[i] != tr.belastung {
					t.Errorf("%s: wrong belastung, expected: %v, got: %v", name, tc.expectedBelastung[i], tr.belastung)
				}

				if tc.expectedGutschrift[i] != tr.gutschrift {
					t.Errorf("%s: wrong gutschift, expected: %v, got: %v", name, tc.expectedGutschrift[i], tr.gutschrift)
				}
			}
		})
	}
}

// TODO: test WriteTo
