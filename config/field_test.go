package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/wartek-id/core/tools/dbgen/config"
	"gitlab.com/wartek-id/core/tools/dbgen/types/field_option"
)

func TestField_GetName(t *testing.T) {
	field := config.Field{
		Name: "users",
	}

	assert.Equal(t, "users", field.GetName())
}

func TestField_IsNotNull(t *testing.T) {
	field := config.Field{
		Name: "users",
		Options: []field_option.FieldOption{
			field_option.NotNull,
		},
	}

	assert.True(t, field.IsNotNull())

	field = config.Field{
		Name: "users",
		Options: []field_option.FieldOption{
			field_option.Nullable,
		},
	}

	assert.False(t, field.IsNotNull())

	field = config.Field{
		Name: "users",
	}

	assert.False(t, field.IsNotNull())
}
