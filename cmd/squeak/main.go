package main

import (
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag"
	"squeak/cmd/squeak/generate"
	"squeak/lib"
)
import "github.com/rs/zerolog"

var rootCmd = cobra.Command{
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		debug, _ := cmd.PersistentFlags().GetBool("debug")
		if debug {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}
	},
}

var SQLDialectIds = map[lib.SQLDialect][]string{
	lib.SQLite:     {"sqlite"},
	lib.MySQL:      {"mysql"},
	lib.PostGreSQL: {"postgresql"},
}

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug logging")
	rootCmd.PersistentFlags().String("config", "", "Path to config file")
	_ = rootCmd.MarkFlagRequired("config")

	// we accept overrides for the different settings under `generate`
	// in the YAML file
	// TODO(manuel): this should be moved to the RootCmd
	rootCmd.PersistentFlags().VarP(
		enumflag.New(&lib.CurrentSQLDialect, "dialect",
			SQLDialectIds, enumflag.EnumCaseInsensitive),
		"dialect", "d",
		"Dialect to use for the generated SQL statements")

	rootCmd.AddCommand(&generate.GenerateCmd)

	_ = rootCmd.Execute()
}
