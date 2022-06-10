package command

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/dir"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/json"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/schema"
)

var (
	outputDir string
	DumpDbCmd = &cobra.Command{
		Use:     "dump:db",
		Short:   "Dump db",
		Long:    "This command is used to dump your current db and save it to your JSON schemas",
		Run:     DumpDb,
		Example: "dump:db -c postgresql://postgres:postgres@localhost:5432/playground -o examples/schemas",  // pragma: allowlist secret
	}
)

func init() {
	DumpDbCmd.Flags().StringVarP(&migrationConnString, "connection", "c", "", "(Required) Set connection string")
	DumpDbCmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output schemas path")
	DumpDbCmd.MarkFlagRequired("connection")
	DumpDbCmd.MarkFlagRequired("output")
}

func DumpDb(cmd *cobra.Command, args []string) {
	fmt.Println("🚀 Dumping database")
	crawler, err := schema.NewSchema(migrationConnString)
	if err != nil {
		fmt.Println(color.RedString("Failed to connect to datasource"))
		fmt.Println("Please see error details below:")
		fmt.Printf("\t%s\n", err)
		os.Exit(1)
	}

	gen := json.NewSchemasGenerator(crawler)
	override, err := dir.CheckDirExists(outputDir)
	if err != nil {
		fmt.Println(color.RedString("Check output path failed"))
		fmt.Println("Please see error details below:")
		fmt.Printf("\t%s\n", err)
		os.Exit(1)
	}

	if override {
		err = gen.GenerateBySchemas(outputDir)
		if err != nil {
			fmt.Println(color.RedString("Failed to generate schemas to JSON"))
			fmt.Println("Please see error details below:")
			fmt.Printf("\t%s\n", err)
			os.Exit(1)
		}
	}
	fmt.Println("Succeed dumping database")
}
