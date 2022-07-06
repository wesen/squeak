package main

import (
	"context"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag"
	"os"
	"squeak/cmd/squeak/generate"
	"squeak/lib"
	"time"
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
		if isatty.IsTerminal(os.Stderr.Fd()) {
			log.Debug().Msg("stderr is a terminal")
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
		} else {
			log.Debug().Msg("stderr is not a terminal")
			log.Logger = log.Output(os.Stderr)
		}

		ctx := cmd.Context()
		if ctx == nil {
			log.Fatal().Msg("No context")
		}
		config, cast := ctx.Value("config").(*lib.Config)
		if !cast {
			log.Fatal().Msg("No config set in context")
		}

		if cmd.Flags().Lookup("dialect") != nil {
			config.Dialect = overrideSQLDialect
		}
	},
}

var SQLDialectIds = map[lib.SQLDialect][]string{
	lib.SQLite:     {"sqlite"},
	lib.MySQL:      {"mysql"},
	lib.PostGreSQL: {"postgresql"},
}
var overrideSQLDialect lib.SQLDialect = lib.SQLite

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug logging")
	rootCmd.PersistentFlags().String("config", "", "Path to config file")
	_ = rootCmd.MarkFlagRequired("config")

	// we accept overrides for the different settings in the YAML file
	rootCmd.PersistentFlags().VarP(
		enumflag.New(&overrideSQLDialect, "dialect",
			SQLDialectIds, enumflag.EnumCaseInsensitive),
		"dialect", "d",
		"Dialect to use for the generated SQL statements")

	rootCmd.AddCommand(&generate.GenerateCmd)

	ctx := context.WithValue(context.TODO(), "config", lib.NewConfig())
	_ = rootCmd.ExecuteContext(ctx)
}
