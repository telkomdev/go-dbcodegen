package sqlgen_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/wartek-id/core/tools/dbgen/config"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/dialect"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/sb"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/step"
	"gitlab.com/wartek-id/core/tools/dbgen/types/field_option"
	"gitlab.com/wartek-id/core/tools/dbgen/types/field_type"
)

func TestAlterTableGenerator_Dialect(t *testing.T) {
	dial := "postgres"
	do := dialect.DefaultDialectOption()

	sqlGen := sqlgen.NewAlterTableGenerator(dial, do)
	assert.Equal(t, dial, sqlGen.Dialect())
}

func TestAlterTableGenerator_DialectOptions(t *testing.T) {
	dial := "postgres"
	do := dialect.DefaultDialectOption()

	sqlGen := sqlgen.NewAlterTableGenerator(dial, do)
	assert.Equal(t, do, sqlGen.DialectOptions())
}

func TestAlterTableGenerator_ExpressionSQLGenerator(t *testing.T) {
	dial := "postgres"
	do := dialect.DefaultDialectOption()

	sqlGen := sqlgen.NewAlterTableGenerator(dial, do)
	assert.NotNil(t, sqlGen.ExpressionSQLGenerator())
}

func TestAlterSchemaGenerator_Generate(t *testing.T) {
	alterStep := step.AlterSchema{
		Name: "users",
		AddedColumns: []*config.Field{
			{
				Name:  "location",
				Type:  field_type.Varchar,
				Limit: 50,
			},
			{
				Name:    "address",
				Type:    field_type.Varchar,
				Limit:   100,
				Options: []field_option.FieldOption{field_option.NotNull},
			},
		},
		DroppedColumns: []*config.Field{
			{
				Name: "stock",
			},
			{
				Name: "author",
			},
		},
		AlteredColumns: []*step.AlterColumn{
			{
				Name: "name",
				Field: &config.Field{
					Name:    "name",
					Type:    field_type.Varchar,
					Limit:   200,
					Default: "Alfred",
					Options: []field_option.FieldOption{
						field_option.NotNull,
					},
				},
				ChangedType:         true,
				ChangedDefaultValue: true,
				ChangedOptions: []step.OptionAction{
					step.SetNotNull,
				},
			},
			{
				Name: "price",
				Field: &config.Field{
					Name:    "price",
					Type:    field_type.BigInt,
					Default: 10000,
					Options: []field_option.FieldOption{
						field_option.Nullable,
					},
				},
				ChangedType:         true,
				ChangedDefaultValue: true,
				ChangedOptions: []step.OptionAction{
					step.DropNotNull,
				},
			},
		},
	}

	gen := sqlgen.NewAlterTableGenerator("postgres", dialect.DefaultDialectOption())
	buf := sb.NewSQLBuilder()
	gen.Generate(buf, &alterStep)
	result := fmt.Sprintf("%s\n%s,\n%s,\n%s,\n%s,\n%s,\n%s;",
		"ALTER TABLE IF EXISTS \"users\"",
		"\tADD COLUMN \"location\" VARCHAR(50)",
		"\tADD COLUMN \"address\" VARCHAR(100) NOT NULL",
		"\tDROP COLUMN \"stock\"",
		"\tDROP COLUMN \"author\"",
		"\tALTER COLUMN \"name\" SET DATA TYPE VARCHAR(200),\n\tALTER COLUMN \"name\" SET NOT NULL,\n\tALTER COLUMN \"name\" SET DEFAULT 'Alfred'",
		"\tALTER COLUMN \"price\" SET DATA TYPE BIGINT,\n\tALTER COLUMN \"price\" DROP NOT NULL,\n\tALTER COLUMN \"price\" SET DEFAULT 10000",
	)
	assert.Equal(t, result, buf.String())
}

func TestAlterSchemaGenerator_Rollback(t *testing.T) {
	alterStep := step.AlterSchema{
		Name: "users",
		AddedColumns: []*config.Field{
			{
				Name:  "location",
				Type:  field_type.Varchar,
				Limit: 50,
			},
			{
				Name:    "address",
				Type:    field_type.Varchar,
				Limit:   100,
				Options: []field_option.FieldOption{field_option.NotNull},
			},
		},
		DroppedColumns: []*config.Field{
			{
				Name:  "stock",
				Type:  field_type.Decimal,
				Limit: 2,
				Scale: 10,
			},
			{
				Name:  "author",
				Type:  field_type.Varchar,
				Limit: 20,
			},
		},
		AlteredColumns: []*step.AlterColumn{
			{
				Name: "name",
				Field: &config.Field{
					Name:    "name",
					Type:    field_type.Varchar,
					Limit:   200,
					Default: "Alfred",
					Options: []field_option.FieldOption{
						field_option.NotNull,
					},
				},
				LastField: &config.Field{
					Name:    "name",
					Type:    field_type.Varchar,
					Limit:   100,
					Default: "Tejo",
					Options: []field_option.FieldOption{
						field_option.Nullable,
					},
				},
				ChangedType:         true,
				ChangedDefaultValue: true,
				ChangedOptions: []step.OptionAction{
					step.SetNotNull,
				},
			},
			{
				Name: "price",
				Field: &config.Field{
					Name:    "price",
					Type:    field_type.BigInt,
					Default: 10000,
					Options: []field_option.FieldOption{
						field_option.Nullable,
					},
				},
				LastField: &config.Field{
					Name:    "price",
					Type:    field_type.Decimal,
					Limit:   2,
					Scale:   10,
					Default: 50000,
					Options: []field_option.FieldOption{
						field_option.NotNull,
					},
				},
				ChangedType:         true,
				ChangedDefaultValue: true,
				ChangedOptions: []step.OptionAction{
					step.DropNotNull,
				},
			},
		},
	}

	gen := sqlgen.NewAlterTableGenerator("postgres", dialect.DefaultDialectOption())
	buf := sb.NewSQLBuilder()
	gen.Rollback(buf, &alterStep)
	result := fmt.Sprintf("%s\n%s,\n%s,\n%s,\n%s,\n%s,\n%s;",
		"ALTER TABLE IF EXISTS \"users\"",
		"\tDROP COLUMN \"location\"",
		"\tDROP COLUMN \"address\"",
		"\tADD COLUMN \"stock\" DECIMAL(2, 10)",
		"\tADD COLUMN \"author\" VARCHAR(20)",
		"\tALTER COLUMN \"name\" SET DATA TYPE VARCHAR(100),\n\tALTER COLUMN \"name\" DROP NOT NULL,\n\tALTER COLUMN \"name\" SET DEFAULT 'Tejo'",
		"\tALTER COLUMN \"price\" SET DATA TYPE DECIMAL(2, 10),\n\tALTER COLUMN \"price\" SET NOT NULL,\n\tALTER COLUMN \"price\" SET DEFAULT 50000",
	)
	assert.Equal(t, result, buf.String())
}
