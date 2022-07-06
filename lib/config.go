package lib

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/thediveo/enumflag"
	_ "gopkg.in/yaml.v3"
)

var OutputTypeIds = map[OutputType][]string{
	OutputTypeSQL:    {"sql"},
	OutputTypeCSV:    {"csv"},
	OutputTypeSQLite: {"db"},
}

type OutputType enumflag.Flag

const (
	OutputTypeSQL OutputType = iota
	OutputTypeCSV
	OutputTypeSQLite
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
	Output       OutputType
}

func parseEnumFlag(
	flag interface{}, typename string, mapping interface{},
	value string,
) (interface{}, error) {
	flagValue := enumflag.New(flag, typename, mapping, enumflag.EnumCaseInsensitive)
	err := flagValue.Set(value)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not parse flag %s", typename)
	}
	v := flagValue.Get()
	ret, castSuccessful := v.(*SQLDialect)
	if !castSuccessful {
		return nil, errors.Wrapf(err, "Could not parse flag %s", typename)
	}

	return ret, nil
}

func (config *Config) FromConfigFile(configFile *ConfigFile) {
	config.CreateTables = configFile.Generate.CreateTables

	// crazy boilerplate
	_, err := parseEnumFlag(&config.Dialect, "SQLDialect", SQLDialectIds, configFile.Generate.Dialect)
	if err != nil {
		log.Fatal().Err(err).Msgf("Could not parse SQLDialect: %s", configFile.Generate.Dialect)
	}
	_, err = parseEnumFlag(&config.Output, "OutputType", OutputTypeIds, configFile.Generate.Output)
	if err != nil {
		log.Fatal().Err(err).Msgf("Could not parse OutputType: %s", configFile.Generate.Output)
	}
}

func NewConfig() *Config {
	return &Config{
		Dialect:      SQLite,
		Output:       OutputTypeSQL,
		CreateTables: false,
	}
}
