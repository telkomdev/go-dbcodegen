package diff

import (
	"errors"

	"github.com/google/go-cmp/cmp"
	"gitlab.com/wartek-id/core/tools/dbgen/config"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/step"
	"gitlab.com/wartek-id/core/tools/dbgen/types/field_type"
)

var (
	ErrMissingCurrentTable = errors.New("current table is not exists")
	ErrMissingTargetTable  = errors.New("missing target table")
)

type Nameable interface {
	GetName() string
}

type diffSchema struct {
	name    string
	schema  *config.Schema
	fields  map[string]*config.Field
	indexes map[string]*config.Index
}

type Schema struct {
	from   map[string]*diffSchema
	target map[string]*diffSchema
}

func NewSchema(from, target []*config.Schema) *Schema {
	return &Schema{
		from:   buildSchema(from),
		target: buildSchema(target),
	}
}

func (diff *Schema) GeneratePlan() (*step.MigrationPlanner, error) {
	planner := step.NewMigrationPlanner()
	planner.CreateTable = diff.CreatedTable()
	planner.DropTable = diff.DroppedTable()

	for name := range diff.target {
		existingTable := diff.from[name]
		if existingTable == nil {
			continue
		}
		alterTable, err := diff.AlteredSchema(name)
		if err != nil {
			return nil, err
		}

		if alterTable.HasChanges() {
			planner.AlterSchema[name] = alterTable
		}
	}

	return planner, nil
}

func (diff *Schema) CreatedTable() []*config.Schema {
	createdTable := make([]*config.Schema, 0)
	for name, diffSchema := range diff.target {
		if diff.from[name] == nil {
			createdTable = append(createdTable, diffSchema.schema)
		}
	}
	return createdTable
}

func (diff *Schema) DroppedTable() []*config.Schema {
	droppedTable := make([]*config.Schema, 0)
	for name, diffSchema := range diff.from {
		if diff.target[name] == nil {
			droppedTable = append(droppedTable, diffSchema.schema)
		}
	}
	return droppedTable
}

func (diff *Schema) AlteredIndexes(existing, target map[string]*config.Index, planner *step.AlterSchema) {
	for name, index := range existing {
		if target[name] == nil {
			planner.DroppedIndices = append(planner.DroppedIndices, index)
		}
	}

	for name, targetIndex := range target {
		existingIndex := existing[name]
		if existingIndex == nil {
			planner.AddedIndices = append(planner.AddedIndices, targetIndex)
			continue
		}

		if !cmp.Equal(existingIndex, targetIndex) {
			planner.DroppedIndices = append(planner.DroppedIndices, existingIndex)
			planner.AddedIndices = append(planner.AddedIndices, targetIndex)
		}
	}
}

func (diff *Schema) AlteredSchema(table string) (*step.AlterSchema, error) {
	tableFrom := diff.from[table]
	if tableFrom == nil {
		return nil, ErrMissingCurrentTable
	}
	tableTarget := diff.target[table]
	if tableTarget == nil {
		return nil, ErrMissingTargetTable
	}

	existingFields := tableFrom.fields
	targetFields := tableTarget.fields
	migrationSteps := step.NewAlterSchema(table)

	for name, field := range targetFields {
		existingField := existingFields[name]
		if existingField == nil {
			migrationSteps.AddedColumns = append(migrationSteps.AddedColumns, field)
			continue
		}

		alteredColumn := diff.alteredColumn(existingField, field)
		if alteredColumn.HasChanges() {
			migrationSteps.AlteredColumns = append(migrationSteps.AlteredColumns, alteredColumn)
		}
	}

	for name, field := range existingFields {
		if targetFields[name] == nil {
			migrationSteps.DroppedColumns = append(migrationSteps.DroppedColumns, field)
			continue
		}
	}

	diff.AlteredIndexes(tableFrom.indexes, tableTarget.indexes, migrationSteps)
	return migrationSteps, nil
}

func (diff *Schema) alteredColumn(from, target *config.Field) *step.AlterColumn {
	alterColumn := step.AlterColumn{
		Name:                target.Name,
		Field:               target,
		LastField:           from,
		ChangedType:         !diff.isSameFieldType(from, target),
		ChangedDefaultValue: diff.changedDefaultValue(from, target),
		ChangedOptions:      diff.changedOptions(from, target),
	}

	return &alterColumn
}

func (diff *Schema) isSameFieldType(from, target *config.Field) bool {
	if from.Type != target.Type {
		return false
	}

	switch target.Type {
	case field_type.Varchar, field_type.Decimal:
		return from.Limit == target.Limit &&
			from.Scale == target.Scale
	}

	return true
}

func (diff *Schema) changedOptions(from, target *config.Field) []step.OptionAction {
	options := make([]step.OptionAction, 0)
	if from.IsNotNull() != target.IsNotNull() {
		if target.IsNotNull() {
			options = append(options, step.SetNotNull)
		} else {
			options = append(options, step.DropNotNull)
		}
	}

	return options
}

func (diff *Schema) changedDefaultValue(from, target *config.Field) bool {
	return from.Default != target.Default
}

func buildSchema(schemas []*config.Schema) map[string]*diffSchema {
	cmpSchema := make(map[string]*diffSchema)
	for _, sc := range schemas {
		cmpSchema[sc.Name] = &diffSchema{
			name:    sc.Name,
			schema:  sc,
			fields:  nameableMapper(sc.Fields),
			indexes: nameableMapper(sc.Index),
		}
	}
	return cmpSchema
}

func nameableMapper[T Nameable](elements []T) map[string]T {
	result := make(map[string]T)
	for _, element := range elements {
		result[element.GetName()] = element
	}

	return result
}
