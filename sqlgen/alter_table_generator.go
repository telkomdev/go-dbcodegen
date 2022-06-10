package sqlgen

import (
	"bytes"

	"gitlab.com/wartek-id/core/tools/dbgen/config"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/dialect"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/exp"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/sb"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/step"
)

type AlterTableGenerator interface {
	Dialect() string
	DialectOptions() *dialect.DialectOption
	ExpressionSQLGenerator() exp.ExpressionSQLGenerator
	Generate(b sb.SQLBuilder, at *step.AlterSchema) error
	Rollback(b sb.SQLBuilder, at *step.AlterSchema) error
}

type alterTableGenerator struct {
	dialect        string
	esg            exp.ExpressionSQLGenerator
	dialectOptions *dialect.DialectOption
}

func NewAlterTableGenerator(dialect string, do *dialect.DialectOption) AlterTableGenerator {
	return &alterTableGenerator{
		dialect:        dialect,
		dialectOptions: do,
		esg:            exp.NewExpressionSQLGenerator(dialect, do),
	}
}

func (atg *alterTableGenerator) Dialect() string {
	return atg.dialect
}

func (atg *alterTableGenerator) DialectOptions() *dialect.DialectOption {
	return atg.dialectOptions
}

func (atg *alterTableGenerator) ExpressionSQLGenerator() exp.ExpressionSQLGenerator {
	return atg.esg
}

func (atg *alterTableGenerator) Generate(b sb.SQLBuilder, at *step.AlterSchema) error {
	if !at.FieldChanged() {
		return nil
	}

	atg.alterTableTemplate(b, at.Name)
	queries := make([][]byte, 0)

	if at.IsColumnsAdded() {
		buf := sb.NewSQLBuilder()
		atg.generateColumns(buf, at.AddedColumns)
		queries = append(queries, buf.Bytes())
	}

	if at.IsColumnsDropped() {
		buf := sb.NewSQLBuilder()
		atg.dropColumns(buf, at.DroppedColumns)
		queries = append(queries, buf.Bytes())
	}

	if at.IsColumnsAltered() {
		buf := sb.NewSQLBuilder()
		atg.alterColumns(buf, at.AlteredColumns)
		queries = append(queries, buf.Bytes())
	}

	b.Write(bytes.Join(queries, atg.dialectOptions.CommaNewLineFragment))
	b.WriteRunes(atg.dialectOptions.SemiColonRune)
	return nil
}

func (atg *alterTableGenerator) generateColumns(b sb.SQLBuilder, fields []*config.Field) {
	for i, field := range fields {
		b.WriteRunes(atg.dialectOptions.TabRune)
		b.Write(atg.dialectOptions.AddColumnTemplate())
		atg.ExpressionSQLGenerator().LiteralExpression(b, field.Name)
		b.WriteRunes(atg.dialectOptions.SpaceRune)
		b.Write(atg.ExpressionSQLGenerator().GetTypeFragment(field))
		b.Write(atg.ExpressionSQLGenerator().GetOptionsFragment(field))

		if i != len(fields)-1 {
			b.WriteRunes(atg.dialectOptions.CommaRune)
			b.WriteRunes(atg.dialectOptions.NewLineRune)
		}
	}
}

func (atg *alterTableGenerator) dropColumns(b sb.SQLBuilder, fields []*config.Field) {
	for i, field := range fields {
		b.WriteRunes(atg.dialectOptions.TabRune)
		b.Write(atg.dialectOptions.DropColumnTemplate())
		atg.ExpressionSQLGenerator().LiteralExpression(b, field.Name)

		if i != len(fields)-1 {
			b.WriteRunes(atg.dialectOptions.CommaRune)
			b.WriteRunes(atg.dialectOptions.NewLineRune)
		}
	}
}

func (atg *alterTableGenerator) alterColumns(b sb.SQLBuilder, fields []*step.AlterColumn) {
	for i, field := range fields {
		atg.alterColumn(b, field)
		if i != len(fields)-1 {
			b.Write(atg.dialectOptions.CommaNewLineFragment)
		}
	}
}

func (atg *alterTableGenerator) alterColumn(b sb.SQLBuilder, field *step.AlterColumn) {
	changes := [][]byte{}
	if field.ChangedType {
		buf := sb.NewSQLBuilder()
		atg.changeColumnType(buf, field.Field)
		changes = append(changes, buf.Bytes())
	}

	if field.IsOptionsChanged() {
		buf := sb.NewSQLBuilder()
		atg.ChangeColumnOptions(buf, field)
		changes = append(changes, buf.Bytes())
	}

	if field.ChangedDefaultValue {
		buf := sb.NewSQLBuilder()
		atg.changeColumnDefault(buf, field.Field)
		changes = append(changes, buf.Bytes())
	}

	b.Write(bytes.Join(changes, atg.dialectOptions.CommaNewLineFragment))
}

func (atg *alterTableGenerator) changeColumnType(b sb.SQLBuilder, field *config.Field) {
	atg.alterColumnTemplate(b, field.Name)
	b.Write(atg.dialectOptions.SetFragment)
	b.Write(atg.dialectOptions.DataTypeFragment)
	b.Write(atg.ExpressionSQLGenerator().GetTypeFragment(field))
}

func (atg *alterTableGenerator) changeColumnDefault(b sb.SQLBuilder, field *config.Field) {
	atg.alterColumnTemplate(b, field.Name)
	b.Write(atg.dialectOptions.SetFragment)
	b.Write(atg.dialectOptions.DefaultFragment)
	b.Write(atg.ExpressionSQLGenerator().GetDefaultValue(field.Default))
}

func (atg *alterTableGenerator) ChangeColumnOptions(b sb.SQLBuilder, field *step.AlterColumn) {
	for _, option := range field.ChangedOptions {
		switch option {
		case step.SetNotNull:
			atg.setNotNull(b, field.Name)
		case step.DropNotNull:
			atg.dropNotNull(b, field.Name)
		}
	}
}

func (atg *alterTableGenerator) Rollback(b sb.SQLBuilder, at *step.AlterSchema) error {
	atg.alterTableTemplate(b, at.Name)
	queries := make([][]byte, 0)

	if at.IsColumnsAdded() {
		buf := sb.NewSQLBuilder()
		atg.dropColumns(buf, at.AddedColumns)
		queries = append(queries, buf.Bytes())
	}

	if at.IsColumnsDropped() {
		buf := sb.NewSQLBuilder()
		atg.generateColumns(buf, at.DroppedColumns)
		queries = append(queries, buf.Bytes())
	}

	if at.IsColumnsAltered() {
		buf := sb.NewSQLBuilder()
		atg.rollbackAlterColumns(buf, at.AlteredColumns)
		queries = append(queries, buf.Bytes())
	}

	b.Write(bytes.Join(queries, atg.dialectOptions.CommaNewLineFragment))
	b.WriteRunes(atg.dialectOptions.SemiColonRune)
	return nil
}

func (atg *alterTableGenerator) rollbackAlterColumns(b sb.SQLBuilder, fields []*step.AlterColumn) {
	for i, field := range fields {
		atg.rollbackAlterColumn(b, field)
		if i != len(fields)-1 {
			b.Write(atg.dialectOptions.CommaNewLineFragment)
		}
	}
}

func (atg *alterTableGenerator) rollbackAlterColumn(b sb.SQLBuilder, field *step.AlterColumn) {
	changes := [][]byte{}
	if field.ChangedType {
		buf := sb.NewSQLBuilder()
		atg.changeColumnType(buf, field.LastField)
		changes = append(changes, buf.Bytes())
	}

	if field.IsOptionsChanged() {
		buf := sb.NewSQLBuilder()
		atg.rollbackChangeColumnOptions(buf, field)
		changes = append(changes, buf.Bytes())
	}

	if field.ChangedDefaultValue {
		buf := sb.NewSQLBuilder()
		atg.changeColumnDefault(buf, field.LastField)
		changes = append(changes, buf.Bytes())
	}

	b.Write(bytes.Join(changes, atg.dialectOptions.CommaNewLineFragment))
}

func (atg *alterTableGenerator) rollbackChangeColumnOptions(b sb.SQLBuilder, field *step.AlterColumn) {
	for _, option := range field.ChangedOptions {
		switch option {
		case step.SetNotNull:
			atg.dropNotNull(b, field.Name)
		case step.DropNotNull:
			atg.setNotNull(b, field.Name)
		}
	}
}

func (atg *alterTableGenerator) dropNotNull(b sb.SQLBuilder, name string) {
	atg.alterColumnTemplate(b, name)
	b.Write(atg.dialectOptions.DropFragment)
	b.Write(atg.dialectOptions.NotNullFragment)
}

func (atg *alterTableGenerator) setNotNull(b sb.SQLBuilder, name string) {
	atg.alterColumnTemplate(b, name)
	b.Write(atg.dialectOptions.SetFragment)
	b.Write(atg.dialectOptions.NotNullFragment)
}

func (atg *alterTableGenerator) alterColumnTemplate(b sb.SQLBuilder, name string) {
	b.WriteRunes(atg.dialectOptions.TabRune)
	b.Write(atg.dialectOptions.AlterColumnTemplate())
	atg.ExpressionSQLGenerator().LiteralExpression(b, name)
	b.WriteRunes(atg.dialectOptions.SpaceRune)
}

func (atg *alterTableGenerator) alterTableTemplate(b sb.SQLBuilder, name string) {
	b.Write(atg.dialectOptions.AlterClause)
	b.Write(atg.dialectOptions.TableFragment)
	b.Write(atg.dialectOptions.IfExistsFragment)
	atg.ExpressionSQLGenerator().LiteralExpression(b, name)
	b.WriteRunes(atg.dialectOptions.NewLineRune)
}
