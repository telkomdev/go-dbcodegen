package exp_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/wartek-id/core/tools/dbgen/config"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/dialect"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/exp"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/sb"
	"gitlab.com/wartek-id/core/tools/dbgen/types/field_option"
)

func TestGetTypeFragment(t *testing.T) {
	ex := exp.NewExpressionSQLGenerator("", dialect.DefaultDialectOption())
	testCases := []struct {
		input  config.Field
		result string
	}{
		{
			input: config.Field{
				Type:  "decimal",
				Limit: 50,
				Scale: 2,
			},
			result: "DECIMAL(50, 2)",
		},
		{
			input: config.Field{
				Type:  "varchar",
				Limit: 255,
			},
			result: "VARCHAR(255)",
		},
		{
			input: config.Field{
				Type: "int",
			},
			result: "INT",
		},
	}

	for _, tc := range testCases {
		result := ex.GetTypeFragment(&tc.input)
		assert.Equal(t, tc.result, string(result))
	}
}

func TestGetOptionsFragment(t *testing.T) {
	ex := exp.NewExpressionSQLGenerator("", dialect.DefaultDialectOption())
	testCases := []struct {
		input  config.Field
		result string
	}{
		{
			input: config.Field{
				Options: []field_option.FieldOption{
					"not null",
					"auto increment",
				},
			},
			result: " NOT NULL AUTO INCREMENT",
		},
		{
			input: config.Field{
				Options: []field_option.FieldOption{
					"nullable",
				},
			},
			result: " NULL",
		},
		{
			input:  config.Field{},
			result: "",
		},
	}

	for _, tc := range testCases {
		result := ex.GetOptionsFragment(&tc.input)
		assert.Equal(t, tc.result, string(result))
	}
}

func TestLiteralExpression(t *testing.T) {
	ex := exp.NewExpressionSQLGenerator("", dialect.DefaultDialectOption())
	b := sb.NewSQLBuilder()

	ex.LiteralExpression(b, "user")
	assert.Equal(t, []byte("\"user\""), b.Bytes())
}

func TestGetDefaultValue(t *testing.T) {
	ex := exp.NewExpressionSQLGenerator("", dialect.DefaultDialectOption())

	result := ex.GetDefaultValue("user")
	assert.Equal(t, []byte("'user'"), result)

	result = ex.GetDefaultValue(1)
	assert.Equal(t, []byte("1"), result)
}
