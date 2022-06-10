package config_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/wartek-id/core/tools/dbgen/config"
)

func TestIndex_GetName(t *testing.T) {
	index := config.Index{
		Name: "users",
	}

	assert.Equal(t, "users", index.GetName())
}

func TestIndex_UnmarshallJSON(t *testing.T) {
	testCases := []struct {
		input  []byte
		result *config.IndexField
	}{
		{
			input: []byte(`{"column": "name"}`),
			result: &config.IndexField{
				Column: "name",
				Order:  "ASC",
			},
		},
		{
			input: []byte(`{"column": "name", "order": "ASC"}`),
			result: &config.IndexField{
				Column: "name",
				Order:  "ASC",
			},
		},
		{
			input: []byte(`{"column": "name", "order": "DESC"}`),
			result: &config.IndexField{
				Column: "name",
				Order:  "DESC",
			},
		},
	}

	for _, tc := range testCases {
		var result config.IndexField
		err := json.Unmarshal(tc.input, &result)

		assert.Nil(t, err)
		assert.Equal(t, tc.result, &result)
	}
}

func TestIndex_GetColumns(t *testing.T) {
	index := config.Index{
		Name: "index_users_on_name_and_email",
		Fields: []*config.IndexField{
			{
				Column: "name",
			},
			{
				Column: "email",
			},
		},
	}

	assert.Equal(t, []string{"name", "email"}, index.GetColumns())
}
