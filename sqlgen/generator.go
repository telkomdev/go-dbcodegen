package sqlgen

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"gitlab.com/wartek-id/core/tools/dbgen/config"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/dialect"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/diff"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/sb"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/schema"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/step"
)

const (
	DefaultMigrationExt         = ".sql"
	DefaultDialect              = "postgres"
	FullSchemaMigrationFilename = "temp/fullschema/migration.sql"
)

var SectionSeparator = []byte("\n\n")

type SqlGenerator struct {
	dbUpFilename   string
	dbDownFilename string
	flag           *Flag
	generators     *generators
	schemas        []*config.Schema
	crawler        schema.Schema
	dialect        string
	dialectOption  *dialect.DialectOption
}

type generators struct {
	ctg CreateTableGenerator
	cig CreateIndexGenerator
	atg AlterTableGenerator
	dig DropIndexGenerator
	dtg DropTableGenerator
}

func NewGenerator(crawler schema.Schema, schemas []*config.Schema, flag *Flag) *SqlGenerator {
	dbUpFilename, dbDownFilename := getTargetPath(flag.OutputDirectory, flag.OutputTarget)

	return &SqlGenerator{
		dbUpFilename:   dbUpFilename,
		dbDownFilename: dbDownFilename,
		schemas:        schemas,
		flag:           flag,
		generators:     initGenerator(DefaultDialect, dialect.DefaultDialectOption()),
		crawler:        crawler,
		dialect:        DefaultDialect,
		dialectOption:  dialect.DefaultDialectOption(),
	}
}

func initGenerator(dialect string, do *dialect.DialectOption) *generators {
	return &generators{
		ctg: NewCreateTableGenerator(dialect, do),
		cig: NewCreateIndexGenerator(dialect, do),
		atg: NewAlterTableGenerator(dialect, do),
		dig: NewDropIndexGenerator(dialect, do),
		dtg: NewDropTableGenerator(dialect, do),
	}
}

func getTargetPath(dir, target string) (string, string) {
	ext := filepath.Ext(target)
	if ext != "" {
		target = strings.Replace(target, ext, "", -1)
	} else {
		ext = DefaultMigrationExt
	}

	upFilename := fmt.Sprintf("%s/%s.up%s", dir, target, ext)
	downFilename := fmt.Sprintf("%s/%s.down%s", dir, target, ext)
	return upFilename, downFilename
}

func (gen *SqlGenerator) CreateTableGenerator() CreateTableGenerator {
	return gen.generators.ctg
}

func (gen *SqlGenerator) CreateIndexGenerator() CreateIndexGenerator {
	return gen.generators.cig
}

func (gen *SqlGenerator) AlterTableGenerator() AlterTableGenerator {
	return gen.generators.atg
}

func (gen *SqlGenerator) DropIndexGenerator() DropIndexGenerator {
	return gen.generators.dig
}

func (gen *SqlGenerator) DropTableGenerator() DropTableGenerator {
	return gen.generators.dtg
}

func (gen *SqlGenerator) Generate() error {
	currentSchemas, err := gen.crawler.GetSchemas()
	if err != nil {
		return err
	}

	planner := diff.NewSchema(currentSchemas, gen.schemas)
	migrationPlanner, err := planner.GeneratePlan()
	if err != nil {
		return err
	}

	err = gen.UpMigration(migrationPlanner)
	if err != nil {
		return err
	}

	fmt.Println()
	err = gen.DownMigration(migrationPlanner)
	if err != nil {
		return err
	}

	fmt.Println()
	err = gen.FullSchemaMigration()
	if err != nil {
		return err
	}

	fmt.Println("\nDatabase Migration Generation Completed.")
	return nil
}

func (gen *SqlGenerator) FullSchemaMigration() error {
	fmt.Println("🚀 Generating up full schema migration file")

	createTables := gen.GenerateCreateTables(gen.schemas)
	err := gen.Writer(FullSchemaMigrationFilename, getContents(createTables))
	if err != nil {
		fmt.Println(color.RedString("Failed"))
		return err
	}
	fmt.Println(color.GreenString("Succeeded"))
	return nil
}

func (gen *SqlGenerator) UpMigration(plan *step.MigrationPlanner) error {
	fmt.Println("🚀 Generating up database migration files")
	fmt.Printf("Target file: %s\n", color.HiBlueString(gen.dbUpFilename))

	createTables := gen.GenerateCreateTables(plan.CreateTable)
	alterTables := gen.AlterTableUp(plan.AlterSchema)
	dropTables := []byte{}
	if !gen.flag.SkipDropTable {
		dropTables = gen.GenerateDropTables(plan.DropTable)
	}

	content := getContents(createTables, dropTables, alterTables)
	if len(bytes.TrimSpace(content)) == 0 {
		fmt.Println(color.YellowString("No changes being detected, skipping..."))
		return nil
	}

	if gen.dialectOption.SupportTransaction {
		content = getContents(
			gen.dialectOption.BeginClause, content, gen.dialectOption.CommitClause,
		)
	}
	err := gen.Writer(gen.dbUpFilename, content)
	if err != nil {
		fmt.Println(color.RedString("Failed"))
		return err
	}
	fmt.Println(color.GreenString("Succeeded"))
	return nil
}

func (gen *SqlGenerator) AlterTableUp(alterSchemas map[string]*step.AlterSchema) []byte {
	contents := make([][]byte, 0)
	for _, as := range alterSchemas {
		diBuf := sb.NewSQLBuilder()
		for _, idx := range as.DroppedIndices {
			gen.DropIndexGenerator().Generate(diBuf, idx)
			diBuf.WriteNewLine()
		}

		atBuf := sb.NewSQLBuilder()
		gen.AlterTableGenerator().Generate(atBuf, as)

		aiBuf := sb.NewSQLBuilder()
		for _, idx := range as.AddedIndices {
			gen.CreateIndexGenerator().Generate(aiBuf, as.Name, idx)
			aiBuf.WriteNewLine()
		}

		contents = append(contents, getContents(atBuf.Bytes(), diBuf.Bytes(), aiBuf.Bytes()))
	}
	return bytes.Join(contents, SectionSeparator)
}

func (gen *SqlGenerator) AlterTableDown(alterSchemas map[string]*step.AlterSchema) []byte {
	contents := make([][]byte, 0)
	for _, as := range alterSchemas {
		aiBuf := sb.NewSQLBuilder()
		for _, idx := range as.AddedIndices {
			gen.DropIndexGenerator().Generate(aiBuf, idx)
			aiBuf.WriteNewLine()
		}

		atBuf := sb.NewSQLBuilder()
		gen.AlterTableGenerator().Rollback(atBuf, as)

		diBuf := sb.NewSQLBuilder()
		for _, idx := range as.DroppedIndices {
			gen.CreateIndexGenerator().Generate(diBuf, as.Name, idx)
			diBuf.WriteNewLine()
		}

		contents = append(contents, getContents(atBuf.Bytes(), diBuf.Bytes(), aiBuf.Bytes()))
	}
	return bytes.Join(contents, SectionSeparator)
}

func (gen *SqlGenerator) DownMigration(plan *step.MigrationPlanner) error {
	fmt.Println("🚀 Generating down database migration files")
	fmt.Printf("Target file: %s\n", color.HiBlueString(gen.dbDownFilename))

	createTableDown := gen.GenerateDropTables(plan.CreateTable)
	alterTables := gen.AlterTableDown(plan.AlterSchema)
	dropTableDown := []byte{}
	if !gen.flag.SkipDropTable {
		dropTableDown = gen.GenerateCreateTables(plan.DropTable)
	}

	content := getContents(createTableDown, dropTableDown, alterTables)
	if len(bytes.TrimSpace(content)) == 0 {
		fmt.Println(color.YellowString("No changes being detected, skipping..."))
		return nil
	}

	if gen.dialectOption.SupportTransaction {
		content = getContents(
			gen.dialectOption.BeginClause, content, gen.dialectOption.CommitClause,
		)
	}
	err := gen.Writer(gen.dbDownFilename, content)
	if err != nil {
		fmt.Println(color.RedString("Failed"))
		return err
	}
	fmt.Println(color.GreenString("Succeeded"))
	return nil
}

func (gen *SqlGenerator) Writer(filename string, content []byte) error {
	err := os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, content, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (gen *SqlGenerator) GenerateCreateTables(schemas []*config.Schema) []byte {
	sb := sb.NewSQLBuilder()
	for _, schema := range schemas {
		gen.CreateTableGenerator().Generate(sb, schema)
		sb.WriteNewLine()
		sb.WriteNewLine()

		for _, idx := range schema.Index {
			gen.CreateIndexGenerator().Generate(sb, schema.Name, idx)
			sb.WriteNewLine()
		}
		if len(schema.Index) > 0 {
			sb.WriteNewLine()
		}
	}

	return bytes.TrimSpace(sb.Bytes())
}

func (gen *SqlGenerator) GenerateDropTables(schemas []*config.Schema) []byte {
	sb := sb.NewSQLBuilder()
	for _, schema := range schemas {
		gen.DropTableGenerator().Generate(sb, schema)
		sb.WriteNewLine()
		sb.WriteNewLine()
	}

	return bytes.TrimSpace(sb.Bytes())
}

func getContents(contents ...[]byte) []byte {
	container := make([][]byte, 0)

	for _, content := range contents {
		content = bytes.TrimSpace(content)
		if len(content) > 0 {
			container = append(container, content)
		}
	}

	return bytes.TrimSpace(bytes.Join(container, SectionSeparator))
}
