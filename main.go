// This file is part of lutils, utility programs for plaintext accounting
//
// Copyright 2021 Ben Fiedler
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	ledgerFmt "git.sr.ht/~bfiedler/lutil/fmt"
	"git.sr.ht/~bfiedler/lutil/importer"
	"git.sr.ht/~bfiedler/lutil/viseca"
)

var commands = map[string]func(){
	"fmt":    ledgerFmt.RunFmt,
	"import": importer.RunImport,
	"viseca": viseca.RunViseca,
}

func main() {
	if len(os.Args) < 2 {
		printCommands()
		log.Fatal("expect at least argument 'command'")
	}
	cmd := os.Args[1]
	os.Args = os.Args[1:]
	f, ok := commands[cmd]
	if !ok {
		printCommands()
		log.Fatalf("command \"%s\" not found", cmd)
	}
	f()
}

func printCommands() {
	fmt.Fprint(os.Stderr, "The following commands are available: ")
	cmds := make([]string, 0, len(commands))
	for cmd := range commands {
		cmds = append(cmds, cmd)
	}
	fmt.Fprintf(os.Stderr, strings.Join(cmds, ", "))
	fmt.Fprintln(os.Stderr)
}
