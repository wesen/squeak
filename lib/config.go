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
	dialectValue := enumflag.New(&config.Dialect, "dialect", SQLDialectIds, enumflag.EnumCaseInsensitive)
	err := dialectValue.Set(configFile.Generate.Dialect)
	if err != nil {
		log.Fatal().Msgf("Could not parse SQL dialect: %s", configFile.Generate.Dialect)
	}

	v := dialectValue.Get()
	sqlDialect, castSuccessful := v.(*SQLDialect)
	if !castSuccessful {
		log.Fatal().Msgf("Could not parse SQL dialect: %s", configFile.Generate.Dialect)
	}
	config.Dialect = *sqlDialect
	config.Output = configFile.Generate.Output
}

func NewConfig() *Config {
	return &Config{
		Dialect:      SQLite,
		Output:       "sql",
		CreateTables: false,
	}
}
