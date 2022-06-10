package field_type

import (
	"encoding/json"
	"fmt"
	"strings"
)

type FieldType string

const (
	Boolean FieldType = "bool"
	Varchar FieldType = "varchar"
	Text    FieldType = "text"

	SmallInt FieldType = "smallint"
	Int      FieldType = "int"
	BigInt   FieldType = "bigint"

	Json  FieldType = "json"
	Jsonb FieldType = "jsonb"

	Float   FieldType = "float"
	Decimal FieldType = "decimal"

	Timestamp   FieldType = "timestamp"
	Timestamptz FieldType = "timestamptz"

	BigSerial   FieldType = "bigserial"
	Serial      FieldType = "serial"
	SmallSerial FieldType = "smallserial"

	FieldTypeString  = "string"
	FieldTypeNumeric = "numeric"
	FieldTypeBinary  = "binary"
)

var SupportedFieldType = []FieldType{
	Boolean,
	Varchar,
	Text,
	SmallInt,
	Int,
	BigInt,
	Json,
	Jsonb,
	Decimal,
	Float,
	Timestamp,
	Timestamptz,
	BigSerial,
	Serial,
	SmallSerial,
}

func (t *FieldType) UnmarshalJSON(data []byte) error {
	var strType string
	err := json.Unmarshal(data, &strType)
	if err != nil {
		return err
	}

	ft := FieldType(strings.ToLower(strType))
	for _, typ := range SupportedFieldType {
		if ft == typ {
			*t = ft
			return nil
		}
	}

	return fmt.Errorf("invalid \"%s\" as field type", strType)
}

func (t FieldType) Type() string {
	switch t {
	case Varchar, Text, Json:
		return FieldTypeString
	case Jsonb:
		return FieldTypeBinary
	}
	return FieldTypeNumeric
}

func (t FieldType) HasLimit() bool {
	return t == Varchar || t == Decimal
}

func (t FieldType) HasScale() bool {
	return t == Decimal
}

func ParseString(ft string) FieldType {
	return FieldType(ft)
}
