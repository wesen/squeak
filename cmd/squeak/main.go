package main

import (
	"github.com/spf13/cobra"
	"squeak/cmd/squeak/generate"
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

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug logging")

	rootCmd.AddCommand(&generate.GenerateCmd)

	_ = rootCmd.Execute()
}
