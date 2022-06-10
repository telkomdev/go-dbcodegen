package command

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gitlab.com/wartek-id/core/tools/dbgen/config"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/schema"
)

const (
	DefaultConnectionString = ""
	DefaultOutputDirectory  = "db/migration"
	DefaultOutputName       = ""
	DefaultSkipTable        = false
)

var GenMigration = &cobra.Command{
	Use:     "gen:migration [flags] file(s)/folder(s)...",
	Short:   "Generate SQL migration",
	Long:    "This command is used to generate new migrations and the full schema",
	Run:     GenerateMigration,
	Example: "gen:migration -c postgresql://postgres:postgres@localhost:5432/playground -o migration user.json logging.json",  // pragma: allowlist secret
}

var (
	migrationConnString,
	migrationDir,
	migrationOutput string

	skipDropTable bool
)

func init() {
	GenMigration.Flags().StringVarP(&migrationConnString, "connection", "c", DefaultConnectionString, "set connection string")
	GenMigration.Flags().StringVarP(&migrationDir, "dir", "d", DefaultOutputDirectory, "set migration directory")
	GenMigration.Flags().StringVarP(&migrationOutput, "output", "o", DefaultOutputName, "set output name")
	GenMigration.Flags().BoolVar(&skipDropTable, "skip-drop-table", DefaultSkipTable, "skip drop table generation query")
}

func GenerateMigration(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Print(color.RedString("Command Failed "))
		fmt.Println("Please specify schema file(s)/folder(s)")
		os.Exit(1)
	}

	schemas, err := config.Parse(args...)
	if err != nil {
		fmt.Println(color.RedString("Database Generation Failed"))
		fmt.Println("Please see error details below:")
		fmt.Printf("\t%s\n", err)
		os.Exit(1)
	}

	crawler, err := schema.NewSchema(migrationConnString)
	if err != nil {
		fmt.Println(color.RedString("Failed to connect to datasource"))
		fmt.Println("Please see error details below:")
		fmt.Printf("\t%s\n", err)
		os.Exit(1)
	}

	flag, err := sqlgen.NewFlag(migrationDir, migrationOutput, skipDropTable)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	gen := sqlgen.NewGenerator(crawler, schemas, flag)
	err = gen.Generate()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
