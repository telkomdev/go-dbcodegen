package sqlgen_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/wartek-id/core/tools/dbgen/config"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/dialect"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/sb"
)

func TestCreateIndexGenerator_Dialect(t *testing.T) {
	dial := "postgres"
	do := dialect.DefaultDialectOption()

	sqlGen := sqlgen.NewCreateIndexGenerator(dial, do)
	assert.Equal(t, dial, sqlGen.Dialect())
}

func TestCreateIndexGenerator_DialectOptions(t *testing.T) {
	dial := "postgres"
	do := dialect.DefaultDialectOption()

	sqlGen := sqlgen.NewCreateIndexGenerator(dial, do)
	assert.Equal(t, do, sqlGen.DialectOptions())
}

func TestCreateIndexGenerator_ExpressionSQLGenerator(t *testing.T) {
	dial := "postgres"
	do := dialect.DefaultDialectOption()

	sqlGen := sqlgen.NewCreateIndexGenerator(dial, do)
	assert.NotNil(t, sqlGen.ExpressionSQLGenerator())
}

func TestCreateIndexGenerator_Generate(t *testing.T) {
	doConcurrent := dialect.DefaultDialectOption()
	doConcurrent.SupportConcurrently = true

	testCases := []struct {
		dialect *dialect.DialectOption
		input   *config.Index
		result  string
	}{
		{
			dialect: dialect.DefaultDialectOption(),
			input: &config.Index{
				Name: "idx_name",
				Fields: []*config.IndexField{
					{
						Column: "name",
						Order:  "ASC",
					},
				},
				Unique: true,
			},
			result: `CREATE UNIQUE INDEX IF NOT EXISTS "idx_name" ON "user"("name" ASC);`,
		},
		{
			dialect: dialect.DefaultDialectOption(),
			input: &config.Index{
				Name: "idx_name",
				Fields: []*config.IndexField{
					{
						Column: "name",
						Order:  "ASC",
					},
				},
			},
			result: `CREATE INDEX IF NOT EXISTS "idx_name" ON "user"("name" ASC);`,
		},
		{
			dialect: dialect.DefaultDialectOption(),
			input: &config.Index{
				Name: "idx_name",
				Fields: []*config.IndexField{
					{
						Column: "name",
						Order:  "ASC",
					},
					{
						Column: "age",
						Order:  "ASC",
					},
				},
			},
			result: `CREATE INDEX IF NOT EXISTS "idx_name" ON "user"("name" ASC, "age" ASC);`,
		},
		{
			dialect: doConcurrent,
			input: &config.Index{
				Name: "idx_name",
				Fields: []*config.IndexField{
					{
						Column: "name",
						Order:  "ASC",
					},
					{
						Column: "age",
					},
				},
			},
			result: `CREATE INDEX CONCURRENTLY IF NOT EXISTS "idx_name" ON "user"("name" ASC, "age");`,
		},
	}

	for _, tc := range testCases {
		buf := sb.NewSQLBuilder()
		sqlGen := sqlgen.NewCreateIndexGenerator("postgres", tc.dialect)
		sqlGen.Generate(buf, "user", tc.input)
		result, err := buf.ToSQL()
		assert.Nil(t, err)
		assert.Equal(t, tc.result, result)
	}
}
