// The generate command is used to scaffold both CREATE and INSERT
// statements
package generate

import (
	"fmt"
	"github.com/spf13/cobra"
	"squeak/lib"
)

var GenerateCmd = cobra.Command{
	Use:   "generate",
	Short: "Generate CREATE and INSERT statements",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("sql dialect is: %d=%q\n",
			lib.CurrentSQLDialect,
			cmd.Flags().Lookup("dialect").Value.String())
	},
}

func init() {
}
