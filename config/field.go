package config

import (
	"gitlab.com/wartek-id/core/tools/dbgen/types/field_option"
	"gitlab.com/wartek-id/core/tools/dbgen/types/field_type"
)

type Field struct {
	Name    string                     `json:"name"`
	Type    field_type.FieldType       `json:"type"`
	Scale   int                        `json:"scale"`
	Limit   int                        `json:"limit"`
	Default interface{}                `json:"default"`
	Options []field_option.FieldOption `json:"options"`
}

func (f *Field) GetName() string {
	return f.Name
}

func (f *Field) IsNotNull() bool {
	for _, opt := range f.Options {
		switch opt {
		case field_option.NotNull, field_option.PrimaryKey:
		return true
		case field_option.Nullable:
			return false
		}
	}

	return false
}
