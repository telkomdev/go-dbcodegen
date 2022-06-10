package sqlgen

import (
	"gitlab.com/wartek-id/core/tools/dbgen/config"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/dialect"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/exp"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/sb"
)

type CreateTableGenerator interface {
	Dialect() string
	DialectOptions() *dialect.DialectOption
	ExpressionSQLGenerator() exp.ExpressionSQLGenerator
	Generate(sb.SQLBuilder, *config.Schema)
}

type createTableGenerator struct {
	dialect        string
	esg            exp.ExpressionSQLGenerator
	dialectOptions *dialect.DialectOption
}

func NewCreateTableGenerator(dialect string, do *dialect.DialectOption) CreateTableGenerator {
	return &createTableGenerator{
		dialect:        dialect,
		dialectOptions: do,
		esg:            exp.NewExpressionSQLGenerator(dialect, do),
	}
}

func (ctg *createTableGenerator) Dialect() string {
	return ctg.dialect
}

func (ctg *createTableGenerator) DialectOptions() *dialect.DialectOption {
	return ctg.dialectOptions
}

func (ctg *createTableGenerator) Generate(b sb.SQLBuilder, schema *config.Schema) {
	b.Write(ctg.dialectOptions.CreateClause).
		Write(ctg.dialectOptions.TableFragment).
		Write(ctg.dialectOptions.IfNotExistsFragment)

	ctg.ExpressionSQLGenerator().LiteralExpression(b, schema.Name)

	b.WriteRunes(ctg.dialectOptions.SpaceRune)
	b.WriteRunes(ctg.dialectOptions.LeftParenRune)
	b.WriteRunes(ctg.dialectOptions.NewLineRune)
	ctg.FieldSQL(b, schema.Fields)
	b.WriteRunes(ctg.dialectOptions.NewLineRune)
	b.WriteRunes(ctg.dialectOptions.RightParenRune)
	b.WriteRunes(ctg.dialectOptions.SemiColonRune)
}

func (ctg *createTableGenerator) FieldSQL(b sb.SQLBuilder, fields []*config.Field) {
	for i, field := range fields {
		b.WriteRunes(ctg.dialectOptions.TabRune)
		ctg.ExpressionSQLGenerator().LiteralExpression(b, field.Name)
		b.WriteRunes(ctg.dialectOptions.SpaceRune)
		b.Write(ctg.esg.GetTypeFragment(field))
		b.Write(ctg.esg.GetOptionsFragment(field))

		if i != len(fields)-1 {
			b.Write(ctg.dialectOptions.CommaNewLineFragment)
		}
	}
}

func (ctg *createTableGenerator) ExpressionSQLGenerator() exp.ExpressionSQLGenerator {
	return ctg.esg
}
