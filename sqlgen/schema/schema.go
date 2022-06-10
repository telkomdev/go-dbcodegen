package schema

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.com/wartek-id/core/tools/dbgen/config"
)

var (
	ErrConstraintNotExists   = errors.New("constraint is not exists")
	ErrConstraintHasNoFields = errors.New("constraint has no fields")
	ErrUnsupportedDriver     = errors.New("unsupported driver")
)

type Schema interface {
	GetSchemas() ([]*config.Schema, error)
	GetTables() ([]string, error)
	GetIndices() (map[string]*Indices, error)
	GetFields(tblName string) ([]*config.Field, error)
	GetPrimaryKeys() (map[string]*PrimaryKey, error)
}

func NewSchema(connString string) (Schema, error) {
	driver := strings.Split(connString, "://")
	switch driver[0] {
	case "postgresql":
		pool, err := pgxpool.Connect(context.Background(), connString)
		if err != nil {
			return nil, err
		}

		return NewPostgresSchema(pool), nil
	}

	return nil, ErrUnsupportedDriver
}

type PrimaryKey struct {
	Table  string
	Name   string
	Column string
}

type Indices struct {
	Table   string
	Indices map[string]*config.Index
}

type TableStructure struct {
	ColumnName    string
	ColumnDefault sql.NullString
	IsNullable    string
	DataType      string
	CharMaxLen    sql.NullInt32
	NumPrecision  sql.NullInt32
	NumScale      sql.NullInt32
}

func (i *Indices) GetByConstraintName(name string) (*config.Index, error) {
	if i.Indices[name] == nil {
		return nil, ErrConstraintNotExists
	}

	idx := i.Indices[name]
	if len(idx.Fields) == 0 {
		return nil, ErrConstraintHasNoFields
	}

	return i.Indices[name], nil
}
