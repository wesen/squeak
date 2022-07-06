package lib

import (
	"github.com/thediveo/enumflag"
)

var SQLDialectIds = map[SQLDialect][]string{
	SQLite:     {"sqlite"},
	MySQL:      {"mysql"},
	PostGreSQL: {"postgresql"},
}

type SQLDialect enumflag.Flag

const (
	SQLite SQLDialect = iota
	PostGreSQL
	MySQL
)
