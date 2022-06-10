package sqlgen

import (
	"gitlab.com/wartek-id/core/tools/dbgen/config"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/dialect"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/exp"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/sb"
)

type DropIndexGenerator interface {
	Dialect() string
	DialectOptions() *dialect.DialectOption
	ExpressionSQLGenerator() exp.ExpressionSQLGenerator
	Generate(sb.SQLBuilder, *config.Index)
}

type dropIndexGenerator struct {
	dialect        string
	esg            exp.ExpressionSQLGenerator
	dialectOptions *dialect.DialectOption
}

func NewDropIndexGenerator(dialect string, do *dialect.DialectOption) DropIndexGenerator {
	return &dropIndexGenerator{
		dialect:        dialect,
		dialectOptions: do,
		esg:            exp.NewExpressionSQLGenerator(dialect, do),
	}
}

func (dig *dropIndexGenerator) Dialect() string {
	return dig.dialect
}

func (dig *dropIndexGenerator) DialectOptions() *dialect.DialectOption {
	return dig.dialectOptions
}

func (dig *dropIndexGenerator) ExpressionSQLGenerator() exp.ExpressionSQLGenerator {
	return dig.esg
}

func (dig *dropIndexGenerator) Generate(b sb.SQLBuilder, index *config.Index) {
	b.Write(dig.dialectOptions.DropClause).
		Write(dig.dialectOptions.IndexFragment).
		Write(dig.dialectOptions.IfExistsFragment)

	dig.ExpressionSQLGenerator().LiteralExpression(b, index.Name)
	b.WriteRunes(dig.dialectOptions.SemiColonRune)
}
