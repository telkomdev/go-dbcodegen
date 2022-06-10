package diff_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/wartek-id/core/tools/dbgen/config"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/diff"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/step"
	"gitlab.com/wartek-id/core/tools/dbgen/types/field_option"
)

func TestCreatedTable(t *testing.T) {
	from := []*config.Schema{
		{
			Name: "user",
		},
		{
			Name: "module",
		},
	}

	target := []*config.Schema{
		{
			Name: "user",
		},
		{
			Name: "module",
		},
		{
			Name: "toggle",
		},
	}

	diffSchema := diff.NewSchema(from, target)
	result := diffSchema.CreatedTable()

	assert.Equal(t, []*config.Schema{{Name: "toggle"}}, result)
}

func TestDroppedTable(t *testing.T) {
	from := []*config.Schema{
		{
			Name: "user",
		},
		{
			Name: "module",
		},
		{
			Name: "toggle",
		},
	}

	target := []*config.Schema{
		{
			Name: "user",
		},
		{
			Name: "module",
		},
	}

	diffSchema := diff.NewSchema(from, target)
	result := diffSchema.DroppedTable()

	assert.Equal(t, []*config.Schema{{Name: "toggle"}}, result)
}

func TestAlterSchema(t *testing.T) {
	testCases := []struct {
		existing []*config.Schema
		target   []*config.Schema
		input    string
		result   *step.AlterSchema
		err      error
	}{
		{
			existing: []*config.Schema{
				{
					Name: "users",
					Fields: []*config.Field{
						{
							Name: "id",
							Type: "bigserial",
							Options: []field_option.FieldOption{
								field_option.NotNull,
							},
						},
						{
							Name:  "name",
							Type:  "varchar",
							Limit: 20,
						},
						{
							Name:  "address",
							Type:  "varchar",
							Limit: 50,
						},
						{
							Name:  "school",
							Type:  "varchar",
							Limit: 50,
						},
						{
							Name: "age",
							Type: "int",
							Options: []field_option.FieldOption{
								field_option.NotNull,
							},
						},
					},
				},
			},
			target: []*config.Schema{
				{
					Name: "users",
					Fields: []*config.Field{
						{
							Name: "id",
							Type: "bigserial",
							Options: []field_option.FieldOption{
								field_option.NotNull,
							},
						},
						{
							Name:  "name",
							Type:  "varchar",
							Limit: 20,
							Options: []field_option.FieldOption{
								field_option.NotNull,
							},
						},
						{
							Name:  "address",
							Type:  "varchar",
							Limit: 150,
						},
						{
							Name:  "school",
							Type:  "varchar",
							Limit: 100,
							Options: []field_option.FieldOption{
								field_option.NotNull,
							},
						},
						{
							Name: "created_at",
							Type: "timestamp",
							Options: []field_option.FieldOption{
								field_option.NotNull,
							},
						},
					},
				},
			},
			input: "users",
			result: &step.AlterSchema{
				Name: "users",
				AddedColumns: []*config.Field{
					{
						Name: "created_at",
						Type: "timestamp",
						Options: []field_option.FieldOption{
							field_option.NotNull,
						},
					},
				},
				AlteredColumns: []*step.AlterColumn{
					{
						Name: "name",
						Field: &config.Field{
							Name:  "name",
							Type:  "varchar",
							Limit: 20,
							Options: []field_option.FieldOption{
								field_option.NotNull,
							},
						},
						LastField: &config.Field{
							Name:  "name",
							Type:  "varchar",
							Limit: 20,
						},
						ChangedOptions: []step.OptionAction{step.SetNotNull},
					},
					{
						Name: "address",
						Field: &config.Field{
							Name:  "address",
							Type:  "varchar",
							Limit: 150,
						},
						LastField: &config.Field{
							Name:  "address",
							Type:  "varchar",
							Limit: 50,
						},
						ChangedType:    true,
						ChangedOptions: []step.OptionAction{},
					},
					{
						Name: "school",
						Field: &config.Field{
							Name:  "school",
							Type:  "varchar",
							Limit: 100,
							Options: []field_option.FieldOption{
								field_option.NotNull,
							},
						},
						LastField: &config.Field{
							Name:  "school",
							Type:  "varchar",
							Limit: 50,
						},
						ChangedType:    true,
						ChangedOptions: []step.OptionAction{step.SetNotNull},
					},
				},
				DroppedColumns: []*config.Field{
					{
						Name: "age",
						Type: "int",
						Options: []field_option.FieldOption{
							field_option.NotNull,
						},
					},
				},
			},
		},
		{
			existing: []*config.Schema{
				{
					Name: "users",
				},
			},
			target: []*config.Schema{},
			input:  "users",
			err:    diff.ErrMissingTargetTable,
		},
		{
			target: []*config.Schema{
				{
					Name: "users",
				},
			},
			existing: []*config.Schema{},
			input:    "users",
			err:      diff.ErrMissingCurrentTable,
		},
		{
			existing: []*config.Schema{
				{
					Name: "users",
					Index: []*config.Index{
						{
							Name: "index_on_name",
							Fields: []*config.IndexField{
								{
									Column: "name",
								},
							},
						},
						{
							Name: "index_on_name_age",
							Fields: []*config.IndexField{
								{
									Column: "age",
								},
								{
									Column: "name",
									Order:  "ASC",
								},
							},
						},
						{
							Name: "index_on_dashboard",
							Fields: []*config.IndexField{
								{
									Column: "name",
								},
								{
									Column: "email",
								},
							},
						},
						{
							Name: "index_on_location",
							Fields: []*config.IndexField{
								{
									Column: "location",
								},
							},
						},
					},
				},
			},
			target: []*config.Schema{
				{
					Name: "users",
					Index: []*config.Index{
						{
							Name: "index_on_name",
							Fields: []*config.IndexField{
								{
									Column: "name",
								},
							},
						},
						{
							Name: "index_on_name_age",
							Fields: []*config.IndexField{
								{
									Column: "name",
									Order:  "ASC",
								},
								{
									Column: "age",
								},
							},
						},
						{
							Name: "index_on_dashboard",
							Fields: []*config.IndexField{
								{
									Column: "name",
								},
								{
									Column: "email",
								},
								{
									Column: "location",
								},
							},
						},
						{
							Name: "index_on_email",
							Fields: []*config.IndexField{
								{
									Column: "email",
								},
							},
						},
					},
				},
			},
			input: "users",
			result: &step.AlterSchema{
				Name: "users",
				DroppedIndices: []*config.Index{
					{
						Name: "index_on_location",
						Fields: []*config.IndexField{
							{
								Column: "location",
							},
						},
					},
					{
						Name: "index_on_dashboard",
						Fields: []*config.IndexField{
							{
								Column: "name",
							},
							{
								Column: "email",
							},
						},
					},
				},
				AddedIndices: []*config.Index{
					{
						Name: "index_on_email",
						Fields: []*config.IndexField{
							{
								Column: "email",
							},
						},
					},
					{
						Name: "index_on_dashboard",
						Fields: []*config.IndexField{
							{
								Column: "name",
							},
							{
								Column: "email",
							},
							{
								Column: "location",
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		diffSchema := diff.NewSchema(tc.existing, tc.target)
		result, err := diffSchema.AlteredSchema(tc.input)
		assert.Equal(t, tc.err, err)
		if tc.result != nil {
			assert.Equal(t, tc.result.Name, result.Name)
			assert.ElementsMatch(t, tc.result.AddedColumns, result.AddedColumns)
			assert.ElementsMatch(t, tc.result.AlteredColumns, result.AlteredColumns)
			assert.ElementsMatch(t, tc.result.DroppedColumns, result.DroppedColumns)
		}
	}
}

func TestGeneratePlan(t *testing.T) {
	testCases := []struct {
		existingSchema []*config.Schema
		targetSchema   []*config.Schema
		result         *step.MigrationPlanner
		err            error
	}{
		{
			existingSchema: []*config.Schema{
				{
					Name: "users",
					Fields: []*config.Field{
						{
							Name: "id",
							Type: "bigserial",
							Options: []field_option.FieldOption{
								field_option.NotNull,
							},
						},
						{
							Name:  "name",
							Type:  "varchar",
							Limit: 20,
						},
						{
							Name:  "address",
							Type:  "varchar",
							Limit: 50,
						},
						{
							Name:  "school",
							Type:  "varchar",
							Limit: 50,
						},
						{
							Name: "age",
							Type: "int",
							Options: []field_option.FieldOption{
								field_option.NotNull,
							},
						},
					},
				},
				{
					Name: "example",
					Fields: []*config.Field{
						{
							Name: "id",
							Type: "bigserial",
						},
						{
							Name:  "name",
							Type:  "varchar",
							Limit: 20,
						},
					},
				},
			},
			targetSchema: []*config.Schema{
				{
					Name: "users",
					Fields: []*config.Field{
						{
							Name: "id",
							Type: "bigserial",
							Options: []field_option.FieldOption{
								field_option.NotNull,
							},
						},
						{
							Name:  "name",
							Type:  "varchar",
							Limit: 20,
							Options: []field_option.FieldOption{
								field_option.NotNull,
							},
						},
						{
							Name:  "address",
							Type:  "varchar",
							Limit: 150,
						},
						{
							Name:  "school",
							Type:  "varchar",
							Limit: 100,
							Options: []field_option.FieldOption{
								field_option.NotNull,
							},
						},
						{
							Name: "created_at",
							Type: "timestamp",
							Options: []field_option.FieldOption{
								field_option.NotNull,
							},
						},
					},
				},
				{
					Name: "documents",
					Fields: []*config.Field{
						{
							Name: "id",
							Type: "bigserial",
						},
						{
							Name:  "name",
							Type:  "varchar",
							Limit: 20,
						},
					},
				},
			},
			result: &step.MigrationPlanner{
				CreateTable: []*config.Schema{
					{
						Name: "documents",
						Fields: []*config.Field{
							{
								Name: "id",
								Type: "bigserial",
							},
							{
								Name:  "name",
								Type:  "varchar",
								Limit: 20,
							},
						},
					},
				},
				DropTable: []*config.Schema{
					{
						Name: "example",
						Fields: []*config.Field{
							{
								Name: "id",
								Type: "bigserial",
							},
							{
								Name:  "name",
								Type:  "varchar",
								Limit: 20,
							},
						},
					},
				},
				AlterSchema: map[string]*step.AlterSchema{
					"users": {
						Name: "users",
						AddedColumns: []*config.Field{
							{
								Name: "created_at",
								Type: "timestamp",
								Options: []field_option.FieldOption{
									field_option.NotNull,
								},
							},
						},
						AlteredColumns: []*step.AlterColumn{
							{
								Name: "name",
								Field: &config.Field{
									Name:  "name",
									Type:  "varchar",
									Limit: 20,
									Options: []field_option.FieldOption{
										field_option.NotNull,
									},
								},
								LastField: &config.Field{
									Name:  "name",
									Type:  "varchar",
									Limit: 20,
								},
								ChangedOptions: []step.OptionAction{step.SetNotNull},
							},
							{
								Name: "address",
								Field: &config.Field{
									Name:  "address",
									Type:  "varchar",
									Limit: 150,
								},
								LastField: &config.Field{
									Name:  "address",
									Type:  "varchar",
									Limit: 50,
								},
								ChangedType:    true,
								ChangedOptions: []step.OptionAction{},
							},
							{
								Name: "school",
								Field: &config.Field{
									Name:  "school",
									Type:  "varchar",
									Limit: 100,
									Options: []field_option.FieldOption{
										field_option.NotNull,
									},
								},
								LastField: &config.Field{
									Name:  "school",
									Type:  "varchar",
									Limit: 50,
								},
								ChangedType:    true,
								ChangedOptions: []step.OptionAction{step.SetNotNull},
							},
						},
						DroppedColumns: []*config.Field{
							{
								Name: "age",
								Type: "int",
								Options: []field_option.FieldOption{
									field_option.NotNull,
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		schemaDiff := diff.NewSchema(tc.existingSchema, tc.targetSchema)
		result, err := schemaDiff.GeneratePlan()
		assert.Equal(t, tc.err, err)
		assert.ElementsMatch(t, tc.result.CreateTable, result.CreateTable)
		assert.ElementsMatch(t, tc.result.DropTable, result.DropTable)
		for name, alter := range tc.result.AlterSchema {
			res := result.AlterSchema[name]
			assert.Equal(t, alter.Name, res.Name)
			assert.ElementsMatch(t, alter.AddedColumns, res.AddedColumns)
			assert.ElementsMatch(t, alter.DroppedColumns, res.DroppedColumns)
			assert.ElementsMatch(t, alter.AlteredColumns, res.AlteredColumns)
		}
	}
}
