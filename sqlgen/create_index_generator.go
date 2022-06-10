package sqlgen

import (
	"gitlab.com/wartek-id/core/tools/dbgen/config"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/dialect"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/exp"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/sb"
)

type CreateIndexGenerator interface {
	Dialect() string
	DialectOptions() *dialect.DialectOption
	ExpressionSQLGenerator() exp.ExpressionSQLGenerator
	Generate(sb.SQLBuilder, string, *config.Index)
}

type createIndexGenerator struct {
	dialect        string
	esg            exp.ExpressionSQLGenerator
	dialectOptions *dialect.DialectOption
}

func NewCreateIndexGenerator(dialect string, do *dialect.DialectOption) CreateIndexGenerator {
	return &createIndexGenerator{
		dialect:        dialect,
		dialectOptions: do,
		esg:            exp.NewExpressionSQLGenerator(dialect, do),
	}
}

func (cig *createIndexGenerator) Dialect() string {
	return cig.dialect
}

func (cig *createIndexGenerator) DialectOptions() *dialect.DialectOption {
	return cig.dialectOptions
}

func (cig *createIndexGenerator) ExpressionSQLGenerator() exp.ExpressionSQLGenerator {
	return cig.esg
}

func (cig *createIndexGenerator) Generate(b sb.SQLBuilder, tblName string, idx *config.Index) {
	b.Write(cig.dialectOptions.CreateClause)
	if idx.Unique {
		b.Write(cig.dialectOptions.UniqueFragment).
			WriteRunes(cig.dialectOptions.SpaceRune)
	}
	b.Write(cig.dialectOptions.IndexFragment)

	if cig.dialectOptions.SupportConcurrently {
		b.Write(cig.dialectOptions.ConcurrentlyFragment).
			WriteRunes(cig.dialectOptions.SpaceRune)
	}

	b.Write(cig.dialectOptions.IfNotExistsFragment)
	cig.ExpressionSQLGenerator().LiteralExpression(b, idx.Name)
	b.Write(cig.dialectOptions.OnFragment)
	cig.ExpressionSQLGenerator().LiteralExpression(b, tblName)
	b.WriteRunes(cig.dialectOptions.LeftParenRune)
	cig.FieldSQL(b, idx.Fields)
	b.WriteRunes(cig.dialectOptions.RightParenRune)
	b.WriteRunes(cig.dialectOptions.SemiColonRune)
}

func (cig *createIndexGenerator) FieldSQL(b sb.SQLBuilder, fields []*config.IndexField) {
	for i, field := range fields {
		cig.ExpressionSQLGenerator().LiteralExpression(b, field.Column)
		if field.Order != "" {
			b.WriteRunes(cig.dialectOptions.SpaceRune)
			b.WriteString(field.Order)
		}
		if i != len(fields)-1 {
			b.WriteRunes(cig.dialectOptions.CommaRune)
			b.WriteRunes(cig.dialectOptions.SpaceRune)
		}
	}
}
