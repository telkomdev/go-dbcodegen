package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/wartek-id/core/tools/dbgen/config"
)

func TestParseSchema(t *testing.T) {
	schema, err := config.ParseSchema("../examples/schemas/example.json")
	assert.NoError(t, err)
	assert.Equal(t, "example", schema.Name)
	assert.Len(t, schema.Fields, 2)
}

func TestParseDir(t *testing.T) {
	schemas, err := config.ParseDir("../examples/schemas")
	assert.NoError(t, err)
	assert.Len(t, schemas, 4)
}

func TestParse(t *testing.T) {
	// Parse Dir
	schemas, err := config.Parse("../examples/schemas")
	assert.NoError(t, err)
	assert.Len(t, schemas, 4)

	// Parse File
	schemas, err = config.Parse("../examples/schemas/example.json")
	assert.NoError(t, err)
	assert.Len(t, schemas, 1)

	// DirNotExists
	schemas, err = config.Parse("invalid_dir")
	assert.NotNil(t, err)
	assert.Empty(t, schemas)
}

func TestSchema_GetName(t *testing.T) {
	schema := config.Schema{
		Name: "users",
	}

	assert.Equal(t, "users", schema.GetName())
}
