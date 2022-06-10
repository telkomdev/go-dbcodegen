package sqlgen

import (
	"gitlab.com/wartek-id/core/tools/dbgen/config"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/dialect"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/exp"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/sb"
)

type DropTableGenerator interface {
	Dialect() string
	DialectOptions() *dialect.DialectOption
	ExpressionSQLGenerator() exp.ExpressionSQLGenerator
	Generate(sb.SQLBuilder, *config.Schema)
}

type dropTableGenerator struct {
	dialect        string
	esg            exp.ExpressionSQLGenerator
	dialectOptions *dialect.DialectOption
}

func NewDropTableGenerator(dialect string, do *dialect.DialectOption) DropTableGenerator {
	return &dropTableGenerator{
		dialect:        dialect,
		dialectOptions: do,
		esg:            exp.NewExpressionSQLGenerator(dialect, do),
	}
}

func (dtg *dropTableGenerator) Dialect() string {
	return dtg.dialect
}

func (dtg *dropTableGenerator) DialectOptions() *dialect.DialectOption {
	return dtg.dialectOptions
}

func (dtg *dropTableGenerator) ExpressionSQLGenerator() exp.ExpressionSQLGenerator {
	return dtg.esg
}

func (dtg *dropTableGenerator) Generate(b sb.SQLBuilder, schema *config.Schema) {
	b.Write(dtg.dialectOptions.DropClause).
		Write(dtg.dialectOptions.TableFragment).
		Write(dtg.dialectOptions.IfExistsFragment)

	dtg.ExpressionSQLGenerator().LiteralExpression(b, schema.Name)
	b.WriteRunes(dtg.dialectOptions.SemiColonRune)
}
