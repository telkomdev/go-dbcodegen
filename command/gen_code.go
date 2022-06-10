package command

import (
	"bytes"
	"embed"
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gitlab.com/wartek-id/core/tools/dbgen/config"
	"gitlab.com/wartek-id/core/tools/dbgen/generator"
)

const (
	SqlcTemplatePath        = "template/sqlc.tmpl"
	DefaultPackageName      = "db"
	DefaultPackagePath      = "internal/db"
	DefaultSqlPackage       = "pgx/v4"
	SqlcPath                = "./sqlc.yaml"
	TemporaryPath           = "./temp"
	OutputQueriesPath       = TemporaryPath + "/query/"
	FullSchemaMigrationPath = TemporaryPath + "/fullschema/"
)

type SqlcConfig struct {
	PackageName string
	PackagePath string
	SqlPackage  string
}

var (
	packageName string
	packagePath string
	sqlPackage  string
	inputPath   string
	GenCode     = &cobra.Command{
		Use:     "gen:code",
		Short:   "Generate code",
		Long:    "This command is to generate the CRUD code",
		Run:     GenerateCode,
		Example: "gen:code -i examples/schemas",
	}
	//go:embed template/*
	sqlcTemplate embed.FS
)

func init() {
	GenCode.Flags().StringVar(&packageName, "packageName", DefaultPackageName, "(Optional) Package name for database in Go")
	GenCode.Flags().StringVar(&packagePath, "packagePath", DefaultPackagePath, "(Optional) Output path for database in Go")
	GenCode.Flags().StringVar(&sqlPackage, "sqlPackage", DefaultSqlPackage, "(Optional) SQL Package to specify which library to use")
	GenCode.Flags().StringVarP(&inputPath, "inputPath", "i", "", "(Required) Path to JSON schemas")
	GenCode.MarkFlagRequired("inputPath")
}

func GenerateCode(cmd *cobra.Command, args []string) {
	CreateSqlcConfigFile()
	schemas, err := config.Parse(inputPath)
	function := generator.GenerateQueries(schemas)
	err = SaveQueriesToFile(function, OutputQueriesPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = ExecuteSqlc()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = RemoveUnusedFiles()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ExecuteSqlc() error {
	fmt.Println("🚀 Generating code")
	_, err := exec.LookPath("sqlc")
	if err != nil {
		DownloadSqlc()
	}
	cmd := exec.Command("sqlc", "generate")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	if err = cmd.Run(); err != nil {
		fmt.Println(color.RedString("Generate code failed"))
		return err
	}
	fmt.Println(color.GreenString("Succeeded"))
	return nil
}

func DownloadSqlc() {
	cmd := exec.Command("go", "install", "github.com/kyleconroy/sqlc/cmd/sqlc@v1.13.0")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	if err := cmd.Run(); err != nil {
		fmt.Println("Download sqlc failed")
		panic(err)
	}
}

func CreateSqlcConfigFile() error {
	sqlcConfig := &SqlcConfig{
		PackageName: packageName,
		PackagePath: packagePath,
		SqlPackage:  sqlPackage,
	}

	template, err := template.ParseFS(sqlcTemplate, SqlcTemplatePath)
	if err != nil {
		return err
	}

	f, err := os.Create(SqlcPath)
	if err != nil {
		return err
	}
	err = template.Execute(f, sqlcConfig)
	if err != nil {
		return err
	}
	f.Close()

	return nil
}

func SaveQueriesToFile(function []*config.Function, outputPath string) error {
	var buf bytes.Buffer
	curFilename := ""
	curTable := ""
	err := os.MkdirAll(outputPath, 0755)
	if err != nil {
		return err
	}

	for _, f := range function {
		if curTable != f.TableName {
			if curTable != "" {
				os.WriteFile(curFilename, buf.Bytes(), 0644)
				buf = *&bytes.Buffer{}
				fmt.Printf(color.GreenString("Succeed generate query, target file: %s\n"), color.HiBlueString(curFilename))
			}
			curTable = f.TableName
			curFilename = filepath.Join(outputPath, curTable+".sql")
			fmt.Printf("\n🚀 Generating query to file for table: %s\n", f.TableName)
		}
		buf.WriteString("-- name: " + f.Name + " " + f.SqlcType + "\n")
		buf.WriteString(f.Query + ";\n\n")
	}
	os.WriteFile(curFilename, buf.Bytes(), 0644)
	fmt.Printf(color.GreenString("Succeed generate query, target file: %s\n"), color.HiBlueString(curFilename))
	return nil
}

func RemoveUnusedFiles() error {
	err := os.RemoveAll(SqlcPath)
	if err != nil {
		return err
	}

	err = os.RemoveAll(TemporaryPath)
	if err != nil {
		return err
	}

	return nil
}
