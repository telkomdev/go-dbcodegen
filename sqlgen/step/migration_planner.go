package step

import "gitlab.com/wartek-id/core/tools/dbgen/config"

type MigrationPlanner struct {
	CreateTable []*config.Schema
	DropTable   []*config.Schema
	AlterSchema map[string]*AlterSchema
}

func NewMigrationPlanner() *MigrationPlanner {
	return &MigrationPlanner{
		CreateTable: make([]*config.Schema, 0),
		DropTable:   make([]*config.Schema, 0),
		AlterSchema: make(map[string]*AlterSchema),
	}
}
