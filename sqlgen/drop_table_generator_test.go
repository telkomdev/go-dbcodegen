package sqlgen_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/wartek-id/core/tools/dbgen/config"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/dialect"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/sb"
)

func TestDropTableGenerator_Dialect(t *testing.T) {
	dial := "postgres"
	do := dialect.DefaultDialectOption()

	sqlGen := sqlgen.NewDropTableGenerator(dial, do)
	assert.Equal(t, dial, sqlGen.Dialect())
}

func TestDropTableGenerator_DialectOptions(t *testing.T) {
	dial := "postgres"
	do := dialect.DefaultDialectOption()

	sqlGen := sqlgen.NewDropTableGenerator(dial, do)
	assert.Equal(t, do, sqlGen.DialectOptions())
}

func TestDropTableGenerator_ExpressionSQLGenerator(t *testing.T) {
	dial := "postgres"
	do := dialect.DefaultDialectOption()

	sqlGen := sqlgen.NewDropTableGenerator(dial, do)
	assert.NotNil(t, sqlGen.ExpressionSQLGenerator())
}

func TestDropTableGenerator_Generate(t *testing.T) {
	testCases := []struct {
		dialect *dialect.DialectOption
		input   *config.Schema
		result  string
	}{
		{
			dialect: dialect.DefaultDialectOption(),
			input: &config.Schema{
				Name: "user",
			},
			result: `DROP TABLE IF EXISTS "user";`,
		},
		{
			dialect: dialect.DefaultDialectOption(),
			input: &config.Schema{
				Name: "school",
			},
			result: `DROP TABLE IF EXISTS "school";`,
		},
	}

	for _, tc := range testCases {
		buf := sb.NewSQLBuilder()
		sqlGen := sqlgen.NewDropTableGenerator("postgres", tc.dialect)
		sqlGen.Generate(buf, tc.input)
		result, err := buf.ToSQL()
		assert.Nil(t, err)
		assert.Equal(t, tc.result, result)
	}
}
