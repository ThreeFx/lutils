// This file is part of lutils, utility programs for plaintext accounting
//
// Copyright 2021 Ben Fiedler
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package importer

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	flag "github.com/spf13/pflag"

	"git.sr.ht/~bfiedler/lutil/importer/types"
	"git.sr.ht/~bfiedler/lutil/importer/viseca"
	"git.sr.ht/~bfiedler/lutil/importer/zkb"
)

var ledgerFile *string

var importers = map[string]types.Importer{
	"zkb":    &zkb.Importer{},
	"viseca": &viseca.Importer{},
}

func initFlags() {
	ledgerFile = flag.StringP("ledger", "l", "", "Existing ledger file (eliminate duplicate transactions)")
}

func RunImport() {
	initFlags()
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "wrong number of arguments: want: 1, got: %d\n", flag.NArg())
		flag.Usage()
		os.Exit(-1)
	}

	i := flag.Arg(0)
	importer, ok := importers[strings.ToLower(i)]
	if !ok {
		printImporters()
		fmt.Fprintf(os.Stderr, "unknown importer: %v\n", i)
	}

	ts, err := importer.Parse(os.Stdin)
	if err != nil {
		log.Fatalf("error parsing input: %v", err)
	}

	fmt.Fprintf(os.Stderr, "successfully read transactions")

	var ids map[string]struct{}

	fmt.Fprintf(os.Stderr, "ledgerFile: %s", *ledgerFile)
	if *ledgerFile != "" {
		fmt.Fprintf(os.Stderr, "reading existing ledger '%s'", *ledgerFile)
		ids = readIDs(*ledgerFile, importer.Name())
	}

	fmt.Fprintf(os.Stderr, "beginning import")
	fmt.Printf("\n; importer %s run at %s\n", importer.Name(), time.Now().Format("2006-01-02T15:04:05"))

	for _, t := range ts {
		if _, exists := ids[t.ID()]; exists {
			fmt.Fprintf(os.Stderr, "transaction %s:%s already exists, skipping", importer.Name(), t.ID())
			continue
		}

		_, err := t.WriteTo(os.Stdout)
		if err != nil {
			log.Fatalf("could not write transaction: %v", err)
		}

		_, err = fmt.Fprintf(os.Stdout, "    ; IMPORTED %s:%s\n\n", importer.Name(), t.ID())
		if err != nil {
			log.Fatalf("could not write transaction id: %v", err)
		}
	}

	fmt.Printf("; importer %s finished at %s\n", importer.Name(), time.Now().Format("2006-01-02T15:04:05"))
	fmt.Fprintf(os.Stderr, "finished import")
}

func readIDs(path string, importerName string) map[string]struct{} {
	f, err := os.Open(*ledgerFile)
	if err != nil {
		log.Fatalf("error opening ledger transactions: %v", err)
	}
	defer f.Close()

	ids := make(map[string]struct{}, 0)
	s := bufio.NewScanner(f)
	for s.Scan() {
		if s.Err() != nil {
			log.Fatalf("error scanning ledger transactions: %v", err)
		}

		line := s.Text()
		importPrefix := "    ; IMPORTED " + importerName + ":"
		if !strings.HasPrefix(line, importPrefix) {
			continue
		}

		id := strings.TrimPrefix(line, importPrefix)
		ids[id] = struct{}{}
	}

	return ids
}

func printImporters() {
	fmt.Fprintf(os.Stderr, "known importers:")
	is := make([]string, 0, len(importers))
	for i := range importers {
		is = append(is, i)
	}
	fmt.Fprintf(os.Stderr, strings.Join(is, ", "))
}
