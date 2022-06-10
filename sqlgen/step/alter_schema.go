package step

import (
	"gitlab.com/wartek-id/core/tools/dbgen/config"
)

type AlterSchema struct {
	Name           string
	AddedColumns   []*config.Field
	AlteredColumns []*AlterColumn
	DroppedColumns []*config.Field

	AddedIndices   []*config.Index
	DroppedIndices []*config.Index
}

func NewAlterSchema(name string) *AlterSchema {
	return &AlterSchema{
		Name:           name,
		AddedColumns:   make([]*config.Field, 0),
		AlteredColumns: make([]*AlterColumn, 0),
		DroppedColumns: make([]*config.Field, 0),
	}
}

func (s *AlterSchema) HasChanges() bool {
	return s.FieldChanged() || s.IndicesChanged()
}

func (s *AlterSchema) FieldChanged() bool {
	return s.IsColumnsAdded() ||
		s.IsColumnsAltered() ||
		s.IsColumnsDropped()
}

func (s *AlterSchema) IndicesChanged() bool {
	return s.IsIndicesAdded() ||
		s.IsIndicesDropped()
}

func (s *AlterSchema) IsColumnsAdded() bool {
	return len(s.AddedColumns) != 0
}

func (s *AlterSchema) IsColumnsAltered() bool {
	return len(s.AlteredColumns) != 0
}

func (s *AlterSchema) IsColumnsDropped() bool {
	return len(s.DroppedColumns) != 0
}

func (s *AlterSchema) IsIndicesAdded() bool {
	return len(s.AddedIndices) != 0
}

func (s *AlterSchema) IsIndicesDropped() bool {
	return len(s.DroppedIndices) != 0
}
