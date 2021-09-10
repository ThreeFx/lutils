// This file is part of lutils, utility programs for plaintext accounting
//
// Copyright 2021 Ben Fiedler
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package types

import (
	"io"
)

type Importer interface {
	// Name of the given importer. Used to distinguish transactions from
	// different importers.
	Name() string

	// Parses a list of transactions from this account, and returns them.
	Parse(r io.Reader) ([]Transaction, error)
}

type Transaction interface {
	// Return the transaction's unique identifier. Identifiers may be
	// reused across different importers.
	ID() string

	// Writes a transaction to w in a ledger-compatible format
	WriteTo(w io.Writer) (int64, error)
}
