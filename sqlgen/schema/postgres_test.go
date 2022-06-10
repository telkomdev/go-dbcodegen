package schema_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/wartek-id/core/tools/dbgen/config"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/schema"
	"gitlab.com/wartek-id/core/tools/dbgen/types/field_option"
)

func TestPostgres_GetSchemas(t *testing.T) {
	testCases := map[string]struct {
		tableResult *pgxmock.Rows
		tableErr    error
		fieldResult *pgxmock.Rows
		fieldErr    error
		indexResult *pgxmock.Rows
		indexErr    error
		constResult *pgxmock.Rows
		result      []*config.Schema
		err         error
	}{
		"success": {
			tableResult: pgxmock.NewRows([]string{
				"table_name",
			}).AddRow("example"),
			fieldResult: pgxmock.NewRows([]string{
				"column_name", "column_default", "is_nullable", "data_type", "character_maximum_length",
				"numeric_precision", "numeric_scale",
			}).AddRow(
				"id", "nextval('some_id_sec'::regclass)", "NO", "bigint", nil, 64, nil,
			).AddRow(
				"name", "'Alfred'::character varying", "NO", "character varying", "200", nil, nil,
			).AddRow(
				"price", "100.5", "YES", "numeric", nil, 64, 2,
			),
			indexResult: pgxmock.NewRows([]string{
				"tablename", "indexname", "indexdef",
			}).AddRow(
				"example", "example_pkey", "CREATE UNIQUE INDEX example_pkey ON public.example USING btree (id)",
			),
			constResult: pgxmock.NewRows([]string{
				"table_name", "constraint_name",
			}).AddRow("example", "example_pkey"),
			result: []*config.Schema{
				{
					Name: "example",
					Fields: []*config.Field{
						{
							Name:  "id",
							Type:  "bigserial",
							Limit: 64,
							Options: []field_option.FieldOption{
								field_option.PrimaryKey,
								field_option.NotNull,
							},
						},
						{
							Name:    "name",
							Type:    "varchar",
							Limit:   200,
							Default: "Alfred",
							Options: []field_option.FieldOption{field_option.NotNull},
						},
						{
							Name:    "price",
							Type:    "decimal",
							Scale:   2,
							Limit:   64,
							Default: 100.5,
							Options: []field_option.FieldOption{},
						},
					},
					Index: []*config.Index{},
				},
			},
		},
		"error get table": {
			tableErr: errors.New("error get table"),
			err:      errors.New("error get table"),
		},
		"error get field": {
			tableResult: pgxmock.NewRows([]string{
				"table_name",
			}).AddRow("example"),
			fieldErr: errors.New("error get field"),
			err:      errors.New("error get field"),
		},
		"error get index": {
			tableResult: pgxmock.NewRows([]string{
				"table_name",
			}).AddRow("example"),
			fieldResult: pgxmock.NewRows([]string{
				"column_name", "column_default", "is_nullable", "data_type", "character_maximum_length",
				"numeric_precision", "numeric_scale",
			}),
			indexErr: errors.New("error get index"),
			err:      errors.New("error get index"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mock, err := pgxmock.NewConn()
			assert.Nil(t, err)
			defer mock.Close(context.Background())

			if tc.tableResult != nil {
				mock.ExpectQuery("SELECT [^(FROM)]+FROM \"information_schema\".\"tables\"").
					WillReturnRows(tc.tableResult)
			}
			if tc.tableErr != nil {
				mock.ExpectQuery("SELECT [^(FROM)]+FROM \"information_schema\".\"tables\"").
					WillReturnError(tc.tableErr)
			}
			if tc.fieldResult != nil {
				mock.ExpectQuery("SELECT [^(FROM)]+FROM \"information_schema\".\"columns\"").
					WillReturnRows(tc.fieldResult)
			}
			if tc.fieldErr != nil {
				mock.ExpectQuery("SELECT [^(FROM)]+FROM \"information_schema\".\"columns\"").
					WillReturnError(tc.fieldErr)
			}
			if tc.indexResult != nil {
				mock.ExpectQuery("SELECT [^(FROM)]+FROM \"pg_catalog\".\"pg_indexes\"").
					WillReturnRows(tc.indexResult)
			}
			if tc.indexErr != nil {
				mock.ExpectQuery("SELECT [^(FROM)]+FROM \"pg_catalog\".\"pg_indexes\"").
					WillReturnError(tc.indexErr)
			}
			if tc.constResult != nil {
				mock.ExpectQuery("SELECT [^(FROM)]+FROM \"information_schema\".\"table_constraints\"").
					WillReturnRows(tc.constResult)
			}

			sc := schema.NewPostgresSchema(mock)
			result, err := sc.GetSchemas()
			assert.Equal(t, tc.err, err)
			assert.ElementsMatch(t, tc.result, result)
		})
	}
}

func TestPostgres_GetField(t *testing.T) {
	mock, err := pgxmock.NewConn()
	assert.Nil(t, err)
	defer mock.Close(context.Background())

	fieldsResults := pgxmock.NewRows([]string{
		"column_name", "column_default", "is_nullable", "data_type", "character_maximum_length",
		"numeric_precision", "numeric_scale",
	}).AddRow(
		"id", "nextval('some_id_sec'::regclass)", "NO", "bigint", nil, 64, nil,
	).AddRow(
		"name", "'Alfred'::character varying", "NO", "character varying", "200", nil, nil,
	).AddRow(
		"price", "100.5", "YES", "numeric", nil, 64, 2,
	)
	mock.ExpectQuery("SELECT [^(FROM)]+FROM \"information_schema\".\"columns\"").
		WillReturnRows(fieldsResults)

	indicesResults := pgxmock.NewRows([]string{
		"tablename", "indexname", "indexdef",
	}).AddRow(
		"example", "example_pkey", "CREATE UNIQUE INDEX example_pkey ON public.example USING btree (id)",
	)
	mock.ExpectQuery("SELECT [^(FROM)]+FROM \"pg_catalog\".\"pg_indexes\"").
		WillReturnRows(indicesResults)

	constraintResult := pgxmock.NewRows([]string{
		"table_name", "constraint_name",
	}).AddRow("example", "example_pkey")
	mock.ExpectQuery("SELECT [^(FROM)]+FROM \"information_schema\".\"table_constraints\"").
		WillReturnRows(constraintResult)

	sc := schema.NewPostgresSchema(mock)
	result, err := sc.GetFields("example")
	assert.Nil(t, err)
	assert.ElementsMatch(t, []*config.Field{
		{
			Name:  "id",
			Type:  "bigserial",
			Limit: 64,
			Options: []field_option.FieldOption{
				field_option.PrimaryKey,
				field_option.NotNull,
			},
		},
		{
			Name:    "name",
			Type:    "varchar",
			Limit:   200,
			Default: "Alfred",
			Options: []field_option.FieldOption{field_option.NotNull},
		},
		{
			Name:    "price",
			Type:    "decimal",
			Scale:   2,
			Limit:   64,
			Default: 100.5,
			Options: []field_option.FieldOption{},
		},
	}, result)
}

func TestPostgres_GetTables(t *testing.T) {
	mock, err := pgxmock.NewConn()
	assert.Nil(t, err)
	defer mock.Close(context.Background())

	mock.ExpectQuery("SELECT *").
		WillReturnRows(pgxmock.NewRows([]string{"table_name"}).
			AddRow("users").
			AddRow("histories").
			AddRow("schema_migrations"))
	sc := schema.NewPostgresSchema(mock)
	result, err := sc.GetTables()
	assert.Nil(t, err)
	assert.ElementsMatch(t, []string{"users", "histories"}, result)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgres_GetIndices(t *testing.T) {
	mock, err := pgxmock.NewConn()
	assert.Nil(t, err)
	defer mock.Close(context.Background())

	indicesResults := pgxmock.NewRows([]string{
		"tablename", "indexname", "indexdef",
	}).AddRow(
		"example", "example_pkey", "CREATE UNIQUE INDEX example_pkey ON public.example USING btree (id)",
	).AddRow(
		"example", "index_example_on_name_email", "CREATE INDEX index_example_on_name_email ON public.example(name ASC, email DESC)",
	)
	mock.ExpectQuery("SELECT [^(FROM)]+FROM \"pg_catalog\".\"pg_indexes\"").
		WillReturnRows(indicesResults)

	expectedResult := map[string]*schema.Indices{
		"example": {
			Table: "example",
			Indices: map[string]*config.Index{
				"example_pkey": {
					Name: "example_pkey",
					Fields: []*config.IndexField{
						{Column: "id", Order: "ASC"},
					},
					Unique: true,
				},
				"index_example_on_name_email": {
					Name: "index_example_on_name_email",
					Fields: []*config.IndexField{
						{Column: "name", Order: "ASC"},
						{Column: "email", Order: "DESC"},
					},
				},
			},
		},
	}
	sc := schema.NewPostgresSchema(mock)
	result, err := sc.GetIndices()
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestPostgres_GetPrimaryKeys(t *testing.T) {
	mock, err := pgxmock.NewConn()
	assert.Nil(t, err)
	defer mock.Close(context.Background())

	indicesResults := pgxmock.NewRows([]string{
		"tablename", "indexname", "indexdef",
	}).AddRow(
		"example", "example_pkey", "CREATE UNIQUE INDEX example_pkey ON public.example USING btree (id)",
	).AddRow(
		"example", "index_example_on_name_email", "CREATE INDEX index_example_on_name_email ON public.example(name ASC, email DESC)",
	)
	mock.ExpectQuery("SELECT [^(FROM)]+FROM \"pg_catalog\".\"pg_indexes\"").
		WillReturnRows(indicesResults)
	constraintResult := pgxmock.NewRows([]string{
		"table_name", "constraint_name",
	}).AddRow("example", "example_pkey")
	mock.ExpectQuery("SELECT [^(FROM)]+FROM \"information_schema\".\"table_constraints\"").
		WillReturnRows(constraintResult)

	expectedResult := map[string]*schema.PrimaryKey{
		"example": {
			Table:  "example",
			Name:   "example_pkey",
			Column: "id",
		},
	}
	sc := schema.NewPostgresSchema(mock)
	result, err := sc.GetPrimaryKeys()
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestPostgres_ParseDefaultValue(t *testing.T) {
	testCases := map[string]struct {
		input  string
		result interface{}
	}{
		"int": {
			input:  "100",
			result: 100,
		},
		"float int": {
			input:  "100.00",
			result: 100.00,
		},
		"float": {
			input:  "100.50",
			result: 100.50,
		},
		"pk": {
			input:  "nextval('reg_seq'::regclass)",
			result: nil,
		},
		"varchar": {
			input:  "'100.50'::character varying",
			result: "100.50",
		},
		"text": {
			input:  "'100.50'::text",
			result: "100.50",
		},
	}

	sc := schema.NewPostgresSchema(nil)
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := sc.ParseDefaultValue(tc.input)
			assert.Equal(t, tc.result, result)
		})
	}
}
