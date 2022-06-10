package main

import (
	"log"

	"github.com/spf13/cobra"
	"gitlab.com/wartek-id/core/tools/dbgen/command"
)

var (
	packageName string
	rootCmd     = &cobra.Command{
		Use:     "dbgen",
		Short:   "Wartool database generator",
		Version: release,
	}
)

func init() {
	rootCmd.AddCommand(command.GenCode)
	rootCmd.AddCommand(command.GenMigration)
	rootCmd.AddCommand(command.DumpDbCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
