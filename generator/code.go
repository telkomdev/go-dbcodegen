package generator

import (
	"fmt"
	"strconv"

	"gitlab.com/wartek-id/core/tools/dbgen/config"

	goqu "github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/fatih/color"
	"github.com/iancoleman/strcase"
)

func GenerateQueries(schemas []*config.Schema) []*config.Function {
	dialect := goqu.Dialect("postgres")
	var function []*config.Function
	for _, element := range schemas {
		fmt.Printf("🚀 Generating query for table: %s\n", element.Name)
		for _, f := range GenerateSelectQueryByIndex(dialect, element) {
			function = append(function, f)
		}
		function = append(function, GenerateUpdateQuery(dialect, element))
		function = append(function, GenerateInsertQuery(dialect, element))
		sql := GenerateDeleteQuery(dialect, element)
		if sql != nil {
			function = append(function, sql)
		}

		sql = GenerateRestoreQuery(dialect, element)
		if sql != nil {
			function = append(function, sql)
		}

		sql, err := GenerateDestroyQuery(dialect, element)
		if err != nil {
			fmt.Println(color.RedString("Fail creating query for table: %s\n", element.Name))
			panic(err)
		}
		function = append(function, sql)
		fmt.Println(color.GreenString("Succeeded"))
	}
	return function
}

func GenerateSelectQuery(dialect goqu.DialectWrapper, element *config.Schema, index goqu.Ex) (string, error) {
	ds := dialect.From(element.Name).Where(index)
	sql, _, err := ds.Prepared(true).ToSQL()
	if err != nil {
		return "", err
	}

	return sql, nil
}

func GenerateInsertQuery(dialect goqu.DialectWrapper, element *config.Schema) *config.Function {
	var cols []interface{}
	var vals []interface{}
	for i, e := range element.Fields {
		if e.Name == "id" || e.Name == "deleted_at" || e.Name == "updated_at" {
			continue
		}
		cols = append(cols, e.Name)
		vals = append(vals, "$"+strconv.Itoa(i+1))
	}

	ds := dialect.Insert(element.Name).Cols(cols...).Vals(vals)
	sql, _, _ := ds.Prepared(true).ToSQL()
	functionName := GenerateFunctionName("create_"+element.Name, "")
	return &config.Function{
		Name:      functionName,
		Query:     sql,
		TableName: element.Name,
		SqlcType:  ":exec",
	}
}

func GenerateUpdateQuery(dialect goqu.DialectWrapper, element *config.Schema) *config.Function {
	filter := goqu.And(goqu.Ex{"id": "1"})
	cols := goqu.Record{}
	for _, e := range element.Fields {
		if e.Name == "created_at" || e.Name == "id" {
			continue
		} else if e.Name == "deleted_at" {
			filter = filter.Append(goqu.I("deleted_at").IsNull())
			continue
		}
		cols[e.Name] = e.Name
	}

	ds := dialect.Update(element.Name).Where(filter).Set(cols)
	sql, _, _ := ds.Prepared(true).ToSQL()
	functionName := GenerateFunctionName("update_"+element.Name, "")
	return &config.Function{
		Name:      functionName,
		Query:     sql,
		TableName: element.Name,
		SqlcType:  ":exec",
	}
}

func GenerateDestroyQuery(dialect goqu.DialectWrapper, element *config.Schema) (*config.Function, error) {
	ds := dialect.Delete(element.Name).Where(goqu.Ex{"id": "$1"})
	sql, _, err := ds.Prepared(true).ToSQL()
	if err != nil {
		return nil, err
	}

	functionName := GenerateFunctionName("destroy_"+element.Name, "")
	return &config.Function{
		Name:      functionName,
		Query:     sql,
		TableName: element.Name,
		SqlcType:  ":exec",
	}, nil
}

func GenerateDeleteQuery(dialect goqu.DialectWrapper, element *config.Schema) *config.Function {
	filter := goqu.And(goqu.Ex{"id": "1"})
	for _, f := range element.Fields {
		if f.Name == "deleted_at" {
			filter = filter.Append(goqu.I("deleted_at").IsNull())
			ds := dialect.Update(element.Name).Where(filter).Set(goqu.Record{"deleted_at": goqu.L("NOW()")})
			sql, _, _ := ds.Prepared(true).ToSQL()

			functionName := GenerateFunctionName("delete_"+element.Name, "")
			return &config.Function{
				Name:      functionName,
				Query:     sql,
				TableName: element.Name,
				SqlcType:  ":exec",
			}
		}
	}
	return nil
}

func GenerateRestoreQuery(dialect goqu.DialectWrapper, element *config.Schema) *config.Function {
	filter := goqu.And(goqu.Ex{"id": "1"})
	for _, f := range element.Fields {
		if f.Name == "deleted_at" {
			filter = filter.Append(goqu.I("deleted_at").IsNotNull())
			ds := dialect.Update(element.Name).Where(filter).Set(goqu.Record{"deleted_at": goqu.L("NULL")})
			sql, _, _ := ds.Prepared(true).ToSQL()

			functionName := GenerateFunctionName("restore_"+element.Name, "")
			return &config.Function{
				Name:      functionName,
				Query:     sql,
				TableName: element.Name,
				SqlcType:  ":exec",
			}
		}
	}
	return nil
}

func GenerateSelectQueryByIndex(dialect goqu.DialectWrapper, element *config.Schema) []*config.Function {
	// Initialize FindById
	index := goqu.Ex{"id": "1"}
	var function []*config.Function
	sql, err := GenerateSelectQuery(dialect, element, index)
	if err != nil {
		panic(err)
	}
	functionName := GenerateFunctionName("find_"+element.Name, "_id")
	function = append(function, config.ParseFunction(functionName, sql, element.Name, ":one"))

	// Initialize Find with indexes
	for _, e := range element.Index {
		m := goqu.Ex{}
		cnt := len(e.Fields)
		fields := ""
		for i, field := range e.Fields {
			m[field.Column] = "random"
			fields += "_" + field.Column
			if i < cnt-1 {
				fields += "_and"
			}
		}
		index = goqu.Ex(m)
		sql, err = GenerateSelectQuery(dialect, element, index)
		if err != nil {
			panic(err)
		}
		functionName := GenerateFunctionName("find_"+element.Name, fields)

		function = append(function, config.ParseFunction(functionName, sql, element.Name, ":many"))
	}
	return function
}

func GenerateFunctionName(actionName string, fields string) string {
	if fields != "" {
		actionName += "_by" + fields
	}
	return strcase.ToCamel(actionName)
}
