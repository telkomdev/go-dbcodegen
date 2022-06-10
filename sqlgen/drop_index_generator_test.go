package sqlgen_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/wartek-id/core/tools/dbgen/config"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/dialect"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/sb"
)

func TestDropIndexGenerator_Dialect(t *testing.T) {
	dial := "postgres"
	do := dialect.DefaultDialectOption()

	sqlGen := sqlgen.NewDropIndexGenerator(dial, do)
	assert.Equal(t, dial, sqlGen.Dialect())
}

func TestDropIndexGenerator_DialectOptions(t *testing.T) {
	dial := "postgres"
	do := dialect.DefaultDialectOption()

	sqlGen := sqlgen.NewDropIndexGenerator(dial, do)
	assert.Equal(t, do, sqlGen.DialectOptions())
}

func TestDropIndexGenerator_ExpressionSQLGenerator(t *testing.T) {
	dial := "postgres"
	do := dialect.DefaultDialectOption()

	sqlGen := sqlgen.NewDropIndexGenerator(dial, do)
	assert.NotNil(t, sqlGen.ExpressionSQLGenerator())
}

func TestDropIndexGenerator_Generate(t *testing.T) {
	testCases := []struct {
		dialect *dialect.DialectOption
		input   *config.Index
		result  string
	}{
		{
			dialect: dialect.DefaultDialectOption(),
			input: &config.Index{
				Name: "index_on_name",
			},
			result: `DROP INDEX IF EXISTS "index_on_name";`,
		},
		{
			dialect: dialect.DefaultDialectOption(),
			input: &config.Index{
				Name: "index_on_school",
			},
			result: `DROP INDEX IF EXISTS "index_on_school";`,
		},
	}

	for _, tc := range testCases {
		buf := sb.NewSQLBuilder()
		sqlGen := sqlgen.NewDropIndexGenerator("postgres", tc.dialect)
		sqlGen.Generate(buf, tc.input)
		result, err := buf.ToSQL()
		assert.Nil(t, err)
		assert.Equal(t, tc.result, result)
	}
}
