package sqlgen_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/wartek-id/core/tools/dbgen/config"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen"
	mock_schema "gitlab.com/wartek-id/core/tools/dbgen/sqlgen/mocks/schema"
	"gitlab.com/wartek-id/core/tools/dbgen/types/field_option"
)

func TestSqlGenerator_Generate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	target := filepath.Join(t.TempDir(), "generator")
	upTarget := fmt.Sprintf("%s.up.sql", target)
	downTarget := fmt.Sprintf("%s.down.sql", target)
	mockCrawler := mock_schema.NewMockSchema(ctrl)
	mockCrawler.EXPECT().GetSchemas().Return([]*config.Schema{
		{
			Name: "documents",
			Fields: []*config.Field{
				{
					Name: "id",
					Type: "bigserial",
					Options: []field_option.FieldOption{
						field_option.PrimaryKey,
					},
				},
				{
					Name:  "name",
					Type:  "varchar",
					Limit: 200,
				},
				{
					Name: "release_date",
					Type: "timestamp",
				},
				{
					Name: "created_at",
					Type: "timestamp",
				},
			},
			Index: []*config.Index{
				{
					Name: "index_document_on_name",
					Fields: []*config.IndexField{
						{
							Column: "name",
						},
					},
					Unique: true,
				},
			},
		},
		{
			Name: "example",
			Fields: []*config.Field{
				{
					Name: "id",
					Type: "bigserial",
					Options: []field_option.FieldOption{
						field_option.PrimaryKey,
					},
				},
				{
					Name:  "name",
					Type:  "varchar",
					Limit: 100,
				},
			},
		},
	}, nil).AnyTimes()

	gen := sqlgen.NewGenerator(mockCrawler, []*config.Schema{
		{
			Name: "user",
			Fields: []*config.Field{
				{
					Name: "id",
					Type: "bigserial",
				},
				{
					Name:  "name",
					Type:  "varchar",
					Limit: 200,
				},
				{
					Name: "created_at",
					Type: "timestamp",
				},
			},
			Index: []*config.Index{
				{
					Name: "index_user_on_name",
					Fields: []*config.IndexField{
						{
							Column: "name",
						},
					},
					Unique: true,
				},
			},
		},
		{
			Name: "documents",
			Fields: []*config.Field{
				{
					Name: "id",
					Type: "bigserial",
					Options: []field_option.FieldOption{
						field_option.PrimaryKey,
					},
				},
				{
					Name:  "name",
					Type:  "varchar",
					Limit: 50,
					Options: []field_option.FieldOption{
						field_option.NotNull,
					},
				},
				{
					Name:  "approval",
					Type:  "varchar",
					Limit: 50,
				},
				{
					Name: "created_at",
					Type: "timestamp",
					Options: []field_option.FieldOption{
						field_option.NotNull,
					},
				},
				{
					Name: "updated_at",
					Type: "timestamp",
					Options: []field_option.FieldOption{
						field_option.NotNull,
					},
				},
			},
			Index: []*config.Index{
				{
					Name: "index_document_on_name_approval",
					Fields: []*config.IndexField{
						{
							Column: "name",
							Order:  "ASC",
						},
						{
							Column: "approval",
							Order:  "ASC",
						},
					},
					Unique: true,
				},
			},
		},
	}, &sqlgen.Flag{OutputTarget: target, SkipDropTable: false})
	err := gen.Generate()
	assert.NoError(t, err)

	upMigration, err := os.ReadFile(upTarget)
	comLen := strings.Split(string(upMigration), ";")
	assert.NoError(t, err)
	assert.Len(t, comLen, 9)
	assert.Contains(t, string(upMigration), "BEGIN;")
	assert.Contains(t, string(upMigration), "CREATE TABLE IF NOT EXISTS \"user\" (\n\t\"id\" BIGSERIAL,\n\t\"name\" VARCHAR(200),\n\t\"created_at\" TIMESTAMP\n);")
	assert.Contains(t, string(upMigration), "CREATE UNIQUE INDEX IF NOT EXISTS \"index_user_on_name\" ON \"user\"(\"name\");")
	assert.Contains(t, string(upMigration), "ALTER TABLE IF EXISTS \"documents\"\n")
	assert.Contains(t, string(upMigration), "\tADD COLUMN \"updated_at\" TIMESTAMP NOT NULL")
	assert.Contains(t, string(upMigration), "\tADD COLUMN \"approval\" VARCHAR(50)")
	assert.Contains(t, string(upMigration), "\tDROP COLUMN \"release_date\"")
	assert.Contains(t, string(upMigration), "\tALTER COLUMN \"created_at\" SET NOT NULL")
	assert.Contains(t, string(upMigration), "\tALTER COLUMN \"name\" SET DATA TYPE VARCHAR(50)")
	assert.Contains(t, string(upMigration), "\tALTER COLUMN \"name\" SET NOT NULL")
	assert.Contains(t, string(upMigration), "DROP INDEX IF EXISTS \"index_document_on_name\";")
	assert.Contains(t, string(upMigration), "CREATE UNIQUE INDEX IF NOT EXISTS \"index_document_on_name_approval\" ON \"documents\"(\"name\" ASC, \"approval\" ASC);")
	assert.Contains(t, string(upMigration), "DROP TABLE IF EXISTS \"example\";")
	assert.Contains(t, string(upMigration), "COMMIT;")

	downMigration, err := os.ReadFile(downTarget)
	comLen = strings.Split(string(downMigration), ";")
	assert.NoError(t, err)
	assert.Len(t, comLen, 8)
	assert.Contains(t, string(downMigration), "BEGIN;")
	assert.Contains(t, string(downMigration), "DROP TABLE IF EXISTS \"user\";")
	assert.Contains(t, string(downMigration), "ALTER TABLE IF EXISTS \"documents\"\n")
	assert.Contains(t, string(downMigration), "\tDROP COLUMN \"approval\"")
	assert.Contains(t, string(downMigration), "\tDROP COLUMN \"updated_at\"")
	assert.Contains(t, string(downMigration), "\tADD COLUMN \"release_date\" TIMESTAMP")
	assert.Contains(t, string(downMigration), "\tALTER COLUMN \"name\" DROP NOT NULL")
	assert.Contains(t, string(downMigration), "\tALTER COLUMN \"name\" SET DATA TYPE VARCHAR(200)")
	assert.Contains(t, string(downMigration), "\tALTER COLUMN \"created_at\" DROP NOT NULL")
	assert.Contains(t, string(downMigration), "CREATE UNIQUE INDEX IF NOT EXISTS \"index_document_on_name\" ON \"documents\"(\"name\");")
	assert.Contains(t, string(downMigration), "DROP INDEX IF EXISTS \"index_document_on_name_approval\";")
	assert.Contains(t, string(downMigration), "CREATE TABLE IF NOT EXISTS \"example\" (\n\t\"id\" BIGSERIAL PRIMARY KEY,\n\t\"name\" VARCHAR(100)\n);")
	assert.Contains(t, string(downMigration), "COMMIT;")
}

func TestSqlGenerator_GenerateNoDrop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	target := filepath.Join(t.TempDir(), "generator")
	upTarget := fmt.Sprintf("%s.up.sql", target)
	downTarget := fmt.Sprintf("%s.down.sql", target)
	mockCrawler := mock_schema.NewMockSchema(ctrl)
	mockCrawler.EXPECT().GetSchemas().Return([]*config.Schema{
		{
			Name: "documents",
			Fields: []*config.Field{
				{
					Name: "id",
					Type: "bigserial",
					Options: []field_option.FieldOption{
						field_option.PrimaryKey,
					},
				},
				{
					Name:  "name",
					Type:  "varchar",
					Limit: 200,
				},
				{
					Name: "release_date",
					Type: "timestamp",
				},
				{
					Name: "created_at",
					Type: "timestamp",
				},
			},
			Index: []*config.Index{
				{
					Name: "index_document_on_name",
					Fields: []*config.IndexField{
						{
							Column: "name",
						},
					},
					Unique: true,
				},
			},
		},
		{
			Name: "example",
			Fields: []*config.Field{
				{
					Name: "id",
					Type: "bigserial",
					Options: []field_option.FieldOption{
						field_option.PrimaryKey,
					},
				},
				{
					Name:  "name",
					Type:  "varchar",
					Limit: 100,
				},
			},
		},
	}, nil).AnyTimes()

	gen := sqlgen.NewGenerator(mockCrawler, []*config.Schema{
		{
			Name: "user",
			Fields: []*config.Field{
				{
					Name: "id",
					Type: "bigserial",
				},
				{
					Name:  "name",
					Type:  "varchar",
					Limit: 200,
				},
				{
					Name: "created_at",
					Type: "timestamp",
				},
			},
			Index: []*config.Index{
				{
					Name: "index_user_on_name",
					Fields: []*config.IndexField{
						{
							Column: "name",
						},
					},
					Unique: true,
				},
			},
		},
		{
			Name: "documents",
			Fields: []*config.Field{
				{
					Name: "id",
					Type: "bigserial",
					Options: []field_option.FieldOption{
						field_option.PrimaryKey,
					},
				},
				{
					Name:  "name",
					Type:  "varchar",
					Limit: 50,
					Options: []field_option.FieldOption{
						field_option.NotNull,
					},
				},
				{
					Name:  "approval",
					Type:  "varchar",
					Limit: 50,
				},
				{
					Name: "created_at",
					Type: "timestamp",
					Options: []field_option.FieldOption{
						field_option.NotNull,
					},
				},
				{
					Name: "updated_at",
					Type: "timestamp",
					Options: []field_option.FieldOption{
						field_option.NotNull,
					},
				},
			},
			Index: []*config.Index{
				{
					Name: "index_document_on_name_approval",
					Fields: []*config.IndexField{
						{
							Column: "name",
							Order:  "ASC",
						},
						{
							Column: "approval",
							Order:  "ASC",
						},
					},
					Unique: true,
				},
			},
		},
	}, &sqlgen.Flag{OutputTarget: target, SkipDropTable: true})
	err := gen.Generate()
	assert.NoError(t, err)

	upMigration, err := os.ReadFile(upTarget)
	comLen := strings.Split(string(upMigration), ";")
	assert.NoError(t, err)
	assert.Len(t, comLen, 8)
	assert.Contains(t, string(upMigration), "BEGIN;")
	assert.Contains(t, string(upMigration), "CREATE TABLE IF NOT EXISTS \"user\" (\n\t\"id\" BIGSERIAL,\n\t\"name\" VARCHAR(200),\n\t\"created_at\" TIMESTAMP\n);")
	assert.Contains(t, string(upMigration), "CREATE UNIQUE INDEX IF NOT EXISTS \"index_user_on_name\" ON \"user\"(\"name\");")
	assert.Contains(t, string(upMigration), "ALTER TABLE IF EXISTS \"documents\"\n")
	assert.Contains(t, string(upMigration), "\tADD COLUMN \"updated_at\" TIMESTAMP NOT NULL")
	assert.Contains(t, string(upMigration), "\tADD COLUMN \"approval\" VARCHAR(50)")
	assert.Contains(t, string(upMigration), "\tDROP COLUMN \"release_date\"")
	assert.Contains(t, string(upMigration), "\tALTER COLUMN \"created_at\" SET NOT NULL")
	assert.Contains(t, string(upMigration), "\tALTER COLUMN \"name\" SET DATA TYPE VARCHAR(50)")
	assert.Contains(t, string(upMigration), "\tALTER COLUMN \"name\" SET NOT NULL")
	assert.Contains(t, string(upMigration), "DROP INDEX IF EXISTS \"index_document_on_name\";")
	assert.Contains(t, string(upMigration), "CREATE UNIQUE INDEX IF NOT EXISTS \"index_document_on_name_approval\" ON \"documents\"(\"name\" ASC, \"approval\" ASC);")
	assert.NotContains(t, string(upMigration), "DROP TABLE IF EXISTS \"example\";")
	assert.Contains(t, string(upMigration), "COMMIT;")

	downMigration, err := os.ReadFile(downTarget)
	comLen = strings.Split(string(downMigration), ";")
	assert.NoError(t, err)
	assert.Len(t, comLen, 7)
	assert.Contains(t, string(downMigration), "BEGIN;")
	assert.Contains(t, string(downMigration), "DROP TABLE IF EXISTS \"user\";")
	assert.Contains(t, string(downMigration), "ALTER TABLE IF EXISTS \"documents\"\n")
	assert.Contains(t, string(downMigration), "\tDROP COLUMN \"approval\"")
	assert.Contains(t, string(downMigration), "\tDROP COLUMN \"updated_at\"")
	assert.Contains(t, string(downMigration), "\tADD COLUMN \"release_date\" TIMESTAMP")
	assert.Contains(t, string(downMigration), "\tALTER COLUMN \"name\" DROP NOT NULL")
	assert.Contains(t, string(downMigration), "\tALTER COLUMN \"name\" SET DATA TYPE VARCHAR(200)")
	assert.Contains(t, string(downMigration), "\tALTER COLUMN \"created_at\" DROP NOT NULL")
	assert.Contains(t, string(downMigration), "CREATE UNIQUE INDEX IF NOT EXISTS \"index_document_on_name\" ON \"documents\"(\"name\");")
	assert.Contains(t, string(downMigration), "DROP INDEX IF EXISTS \"index_document_on_name_approval\";")
	assert.NotContains(t, string(downMigration), "CREATE TABLE IF NOT EXISTS \"example\" (\n\t\"id\" BIGSERIAL PRIMARY KEY,\n\t\"name\" VARCHAR(100)\n);")
	assert.Contains(t, string(downMigration), "COMMIT;")
}

func TestSqlGenerator_GenerateNoChanges(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	target := filepath.Join(t.TempDir(), "generator")
	upTarget := fmt.Sprintf("%s.up.sql", target)
	downTarget := fmt.Sprintf("%s.down.sql", target)
	mockCrawler := mock_schema.NewMockSchema(ctrl)
	mockCrawler.EXPECT().GetSchemas().Return([]*config.Schema{
		{
			Name: "documents",
			Fields: []*config.Field{
				{
					Name: "id",
					Type: "bigserial",
					Options: []field_option.FieldOption{
						field_option.PrimaryKey,
					},
				},
				{
					Name:  "name",
					Type:  "varchar",
					Limit: 200,
				},
				{
					Name: "release_date",
					Type: "timestamp",
				},
				{
					Name: "created_at",
					Type: "timestamp",
				},
			},
			Index: []*config.Index{
				{
					Name: "index_document_on_name",
					Fields: []*config.IndexField{
						{
							Column: "name",
						},
					},
					Unique: true,
				},
			},
		},
		{
			Name: "example",
			Fields: []*config.Field{
				{
					Name: "id",
					Type: "bigserial",
					Options: []field_option.FieldOption{
						field_option.PrimaryKey,
					},
				},
				{
					Name:  "name",
					Type:  "varchar",
					Limit: 100,
				},
			},
		},
	}, nil).AnyTimes()

	gen := sqlgen.NewGenerator(mockCrawler, []*config.Schema{
		{
			Name: "documents",
			Fields: []*config.Field{
				{
					Name: "id",
					Type: "bigserial",
					Options: []field_option.FieldOption{
						field_option.PrimaryKey,
					},
				},
				{
					Name:  "name",
					Type:  "varchar",
					Limit: 200,
				},
				{
					Name: "release_date",
					Type: "timestamp",
				},
				{
					Name: "created_at",
					Type: "timestamp",
				},
			},
			Index: []*config.Index{
				{
					Name: "index_document_on_name",
					Fields: []*config.IndexField{
						{
							Column: "name",
						},
					},
					Unique: true,
				},
			},
		},
		{
			Name: "example",
			Fields: []*config.Field{
				{
					Name: "id",
					Type: "bigserial",
					Options: []field_option.FieldOption{
						field_option.PrimaryKey,
					},
				},
				{
					Name:  "name",
					Type:  "varchar",
					Limit: 100,
				},
			},
		},
	}, &sqlgen.Flag{OutputTarget: target, SkipDropTable: false})
	err := gen.Generate()
	assert.NoError(t, err)

	upMig, err := os.ReadFile(upTarget)
	assert.Error(t, err, "not found")
	assert.Nil(t, upMig)

	downMig, err := os.ReadFile(downTarget)
	assert.Error(t, err, "not found")
	assert.Nil(t, downMig)
}

func TestSqlGenerator_CreateTableGenerator(t *testing.T) {
	gen := sqlgen.NewGenerator(nil, []*config.Schema{}, &sqlgen.Flag{OutputTarget: "target"})
	assert.NotNil(t, gen.CreateTableGenerator())
}

func TestSqlGenerator_CreateIndexGenerator(t *testing.T) {
	gen := sqlgen.NewGenerator(nil, []*config.Schema{}, &sqlgen.Flag{OutputTarget: "target"})
	assert.NotNil(t, gen.CreateIndexGenerator())
}

func TestSqlGenerator_AlterTableGenerator(t *testing.T) {
	gen := sqlgen.NewGenerator(nil, []*config.Schema{}, &sqlgen.Flag{OutputTarget: "target"})
	assert.NotNil(t, gen.AlterTableGenerator())
}

func TestSqlGenerator_DropIndexGenerator(t *testing.T) {
	gen := sqlgen.NewGenerator(nil, []*config.Schema{}, &sqlgen.Flag{OutputTarget: "target"})
	assert.NotNil(t, gen.DropIndexGenerator())
}

func TestSqlGenerator_DropTableGenerator(t *testing.T) {
	gen := sqlgen.NewGenerator(nil, []*config.Schema{}, &sqlgen.Flag{OutputTarget: "target"})
	assert.NotNil(t, gen.DropTableGenerator())
}
