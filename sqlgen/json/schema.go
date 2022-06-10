package json

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/schema"
)

type JsonSchemasGenerator struct {
	schemas schema.Schema
}

func NewSchemasGenerator(schemas schema.Schema) *JsonSchemasGenerator {
	return &JsonSchemasGenerator{
		schemas: schemas,
	}
}

func (s *JsonSchemasGenerator) GenerateBySchemas(outputDir string) error {
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return err
	}

	currentSchemas, err := s.schemas.GetSchemas()
	if err != nil {
		return err
	}

	for _, s := range currentSchemas {
		fmt.Println("\nDumping db: " + s.Name)
		filename := filepath.Join(outputDir, s.Name+".json")
		file, err := json.MarshalIndent(s, "", " ")
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(filename, file, 0644)
		if err != nil {
			return err
		}
		fmt.Printf(color.GreenString("Succeed dumping db: %s, target file: %s\n"), color.HiBlueString(s.Name), color.HiBlueString(filename))
	}
	return nil
}
