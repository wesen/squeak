package main

import (
	"context"
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

	config := lib.NewConfig()

	// we accept overrides for the different settings under `generate`
	// in the YAML file.
	// TODO(manuel): The problem here is that we will override the config settings if we pass
	// a pointer to the variable directly
	rootCmd.PersistentFlags().VarP(
		enumflag.New(&config.Dialect, "dialect",
			SQLDialectIds, enumflag.EnumCaseInsensitive),
		"dialect", "d",
		"Dialect to use for the generated SQL statements")

	rootCmd.AddCommand(&generate.GenerateCmd)

	ctx := context.WithValue(context.TODO(), "config", lib.NewConfig())

	_ = rootCmd.ExecuteContext(ctx)
}
