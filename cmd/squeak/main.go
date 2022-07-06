package main

import (
	"context"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"squeak/cmd/squeak/generate"
	"squeak/lib"
	"time"
)
import "github.com/rs/zerolog"

func loadConfigFile(path string) (*lib.ConfigFile, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal().Err(err).Msgf("Could not open config file %s", path)
	}
	defer func() {
		_ = file.Close()
	}()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal().Err(err).Msgf("Could not read config file %s", path)
	}

	configFile := &lib.ConfigFile{}

	err = yaml.Unmarshal(bytes, &configFile)
	if err != nil {
		log.Fatal().Err(err).Msgf("Could not parse config file %s", path)
	}

	return configFile, nil
}

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

		// parse config file
		configValue := cmd.Flags().Lookup("config").Value
		if configValue == nil || configValue.String() == "" {
			log.Fatal().Msg("No config file given")
		}
		configFile, err := loadConfigFile(configValue.String())
		if err != nil {
			log.Fatal().Err(err).Msg("Could not load config file")
		}
		config.FromConfigFile(configFile)

		// handle overrides
		dialectValue := cmd.Flags().Lookup("dialect")
		if dialectValue != nil {
			config.Dialect = overrideSQLDialect
		}

		outputTypeValue := cmd.Flags().Lookup("output")
		if outputTypeValue != nil {
			config.Output = overrideOutputType
		}

	},
}

var overrideSQLDialect lib.SQLDialect = lib.SQLite
var overrideOutputType lib.OutputType = lib.OutputTypeSQL

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug logging")
	rootCmd.PersistentFlags().String("config", "", "Path to config file")
	err := rootCmd.MarkPersistentFlagRequired("config")
	if err != nil {
		log.Fatal().Err(err).Msg("Could not mark config flag required")
	}

	// we accept overrides for the different settings in the YAML file
	rootCmd.PersistentFlags().VarP(
		enumflag.New(&overrideSQLDialect, "dialect",
			lib.SQLDialectIds, enumflag.EnumCaseInsensitive),
		"dialect", "d",
		"Dialect to use for the generated SQL statements (sqlite, mysql, postgresql)")
	rootCmd.PersistentFlags().VarP(
		enumflag.New(&overrideOutputType, "output",
			lib.OutputTypeIds, enumflag.EnumCaseInsensitive),
		"output", "O",
		"Output type (SQL, CSV or SQLite)")

	rootCmd.AddCommand(&generate.GenerateCmd)

	ctx := context.WithValue(context.TODO(), "config", lib.NewConfig())
	_ = rootCmd.ExecuteContext(ctx)
}
