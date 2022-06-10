package step

import (
	"gitlab.com/wartek-id/core/tools/dbgen/config"
)

type OptionAction int

const (
	DropNotNull OptionAction = iota
	SetNotNull
)

type AlterColumn struct {
	Name                string
	Field               *config.Field
	LastField           *config.Field
	ChangedType         bool
	ChangedDefaultValue bool
	ChangedOptions      []OptionAction
}

func (c *AlterColumn) HasChanges() bool {
	return c.ChangedType || c.IsOptionsChanged()
}

func (c *AlterColumn) IsOptionsChanged() bool {
	return len(c.ChangedOptions) != 0
}
