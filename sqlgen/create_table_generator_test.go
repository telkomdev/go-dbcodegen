package sqlgen_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/wartek-id/core/tools/dbgen/config"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/dialect"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/sb"
	"gitlab.com/wartek-id/core/tools/dbgen/types/field_option"
)

func TestCreateTableGenerator_Dialect(t *testing.T) {
	dial := "postgres"
	do := dialect.DefaultDialectOption()

	sqlGen := sqlgen.NewCreateTableGenerator(dial, do)
	assert.Equal(t, dial, sqlGen.Dialect())
}

func TestCreateTableGenerator_DialectOptions(t *testing.T) {
	dial := "postgres"
	do := dialect.DefaultDialectOption()

	sqlGen := sqlgen.NewCreateTableGenerator(dial, do)
	assert.Equal(t, do, sqlGen.DialectOptions())
}

func TestCreateTableGenerator_ExpressionSQLGenerator(t *testing.T) {
	dial := "postgres"
	do := dialect.DefaultDialectOption()

	sqlGen := sqlgen.NewCreateTableGenerator(dial, do)
	assert.NotNil(t, sqlGen.ExpressionSQLGenerator())
}

func TestCreateTableGenerator_Generate(t *testing.T) {
	testCases := []struct {
		dialect *dialect.DialectOption
		input   *config.Schema
		result  string
	}{
		{
			dialect: dialect.DefaultDialectOption(),
			input: &config.Schema{
				Name: "user",
				Fields: []*config.Field{
					{
						Name: "id",
						Type: "bigserial",
						Options: []field_option.FieldOption{
							field_option.NotNull,
						},
					},
					{
						Name:  "name",
						Type:  "varchar",
						Limit: 255,
						Options: []field_option.FieldOption{
							field_option.NotNull,
						},
					},
					{
						Name:  "school",
						Type:  "varchar",
						Limit: 100,
						Options: []field_option.FieldOption{
							field_option.Nullable,
						},
					},
					{
						Name:  "salary",
						Type:  "decimal",
						Limit: 5,
						Scale: 2,
					},
				},
			},
			result: "CREATE TABLE IF NOT EXISTS \"user\" (\n\t\"id\" BIGSERIAL NOT NULL,\n\t\"name\" VARCHAR(255) NOT NULL,\n\t\"school\" VARCHAR(100) NULL,\n\t\"salary\" DECIMAL(5, 2)\n);",
		},
	}

	for _, tc := range testCases {
		buf := sb.NewSQLBuilder()
		sqlGen := sqlgen.NewCreateTableGenerator("postgres", tc.dialect)
		sqlGen.Generate(buf, tc.input)
		result, err := buf.ToSQL()
		assert.Nil(t, err)
		assert.Equal(t, tc.result, result)
	}
}
