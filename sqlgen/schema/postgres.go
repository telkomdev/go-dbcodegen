package schema

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v4"
	pg_query "github.com/pganalyze/pg_query_go/v2"
	"gitlab.com/wartek-id/core/tools/dbgen/config"
	"gitlab.com/wartek-id/core/tools/dbgen/types/field_option"
	"gitlab.com/wartek-id/core/tools/dbgen/types/field_type"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
)

const (
	SchemaMigrationTable = "schema_migrations"
)

type PgInterface interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

var OrderMapper = map[pg_query.SortByDir]string{
	pg_query.SortByDir_SORTBY_ASC:     "ASC",
	pg_query.SortByDir_SORTBY_DEFAULT: "ASC",
	pg_query.SortByDir_SORTBY_DESC:    "DESC",
}

var FieldTypeMapper = map[string]field_type.FieldType{
	"timestamp without time zone": field_type.Timestamp,
	"timestamp with time zone":    field_type.Timestamptz,
	"character varying":           field_type.Varchar,
	"numeric":                     field_type.Decimal,
	"double precision":            field_type.Float,
	"integer":                     field_type.Int,
	"boolean":                     field_type.Boolean,
}

const (
	RegexAutoIncrement = `nextval\(\'[^']+'::regclass\)`
	DefaultSchema      = "public"
)

type postgresSchema struct {
	pool   PgInterface
	schema string

	indicesLoaded bool
	indices       map[string]*Indices

	primaryKeysLoaded bool
	primaryKeys       map[string]*PrimaryKey
}

func NewPostgresSchema(pool PgInterface) *postgresSchema {
	return &postgresSchema{
		pool:   pool,
		schema: DefaultSchema,
	}
}

func (s *postgresSchema) GetSchemas() ([]*config.Schema, error) {
	schemas := make([]*config.Schema, 0)

	tables, err := s.GetTables()
	if err != nil {
		return nil, err
	}

	for _, table := range tables {
		fields, err := s.GetFields(table)
		if err != nil {
			return nil, err
		}

		indices, err := s.GetTablesIndices(table)
		if err != nil {
			return nil, err
		}
		schema := &config.Schema{
			Name:   table,
			Fields: fields,
			Index:  indices,
		}

		schemas = append(schemas, schema)
	}

	return schemas, nil
}

func (s *postgresSchema) GetTables() ([]string, error) {
	query, _, err := goqu.Dialect("postgres").From("information_schema.tables").
		Where(
			goqu.C("table_schema").Eq(s.schema),
			goqu.C("table_type").Eq("BASE TABLE"),
		).Select("table_name").ToSQL()
	if err != nil {
		return nil, err
	}

	tables := make([]string, 0)
	rows, err := s.pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		if err != nil {
			return nil, err
		}

		if table == SchemaMigrationTable {
			continue
		}

		tables = append(tables, table)
	}

	return tables, nil
}

func (s *postgresSchema) GetFields(name string) ([]*config.Field, error) {
	query, _, err := goqu.Dialect("postgres").From("information_schema.columns").
		Where(
			goqu.C("table_schema").Eq(s.schema),
			goqu.C("table_name").Eq(name),
		).Select(
		"column_name", "column_default", "is_nullable",
		"data_type", "character_maximum_length", "numeric_precision",
		"numeric_scale").ToSQL()
	if err != nil {
		return nil, err
	}
	rows, err := s.pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	fields := make([]*config.Field, 0)
	for rows.Next() {
		table := TableStructure{}
		err := rows.Scan(&table.ColumnName, &table.ColumnDefault, &table.IsNullable,
			&table.DataType, &table.CharMaxLen, &table.NumPrecision, &table.NumScale)
		if err != nil {
			return nil, err
		}

		fields = append(fields, s.getField(name, &table))
	}
	return fields, nil
}

func (s *postgresSchema) getField(name string, table *TableStructure) *config.Field {
	ft := FieldTypeMapper[table.DataType]
	if ft == "" {
		ft = field_type.ParseString(table.DataType)
	}

	if s.isAutoIncrement(table.ColumnDefault.String) {
		switch ft {
		case field_type.BigInt:
			ft = field_type.BigSerial
		case field_type.Int:
			ft = field_type.Serial
		case field_type.SmallInt:
			ft = field_type.SmallSerial
		}
	}

	field := &config.Field{
		Name:    table.ColumnName,
		Type:    ft,
		Default: s.ParseDefaultValue(table.ColumnDefault.String),
		Options: s.GetOptions(name, table),
	}

	switch ft.Type() {
	case field_type.FieldTypeString:
		if table.CharMaxLen.Valid {
			field.Limit = int(table.CharMaxLen.Int32)
		}
	case field_type.FieldTypeNumeric:
		if table.NumPrecision.Valid {
			field.Limit = int(table.NumPrecision.Int32)
		}
		if table.NumScale.Valid {
			field.Scale = int(table.NumScale.Int32)
		}
	}

	return field
}

func (s *postgresSchema) isAutoIncrement(defaultValue string) bool {
	match, err := regexp.MatchString(RegexAutoIncrement, defaultValue)
	if err != nil {
		return false
	}
	return match
}

func (s *postgresSchema) ParseDefaultValue(value string) interface{} {
	if value == "" || s.isAutoIncrement(value) {
		return nil
	}

	valSplit := strings.Split(value, "::")
	if len(valSplit) > 1 {
		valValue, valType := valSplit[0], valSplit[1]

		switch valType {
		case "character varying", "text", "timestamp without time zone", "timestamp with time zone":
			// trim ' prefix and suffix
			trimVal := strings.TrimPrefix(strings.TrimSuffix(valValue, "'"), "'")
			// replace escaped '' with '
			return strings.Replace(trimVal, "''", "'", -1)
		}
		return valValue
	}

	valInt, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		return int(valInt)
	}

	valFloat, err := strconv.ParseFloat(value, 64)
	if err == nil {
		return valFloat
	}
	return nil
}

func (s *postgresSchema) GetOptions(name string, table *TableStructure) []field_option.FieldOption {
	options := make([]field_option.FieldOption, 0)
	pk, _ := s.GetPrimaryKey(name)
	if pk != nil && pk.Column == table.ColumnName {
		options = append(options, field_option.PrimaryKey)
	}

	if table.IsNullable == "NO" {
		options = append(options, field_option.NotNull)
	}

	return options
}

func (s *postgresSchema) GetPrimaryKeys() (map[string]*PrimaryKey, error) {
	err := s.LoadPrimaryKeys()
	if err != nil {
		return nil, err
	}
	return s.primaryKeys, nil
}

func (s *postgresSchema) GetPrimaryKey(table string) (*PrimaryKey, error) {
	primaryKeys, err := s.GetPrimaryKeys()
	if err != nil {
		return nil, err
	}
	return primaryKeys[table], nil
}

func (s *postgresSchema) LoadPrimaryKeys() error {
	if s.primaryKeysLoaded {
		return nil
	}

	indices, err := s.GetIndices()
	if err != nil {
		return err
	}

	pkConstraint, err := s.GetPKConstraint()
	if err != nil {
		return err
	}

	primaryKeys := make(map[string]*PrimaryKey)
	for tablename, constraintname := range pkConstraint {
		tableIndices := indices[tablename]
		if indices == nil {
			continue
		}

		constraint, err := tableIndices.GetByConstraintName(constraintname)
		if err != nil {
			continue
		}

		primaryKeys[tablename] = &PrimaryKey{
			Table:  tablename,
			Name:   constraintname,
			Column: constraint.GetColumns()[0],
		}
	}

	s.primaryKeys = primaryKeys
	s.primaryKeysLoaded = true
	return nil
}

func (s *postgresSchema) GetPKConstraint() (map[string]string, error) {
	query, _, err := goqu.Dialect("postgres").
		From("information_schema.table_constraints").
		Where(
			goqu.C("constraint_schema").Eq(s.schema),
			goqu.C("constraint_type").Eq("PRIMARY KEY"),
		).Select("table_name", "constraint_name").ToSQL()

	if err != nil {
		return nil, err
	}

	constraints := make(map[string]string)
	rows, err := s.pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var tablename, constraint string
		err := rows.Scan(&tablename, &constraint)
		if err != nil {
			return nil, err
		}

		constraints[tablename] = constraint
	}

	return constraints, nil
}

func (s *postgresSchema) GetTablesIndices(name string) ([]*config.Index, error) {
	indices, err := s.GetIndices()
	if err != nil {
		return nil, err
	}

	tablePk, err := s.GetPrimaryKey(name)
	if err != nil {
		return nil, err
	}

	result := make([]*config.Index, 0)
	tableIndices := indices[name]
	if tableIndices == nil {
		return result, nil
	}

	for name, index := range tableIndices.Indices {
		if tablePk == nil || name != tablePk.Name {
			result = append(result, index)
		}
	}

	return result, nil
}

func (s *postgresSchema) GetIndices() (map[string]*Indices, error) {
	err := s.LoadIndices()
	if err != nil {
		return nil, err
	}
	return s.indices, nil
}

func (s *postgresSchema) LoadIndices() error {
	if s.indicesLoaded {
		return nil
	}

	query, _, err := goqu.Dialect("postgres").
		From("pg_catalog.pg_indexes").
		Where(goqu.C("schemaname").Eq(s.schema)).
		Select("tablename", "indexname", "indexdef").
		ToSQL()
	if err != nil {
		return err
	}
	rows, err := s.pool.Query(context.Background(), query)
	if err != nil {
		return err
	}

	indices := make(map[string]*Indices)
	for rows.Next() {
		var tablename, indexname, indexdef string
		err := rows.Scan(&tablename, &indexname, &indexdef)
		if err != nil {
			return err
		}

		idxdef, err := pg_query.Parse(indexdef)
		if err != nil {
			return err
		}

		if len(idxdef.Stmts) == 0 {
			return errors.New("invalid statement")
		}

		if indices[tablename] == nil {
			indices[tablename] = &Indices{
				Table:   tablename,
				Indices: make(map[string]*config.Index),
			}
		}

		container := indices[tablename].Indices
		idxStmt := idxdef.Stmts[0].GetStmt().GetIndexStmt()
		index := config.Index{
			Name:   idxStmt.GetIdxname(),
			Fields: []*config.IndexField{},
			Unique: idxStmt.Unique,
		}

		for _, field := range idxStmt.GetIndexParams() {
			ordering := field.GetIndexElem().Ordering
			index.Fields = append(index.Fields, &config.IndexField{
				Column: field.GetIndexElem().GetName(),
				Order:  OrderMapper[ordering],
			})
		}

		container[indexname] = &index
		indices[tablename].Indices = container
	}

	s.indices = indices
	s.indicesLoaded = true

	return nil
}
