// The generate command is used to scaffold both CREATE and INSERT
// statements
package generate

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag"
	"squeak/lib"
)

var GenerateCmd = cobra.Command{
	Use:   "generate",
	Short: "Generate CREATE and INSERT statements",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("sql dialect is: %d=%q\n",
			sqlDialect,
			cmd.PersistentFlags().Lookup("dialect").Value.String())
	},
}

var SQLDialectIds = map[lib.SQLDialect][]string{
	lib.SQLite:     {"sqlite"},
	lib.MySQL:      {"mysql"},
	lib.PostGreSQL: {"postgresql"},
}

var sqlDialect lib.SQLDialect

func init() {
	// TODO(manuel): this should be moved to the RootCmd
	GenerateCmd.PersistentFlags().String("config", "", "Path to config file")
	_ = GenerateCmd.MarkFlagRequired("config")

	// we accept overrides for the different settings under `generate`
	// in the YAML file
	// TODO(manuel): this should be moved to the RootCmd
	GenerateCmd.PersistentFlags().VarP(
		enumflag.New(&sqlDialect, "dialect",
			SQLDialectIds, enumflag.EnumCaseInsensitive),
		"dialect", "d",
		"Dialect to use for the generated SQL statements")
}
