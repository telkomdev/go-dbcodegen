package config

type Function struct {
	Name      string
	Query     string
	TableName string
	SqlcType  string
}

func ParseFunction(name string, query string, tableName string, sqlcType string) *Function {
	return &Function{
		Name:      name,
		Query:     query,
		TableName: tableName,
		SqlcType:  sqlcType,
	}
}
