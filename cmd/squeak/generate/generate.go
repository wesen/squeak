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
		ctx := cmd.Context()
		config, cast := ctx.Value("config").(*lib.Config)
		if !cast {
			cmd.PrintErrf("Could not get config\n")
			return
		}

		if cmd.Flags().Changed("create-tables") {
			createTables, _ := cmd.Flags().GetBool("create-tables")
			config.CreateTables = createTables
		}

		fmt.Printf("sql dialect is: %d=%q, output type is: %d=%q, createTables: %v\n",
			config.Dialect,
			cmd.Flags().Lookup("dialect").Value.String(),
			config.Output,
			cmd.Flags().Lookup("output").Value.String(),
			config.CreateTables,
		)
	},
}

func init() {
	GenerateCmd.PersistentFlags().Bool("create-tables", false, "Generate CREATE TABLE statements")
}
