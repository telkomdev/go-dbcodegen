package field_option

import (
	"encoding/json"
	"fmt"
	"strings"
)

type FieldOption string

const (
	Nullable      FieldOption = "nullable"
	NotNull       FieldOption = "not null"
	AutoIncrement FieldOption = "auto increment"
	Unique        FieldOption = "unique"
	PrimaryKey    FieldOption = "primary key"
)

var SupportedOption = []FieldOption{
	Nullable,
	NotNull,
	AutoIncrement,
	Unique,
	PrimaryKey,
}

func (o *FieldOption) UnmarshalJSON(data []byte) error {
	var strOpt string
	err := json.Unmarshal(data, &strOpt)
	if err != nil {
		return err
	}

	fo := FieldOption(strings.ToLower(strOpt))
	for _, opt := range SupportedOption {
		if fo == opt {
			*o = fo
			return nil
		}
	}
	return fmt.Errorf("invalid \"%s\" as field option", strOpt)
}
