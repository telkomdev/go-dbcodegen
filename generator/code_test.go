package generator_test

import (
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/stretchr/testify/assert"
	"gitlab.com/wartek-id/core/tools/dbgen/config"
	"gitlab.com/wartek-id/core/tools/dbgen/generator"
)

func TestGenerateSelectQuery(t *testing.T) {
	dialect := goqu.Dialect("postgres")
	element := &config.Schema{
		Name: "user",
	}
	idx := goqu.Ex{"id": "1"}
	ds := dialect.From(element.Name).Where(idx)
	sql, _, _ := ds.Prepared(true).ToSQL()

	res, err := generator.GenerateSelectQuery(dialect, element, idx)
	assert.Nil(t, err)
	assert.Equal(t, res, sql)
}

func TestGenerateInsertQuery(t *testing.T) {
	dialect := goqu.Dialect("postgres")
	element := &config.Schema{
		Name: "user",
		Fields: []*config.Field{
			{
				Name: "name",
			}, {
				Name: "id",
			},
		},
	}
	tc := &config.Function{
		Name:      "CreateUser",
		Query:     "INSERT INTO \"user\" (\"name\") VALUES ($1)",
		TableName: "user",
		SqlcType:  ":exec",
	}
	res := generator.GenerateInsertQuery(dialect, element)
	assert.Equal(t, res, tc)
}

func TestGenerateUpdateQuery(t *testing.T) {
	dialect := goqu.Dialect("postgres")
	element := &config.Schema{
		Name: "user",
		Fields: []*config.Field{
			{
				Name: "name",
			},
			{
				Name: "deleted_at",
			},
			{
				Name: "id",
			},
		},
	}
	tc := &config.Function{
		Name:      "UpdateUser",
		Query:     "UPDATE \"user\" SET \"name\"=$1 WHERE ((\"id\" = $2) AND (\"deleted_at\" IS NULL))",
		TableName: "user",
		SqlcType:  ":exec",
	}
	res := generator.GenerateUpdateQuery(dialect, element)
	assert.Equal(t, res, tc)
}

func TestGenerateDestroyQuery(t *testing.T) {
	dialect := goqu.Dialect("postgres")
	element := &config.Schema{
		Name: "user",
		Fields: []*config.Field{
			{
				Name: "name",
			},
		},
	}
	tc := &config.Function{
		Name:      "DestroyUser",
		Query:     "DELETE FROM \"user\" WHERE (\"id\" = $1)",
		TableName: "user",
		SqlcType:  ":exec",
	}
	res, err := generator.GenerateDestroyQuery(dialect, element)
	assert.Nil(t, err)
	assert.Equal(t, res, tc)
}

func TestGenerateDeleteQuery(t *testing.T) {
	dialect := goqu.Dialect("postgres")
	element := &config.Schema{
		Name: "user",
		Fields: []*config.Field{
			{
				Name: "name",
			},
			{
				Name: "deleted_at",
			},
		},
	}
	tc := &config.Function{
		Name:      "DeleteUser",
		Query:     "UPDATE \"user\" SET \"deleted_at\"=NOW() WHERE ((\"id\" = $1) AND (\"deleted_at\" IS NULL))",
		TableName: "user",
		SqlcType:  ":exec",
	}
	res := generator.GenerateDeleteQuery(dialect, element)
	assert.Equal(t, res, tc)
}

func TestGenerateDeleteQuery_ReturnNil(t *testing.T) {
	dialect := goqu.Dialect("postgres")
	element := &config.Schema{
		Name: "user",
		Fields: []*config.Field{
			{
				Name: "name",
			},
		},
	}
	res := generator.GenerateDeleteQuery(dialect, element)
	assert.Nil(t, res)
}

func TestGenerateSelectQueryByIndex(t *testing.T) {
	dialect := goqu.Dialect("postgres")
	element := &config.Schema{
		Name: "user",
		Fields: []*config.Field{
			{
				Name: "name",
			},
		},
		Index: []*config.Index{
			{
				Name: "idx_name_and_email",
				Fields: []*config.IndexField{
					{
						Column: "name",
						Order:  "ASC",
					},
					{
						Column: "email",
						Order:  "ASC",
					},
				},
			},
		},
	}
	tc := []*config.Function{
		{
			Name:      "FindUserById",
			Query:     "SELECT * FROM \"user\" WHERE (\"id\" = $1)",
			TableName: "user",
			SqlcType:  ":one",
		},
		{
			Name:      "FindUserByNameAndEmail",
			Query:     "SELECT * FROM \"user\" WHERE ((\"email\" = $1) AND (\"name\" = $2))",
			TableName: "user",
			SqlcType:  ":many",
		},
	}
	res := generator.GenerateSelectQueryByIndex(dialect, element)
	assert.Equal(t, res, tc)
}

func TestGenerateFunctionName(t *testing.T) {
	tc := "FindById"
	res := generator.GenerateFunctionName("find", "_id")
	assert.Equal(t, res, tc)
}

func TestGenerateQueries(t *testing.T) {
	element := []*config.Schema{
		{
			Name: "user",
			Fields: []*config.Field{
				{
					Name: "name",
				},
				{
					Name: "deleted_at",
				},
			},
			Index: []*config.Index{
				{
					Name: "idx_name_and_email",
					Fields: []*config.IndexField{
						{
							Column: "name",
							Order:  "ASC",
						},
						{
							Column: "email",
							Order:  "ASC",
						},
					},
				},
			},
		},
	}
	tc := []*config.Function{
		{
			Name:      "FindUserById",
			Query:     "SELECT * FROM \"user\" WHERE (\"id\" = $1)",
			TableName: "user",
			SqlcType:  ":one",
		},
		{
			Name:      "FindUserByNameAndEmail",
			Query:     "SELECT * FROM \"user\" WHERE ((\"email\" = $1) AND (\"name\" = $2))",
			TableName: "user",
			SqlcType:  ":many",
		},
		{
			Name:      "UpdateUser",
			Query:     "UPDATE \"user\" SET \"name\"=$1 WHERE ((\"id\" = $2) AND (\"deleted_at\" IS NULL))",
			TableName: "user",
			SqlcType:  ":exec",
		},
		{
			Name:      "CreateUser",
			Query:     "INSERT INTO \"user\" (\"name\") VALUES ($1)",
			TableName: "user",
			SqlcType:  ":exec",
		},
		{
			Name:      "DeleteUser",
			Query:     "UPDATE \"user\" SET \"deleted_at\"=NOW() WHERE ((\"id\" = $1) AND (\"deleted_at\" IS NULL))",
			TableName: "user",
			SqlcType:  ":exec",
		},
		{
			Name:      "RestoreUser",
			Query:     "UPDATE \"user\" SET \"deleted_at\"=NULL WHERE ((\"id\" = $1) AND (\"deleted_at\" IS NOT NULL))",
			TableName: "user",
			SqlcType:  ":exec",
		},
		{
			Name:      "DestroyUser",
			Query:     "DELETE FROM \"user\" WHERE (\"id\" = $1)",
			TableName: "user",
			SqlcType:  ":exec",
		},
	}
	res := generator.GenerateQueries(element)
	assert.Equal(t, res, tc)
}

func TestGenerateRestoreQuery(t *testing.T) {
	dialect := goqu.Dialect("postgres")
	element := &config.Schema{
		Name: "user",
		Fields: []*config.Field{
			{
				Name: "name",
			},
			{
				Name: "deleted_at",
			},
		},
	}
	tc := &config.Function{
		Name:      "RestoreUser",
		Query:     "UPDATE \"user\" SET \"deleted_at\"=NULL WHERE ((\"id\" = $1) AND (\"deleted_at\" IS NOT NULL))",
		TableName: "user",
		SqlcType:  ":exec",
	}
	res := generator.GenerateRestoreQuery(dialect, element)
	assert.Equal(t, res, tc)
}

func TestGenerateRestoreQuery_ReturnNil(t *testing.T) {
	dialect := goqu.Dialect("postgres")
	element := &config.Schema{
		Name: "user",
		Fields: []*config.Field{
			{
				Name: "name",
			},
		},
	}
	res := generator.GenerateRestoreQuery(dialect, element)
	assert.Nil(t, res)
}
