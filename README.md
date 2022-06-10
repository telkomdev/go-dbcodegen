# Go Database Code Generator
This tool is to help you generate schema migrations and CRUD code in Golang

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
1. packageName: Package name for database in Go
2. packagePath: Output path for database in Go
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
