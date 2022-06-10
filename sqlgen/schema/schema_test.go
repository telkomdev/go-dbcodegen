package schema_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/wartek-id/core/tools/dbgen/config"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/schema"
)

func TestNewSchema(t *testing.T) {
	sc, err := schema.NewSchema("mysql://root:root@localhost:3306/playground")
	assert.Equal(t, schema.ErrUnsupportedDriver, err)
	assert.Nil(t, sc)
}

func TestIndices_GetByConstraintName(t *testing.T) {
	pKeys := &config.Index{
		Name: "example_pkeys",
		Fields: []*config.IndexField{
			{Column: "id"},
		},
		Unique: true,
	}
	indices := schema.Indices{
		Table: "example",
		Indices: map[string]*config.Index{
			"example_pkeys": pKeys,
			"empty": {
				Name: "empty",
			},
		},
	}

	result, err := indices.GetByConstraintName("example_pkeys")
	assert.Nil(t, err)
	assert.Equal(t, pKeys, result)

	result, err = indices.GetByConstraintName("something")
	assert.Equal(t, schema.ErrConstraintNotExists, err)
	assert.Nil(t, result)

	result, err = indices.GetByConstraintName("empty")
	assert.Equal(t, schema.ErrConstraintHasNoFields, err)
	assert.Nil(t, result)
}
