package field_type_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/wartek-id/core/tools/dbgen/types/field_type"
)

func TestFieldType_UnmarshallJSON(t *testing.T) {
	testCases := map[string]struct {
		input   []byte
		wantErr error
		result  field_type.FieldType
	}{
		"success": {
			input:  []byte("\"bool\""),
			result: "bool",
		},
		"invalid type": {
			input:   []byte("\"binary\""),
			wantErr: fmt.Errorf("invalid \"binary\" as field type"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var result field_type.FieldType
			err := json.Unmarshal(tc.input, &result)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.result, result)
		})
	}
}

func TestFieldType_Type(t *testing.T) {
	testCases := []struct {
		input  field_type.FieldType
		result string
	}{
		{
			field_type.Varchar,
			field_type.FieldTypeString,
		},
		{
			field_type.Text,
			field_type.FieldTypeString,
		},
		{
			field_type.Json,
			field_type.FieldTypeString,
		},
		{
			field_type.Jsonb,
			field_type.FieldTypeBinary,
		},
		{
			field_type.Int,
			field_type.FieldTypeNumeric,
		},
		{
			field_type.Float,
			field_type.FieldTypeNumeric,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.result, tc.input.Type())
	}
}

func TestFieldType_HasLimit(t *testing.T) {
	assert.True(t, field_type.Varchar.HasLimit())
	assert.True(t, field_type.Decimal.HasLimit())
	assert.False(t, field_type.BigInt.HasLimit())
}

func TestFieldType_HasScale(t *testing.T) {
	assert.False(t, field_type.Varchar.HasScale())
	assert.True(t, field_type.Decimal.HasScale())
	assert.False(t, field_type.BigInt.HasScale())
}

func TestParseString(t *testing.T) {
	ft := field_type.ParseString("bigint")
	assert.Equal(t, field_type.FieldType("bigint"), ft)
}
