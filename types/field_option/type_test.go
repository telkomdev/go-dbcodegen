package field_option_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/wartek-id/core/tools/dbgen/types/field_option"
)

func TestFieldOption_UnmarshallJSON(t *testing.T) {
	testCases := map[string]struct {
		input   []byte
		wantErr error
		result  field_option.FieldOption
	}{
		"success": {
			input:  []byte("\"not null\""),
			result: "not null",
		},
		"invalid type": {
			input:   []byte("\"binary\""),
			wantErr: fmt.Errorf("invalid \"binary\" as field option"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var result field_option.FieldOption
			err := json.Unmarshal(tc.input, &result)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.result, result)
		})
	}
}
