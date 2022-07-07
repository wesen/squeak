package lib

import (
	"bytes"
	"github.com/bxcodec/faker/v3"
	"github.com/feiin/sqlstring"
	"github.com/huandu/go-sqlbuilder"
	"github.com/pkg/errors"
	"github.com/thediveo/enumflag"
	"log"
	"math/rand"
	"strings"
)

type FieldType enumflag.Flag

const (
	FieldTypeText FieldType = iota
	FieldTypeInt
	FieldTypeAuto
	FieldTypeFirstName
	FieldTypeLastName
	FieldTypeEmail
)

var fieldTypeToDefaultSQLType = map[FieldType]string{
	FieldTypeText:      "TEXT",
	FieldTypeInt:       "INTEGER",
	FieldTypeAuto:      "INTEGER",
	FieldTypeFirstName: "TEXT",
	FieldTypeLastName:  "TEXT",
	FieldTypeEmail:     "TEXT",
}

var toFieldType = map[string]FieldType{
	"text":      FieldTypeText,
	"int":       FieldTypeInt,
	"auto":      FieldTypeAuto,
	"firstName": FieldTypeFirstName,
	"lastName":  FieldTypeLastName,
	"email":     FieldTypeEmail,
}

var dialectToFlavor = map[SQLDialect]sqlbuilder.Flavor{
	SQLite:     sqlbuilder.SQLite,
	PostGreSQL: sqlbuilder.PostgreSQL,
	MySQL:      sqlbuilder.MySQL,
}

func CreateTables(
	tableDefinitions map[string]map[string]string,
	dialect SQLDialect,
	output OutputType,
	generateTableOptions map[string]*TableGenerateConfig) (map[string]string, error) {
	res := make(map[string]string)

	for tableName, fieldDefinitions := range tableDefinitions {
		if _, ok := generateTableOptions[tableName]; ok {
			switch output {
			case OutputTypeSQL:
				cb := sqlbuilder.NewCreateTableBuilder()
				cb.SetFlavor(dialectToFlavor[dialect])
				cb.CreateTable(tableName).IfNotExists()
				for field, definition := range fieldDefinitions {
					fieldType, ok := toFieldType[definition]
					if !ok {
						log.Fatalf("Unknown field type %s", definition)
					}
					cb.Define(field, fieldTypeToDefaultSQLType[fieldType])
				}
				res[tableName] = cb.String()
			case OutputTypeSQLite:
				log.Fatal("SQLite format not implemented")
			case OutputTypeCSV:
				log.Fatal("Can't output create tables for CSV")
			}
		}
	}

	return res, nil
}

func GenerateData(
	tableDefinitions map[string]map[string]string,
	dialect SQLDialect,
	output OutputType,
	generateTableOptions map[string]*TableGenerateConfig) (map[string]string, error) {

	res := make(map[string]string)
	for tableName, fieldDefinitions := range tableDefinitions {
		if _opt, ok := generateTableOptions[tableName]; ok {
			data, err := generateTableData(fieldDefinitions, _opt)
			if err != nil {
				return nil, err
			}
			switch output {
			case OutputTypeSQLite:
				log.Fatal("SQLite format not implemented")
			case OutputTypeSQL:
				res[tableName] = toSQL(tableName, fieldDefinitions, data)
			case OutputTypeCSV:
				log.Fatalf("CSV output not implemented")
			}
		}
	}
	return res, nil
}

func generateTableData(
	fields map[string]string, options *TableGenerateConfig,
) (map[string][]interface{}, error) {
	ret := make(map[string][]interface{})
	for k, v := range fields {
		fieldType, ok := toFieldType[v]
		if !ok {
			return nil, errors.Errorf("Unknown field type %s", v)
		}
		ret[k] = make([]interface{}, options.Count)
		for i := 0; i < options.Count; i++ {
			var v interface{} = nil
			switch fieldType {
			case FieldTypeText:
				v = faker.Sentence()
			case FieldTypeInt:
				v = rand.Int()
			case FieldTypeAuto:
				v = i + 1
			case FieldTypeFirstName:
				v = faker.FirstName()
			case FieldTypeLastName:
				v = faker.LastName()
			case FieldTypeEmail:
				v = faker.Email()
			}
			ret[k][i] = v
		}
	}

	return ret, nil
}

func toSQL(tableName string, fieldDefinitions map[string]string, data map[string][]interface{}) string {
	if len(data) == 0 || len(fieldDefinitions) == 0 {
		return ""
	}

	buf := &bytes.Buffer{}
	buf.WriteString("INSERT INTO ")
	buf.WriteString(tableName)
	buf.WriteString("\n( ")

	var fields []string
	for k := range fieldDefinitions {
		fields = append(fields, k)
	}
	buf.WriteString(strings.Join(fields, "\n, "))
	buf.WriteString(")\nVALUES\n")

	// transpose the data matrix
	for i := 0; i < len(data[fields[0]]); i++ {
		var row []string
		for j := 0; j < len(fields); j++ {
			row = append(row, sqlstring.Escape(data[fields[j]][i]))
		}
		if i > 0 {
			buf.WriteString(", ")
		} else {
			buf.WriteString("  ")
		}
		buf.WriteString("(" + strings.Join(row, ", ") + ")\n")
	}
	buf.WriteString(";\n")

	return buf.String()
}
