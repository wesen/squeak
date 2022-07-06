package lib

type Config struct {
	Dialect      SQLDialect
	Output       string
	CreateTables bool
}

func NewConfig() *Config {
	return &Config{
		Dialect:      SQLite,
		Output:       "sql",
		CreateTables: false,
	}
}
