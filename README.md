# Go Database Code Generator
This tool is to help you generate schema migrations and CRUD code in Golang from an entity definition in form of JSON.

## Why
This tool is inspired by an amazing video titled [Design Microservice Architectures the Right Way](https://www.youtube.com/watch?v=j6ow-UemzBc).
Plenty of reasons why this tool is built can be found there. One of the example, according to that video, is when things 
are slow because some queries don't have index, the hero is not the person who fix it by adding index. However, the unsung 
hero is the one who prevents it from ever happened in the first place. One way to do that is to provide a tool to generate high
quality database access code. This dbcodegen is aimed to be such tool.

## Supported Command
## gen:migration
Generate database migration file

Command:
```
dbgen gen:migration -c {connection_string} -d {output directory} -o {output filename} input file(s)/folder(s)
```
Example:
```
dbgen gen:migration -c postgresql://postgres:postgres@localhost:5432/playground -d db/migration -o  users_registrations db/schemas
```

## gen:code
Generate schemas and queries into code

Command:
```
dbgen gen:code -i {input folder}
```

Example:
```
dbgen gen:code -i examples/schemas
```

Notes:
gen:code also has some optional fields:
1. packageName: package name for database in Go
2. packagePath: output path for database in Go
3. sqlPackage: SQL Package to specify which library to use

## dump:db
Dump current schemas into JSON files

Command:
```
dbgen dump:db -c {connection} -o {output path}
```

Example:
```
dump:db -c postgresql://postgres:postgres@localhost:5432/playground -o examples/schemas
```

## Input file Example
The input is JSON file containing structures of an entity. Complete schema spec can be found [here](https://github.com/telkomdev/go-dbcodegen/blob/main/examples/schemas/json-schema-spec.md)
```
{
  "name": "example",
  "fields": [
    {
      "name": "id",
      "type": "bigserial",
      "options": [
        "primary key"
      ]
    },
    {
      "name": "name",
      "type": "varchar",
      "limit": 50,
      "options": [
        "not null"
      ]
    }
  ],
  "indexes": [
    {
      "name": "index_example_on_name",
      "fields": [
        {
          "column": "name"
        }
      ],
      "unique": true
    }
  ]
}
```
Other examples can be found [here](https://github.com/telkomdev/go-dbcodegen/tree/main/examples/schemas)

