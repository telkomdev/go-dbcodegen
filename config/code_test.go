package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/wartek-id/core/tools/dbgen/config"
)

func TestParseFunction(t *testing.T) {
	function := config.Function{
		Name:      "Name",
		Query:     "Query",
		TableName: "TableName",
		SqlcType:  "one",
	}
	result := config.ParseFunction("Name", "Query", "TableName", "one")
	assert.Equal(t, result.Name, function.Name)
	assert.Equal(t, result.Query, function.Query)
	assert.Equal(t, result.TableName, function.TableName)
	assert.Equal(t, result.SqlcType, function.SqlcType)
}
