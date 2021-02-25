package migrate

import (
	"fmt"
	"github.com/borodyadka/db-migrate/drivers"
	"strings"
)

type Table struct {
	schema, name string
	dialect      Dialect
}

func (t *Table) Schema() string {
	if t.schema == "" {
		return `"public"`
	}
	return fmt.Sprintf(`"%s"`, t.schema)
}

func (t *Table) String() string {
	if t.schema != "" {
		return fmt.Sprintf(`"%s"."%s"`, t.schema, t.name)
	}
	return t.name
}

func ParseTableName(t string, dialect Dialect) (*Table, error) {
	var (
		table  string
		schema string
	)
	switch parts := strings.Split(t, "."); len(parts) {
	case 1:
		table = parts[0]
	case 2:
		schema = parts[0]
		table = parts[1]
	default:
		return nil, drivers.NewErrInvalidTableName(t, 0) // TODO
	}
	return &Table{
		schema:  schema,
		name:    table,
		dialect: dialect,
	}, nil
}
