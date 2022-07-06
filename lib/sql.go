package lib

import (
	"github.com/thediveo/enumflag"
)

type SQLDialect enumflag.Flag

const (
	SQLite SQLDialect = iota
	PostGreSQL
	MySQL
)

var CurrentSQLDialect SQLDialect
