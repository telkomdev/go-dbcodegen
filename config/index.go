package config

import "encoding/json"

type Index struct {
	Name   string        `json:"name"`
	Fields []*IndexField `json:"fields"`
	Unique bool          `json:"unique"`
}

type IndexField struct {
	Column string `json:"column"`
	Order  string `json:"order"`
}

func (f *IndexField) UnmarshalJSON(data []byte) error {
	type fieldAlias IndexField
	field := fieldAlias{
		Order: "ASC",
	}
	err := json.Unmarshal(data, &field)
	if err != nil {
		return err
	}
	*f = IndexField(field)
	return nil
}

func (i *Index) GetName() string {
	return i.Name
}

func (i *Index) GetColumns() []string {
	cols := []string{}
	for _, field := range i.Fields {
		cols = append(cols, field.Column)
	}

	return cols
}
