package exp

import (
	"bytes"
	"fmt"

	"gitlab.com/wartek-id/core/tools/dbgen/config"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/dialect"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/sb"
)

type ExpressionSQLGenerator interface {
	GetTypeFragment(field *config.Field) []byte
	GetOptionsFragment(field *config.Field) []byte
	LiteralExpression(buf sb.SQLBuilder, value string)
	GetDefaultValue(value interface{}) []byte
}

type expressionSQLGenerator struct {
	dialect        string
	dialectOptions *dialect.DialectOption
}

func NewExpressionSQLGenerator(dialect string, do *dialect.DialectOption) ExpressionSQLGenerator {
	return &expressionSQLGenerator{
		dialect:        dialect,
		dialectOptions: do,
	}
}

func (ex *expressionSQLGenerator) GetTypeFragment(field *config.Field) []byte {
	buf := sb.NewSQLBuilder()
	buf.Write(ex.dialectOptions.DataTypesLookup[field.Type])
	if field.Limit == 0 && field.Scale == 0 || !field.Type.HasLimit() {
		return buf.Bytes()
	}

	buf.WriteRunes(ex.dialectOptions.LeftParenRune).
		WriteString(fmt.Sprint(field.Limit))
	if field.Scale != 0 && field.Type.HasScale() {
		buf.WriteRunes(ex.dialectOptions.CommaRune, ex.dialectOptions.SpaceRune).
			WriteString(fmt.Sprint(field.Scale))
	}
	buf.WriteRunes(ex.dialectOptions.RightParenRune)
	return buf.Bytes()
}

func (ex *expressionSQLGenerator) GetOptionsFragment(field *config.Field) []byte {
	if len(field.Options) <= 0 {
		return []byte{}
	}
	options := make([][]byte, 0, len(field.Options))
	options = append(options, ex.dialectOptions.EmptyFragment)

	for _, opts := range field.Options {
		options = append(options, ex.dialectOptions.FieldOptionsLookup[opts])
	}

	return bytes.Join(options, []byte(string(ex.dialectOptions.SpaceRune)))
}

func (ex *expressionSQLGenerator) LiteralExpression(buf sb.SQLBuilder, value string) {
	buf.WriteRunes(ex.dialectOptions.QuoteRune)
	buf.WriteString(value)
	buf.WriteRunes(ex.dialectOptions.QuoteRune)
}

func (ex *expressionSQLGenerator) GetDefaultValue(value interface{}) []byte {
	switch v := value.(type) {
	case string:
		buf := bytes.Buffer{}
		buf.WriteRune(ex.dialectOptions.StringQuoteRune)
		buf.WriteString(v)
		buf.WriteRune(ex.dialectOptions.StringQuoteRune)
		return buf.Bytes()
	}

	return []byte(fmt.Sprint(value))
}
