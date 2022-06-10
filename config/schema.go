package config

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type Schema struct {
	Name   string   `json:"name"`
	Fields []*Field `json:"fields"`
	Index  []*Index `json:"indexes"`
}

func (s *Schema) GetName() string {
	return s.Name
}

func ParseSchema(path string) (*Schema, error) {
	var schema Schema
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &schema)
	if err != nil {
		return nil, err
	}
	return &schema, nil
}

func ParseDir(rootPath string) ([]*Schema, error) {
	schemas := make([]*Schema, 0)

	rootPath, err := filepath.Abs(rootPath)
	if err != nil {
		return schemas, err
	}

	wd, _ := os.Getwd()
	err = filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, _ error) error {
		if d.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".json" {
			return nil
		}
		schema, err := ParseSchema(path)
		if err != nil {
			filename, _ := filepath.Rel(wd, path)
			return errors.Wrap(err, filename)
		}

		schemas = append(schemas, schema)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return schemas, nil
}

func Parse(paths ...string) ([]*Schema, error) {
	var schemes []*Schema
	for _, path := range paths {
		path, err := filepath.Abs(path)
		if err != nil {
			return nil, err
		}

		stat, err := os.Stat(path)
		if os.IsNotExist(err) {
			return nil, err
		}

		if stat.IsDir() {
			pathSchemas, err := ParseDir(path)
			if err != nil {
				return nil, err
			}
			schemes = append(schemes, pathSchemas...)
		} else {
			scheme, err := ParseSchema(path)
			if err != nil {
				return nil, err
			}
			schemes = append(schemes, scheme)
		}
	}
	return schemes, nil
}
