package lib

import (
	"github.com/rs/zerolog/log"
	"github.com/thediveo/enumflag"
	_ "gopkg.in/yaml.v3"
)

type GenerateConfig struct {
	CreateTables bool   `yaml:"createTables"`
	Dialect      string `yaml:"dialect"`
	Output       string `yaml:"output"`
}

type ConfigFile struct {
	Generate GenerateConfig `yaml:"generate"`
}

type Config struct {
	CreateTables bool
	Dialect      SQLDialect
	Output       string
}

func (config *Config) FromConfigFile(configFile *ConfigFile) {
	config.CreateTables = configFile.Generate.CreateTables
	dialectValue := enumflag.New(&config.Dialect, configFile.Generate.Dialect, SQLDialectIds, enumflag.EnumCaseInsensitive)
	sqlDialect, castSuccessful := dialectValue.Get().(SQLDialect)
	if !castSuccessful {
		log.Fatal().Msgf("Could not parse SQL dialect: %s", configFile.Generate.Dialect)
	}
	config.Dialect = sqlDialect
	config.Output = configFile.Generate.Output
}

func NewConfig() *Config {
	return &Config{
		Dialect:      SQLite,
		Output:       "sql",
		CreateTables: false,
	}
}
