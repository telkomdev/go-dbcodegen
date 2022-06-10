package dialect

import (
	"bytes"

	"gitlab.com/wartek-id/core/tools/dbgen/types/field_option"
	"gitlab.com/wartek-id/core/tools/dbgen/types/field_type"
)

type DialectOption struct {
	CreateClause []byte
	DropClause   []byte
	AlterClause  []byte
	BeginClause  []byte
	CommitClause []byte

	IndexFragment []byte
	TableFragment []byte

	AlterFragment    []byte
	DropFragment     []byte
	AddFragment      []byte
	ColumnFragment   []byte
	SetFragment      []byte
	DefaultFragment  []byte
	DataTypeFragment []byte

	BooleanFragment     []byte
	VarcharFragment     []byte
	SmallIntFragment    []byte
	IntFragment         []byte
	BigIntFragment      []byte
	JsonFragment        []byte
	JsonbFragment       []byte
	FloatFragment       []byte
	DecimalFragment     []byte
	TimestampFragment   []byte
	TimestamptzFragment []byte

	SmallSerialFragment []byte
	SerialFragment      []byte
	BigSerialFragment   []byte

	PrimaryKeyFragment []byte
	NullableFragment   []byte
	NotNullFragment    []byte
	UniqueFragment     []byte

	ConcurrentlyFragment  []byte
	IfNotExistsFragment   []byte
	IfExistsFragment      []byte
	AutoIncrementFragment []byte

	EmptyFragment        []byte
	OnFragment           []byte
	CommaNewLineFragment []byte
	SupportConcurrently  bool
	SupportTransaction   bool

	LeftParenRune   rune
	RightParenRune  rune
	CommaRune       rune
	SemiColonRune   rune
	SpaceRune       rune
	QuoteRune       rune
	StringQuoteRune rune
	NewLineRune     rune
	TabRune         rune

	DataTypesLookup    map[field_type.FieldType][]byte
	FieldOptionsLookup map[field_option.FieldOption][]byte
}

func DefaultDialectOption() *DialectOption {
	do := &DialectOption{
		CreateClause: []byte("CREATE "),
		DropClause:   []byte("DROP "),
		AlterClause:  []byte("ALTER "),
		BeginClause:  []byte("BEGIN;"),
		CommitClause: []byte("COMMIT;"),

		IndexFragment: []byte("INDEX "),
		TableFragment: []byte("TABLE "),

		AlterFragment:    []byte("ALTER "),
		DropFragment:     []byte("DROP "),
		AddFragment:      []byte("ADD "),
		ColumnFragment:   []byte("COLUMN "),
		SetFragment:      []byte("SET "),
		DefaultFragment:  []byte("DEFAULT "),
		DataTypeFragment: []byte("DATA TYPE "),

		BooleanFragment:     []byte("BOOLEAN"),
		VarcharFragment:     []byte("VARCHAR"),
		SmallIntFragment:    []byte("SMALLINT"),
		IntFragment:         []byte("INT"),
		BigIntFragment:      []byte("BIGINT"),
		JsonFragment:        []byte("JSON"),
		JsonbFragment:       []byte("JSONB"),
		FloatFragment:       []byte("FLOAT"),
		DecimalFragment:     []byte("DECIMAL"),
		TimestampFragment:   []byte("TIMESTAMP"),
		TimestamptzFragment: []byte("TIMESTAMPTZ"),

		SmallSerialFragment: []byte("SMALLSERIAL"),
		SerialFragment:      []byte("SERIAL"),
		BigSerialFragment:   []byte("BIGSERIAL"),

		PrimaryKeyFragment: []byte("PRIMARY KEY"),
		NullableFragment:   []byte("NULL"),
		NotNullFragment:    []byte("NOT NULL"),
		UniqueFragment:     []byte("UNIQUE"),

		ConcurrentlyFragment:  []byte("CONCURRENTLY"),
		IfNotExistsFragment:   []byte("IF NOT EXISTS "),
		IfExistsFragment:      []byte("IF EXISTS "),
		AutoIncrementFragment: []byte("AUTO INCREMENT"),

		LeftParenRune:   '(',
		RightParenRune:  ')',
		CommaRune:       ',',
		SemiColonRune:   ';',
		SpaceRune:       ' ',
		QuoteRune:       '"',
		StringQuoteRune: '\'',
		NewLineRune:     '\n',
		TabRune:         '\t',

		CommaNewLineFragment: []byte(",\n"),
		OnFragment:           []byte(" ON "),
		SupportConcurrently:  false,
		SupportTransaction:   true,
	}

	do.DataTypesLookup = map[field_type.FieldType][]byte{
		field_type.Boolean:     do.BooleanFragment,
		field_type.Varchar:     do.VarcharFragment,
		field_type.SmallInt:    do.SmallIntFragment,
		field_type.Int:         do.IntFragment,
		field_type.BigInt:      do.BigIntFragment,
		field_type.Json:        do.JsonFragment,
		field_type.Jsonb:       do.JsonbFragment,
		field_type.Float:       do.FloatFragment,
		field_type.Decimal:     do.DecimalFragment,
		field_type.Timestamp:   do.TimestampFragment,
		field_type.Timestamptz: do.TimestamptzFragment,
		field_type.BigSerial:   do.BigSerialFragment,
		field_type.Serial:      do.SerialFragment,
		field_type.SmallSerial: do.SmallSerialFragment,
	}

	do.FieldOptionsLookup = map[field_option.FieldOption][]byte{
		field_option.Nullable:      do.NullableFragment,
		field_option.NotNull:       do.NotNullFragment,
		field_option.AutoIncrement: do.AutoIncrementFragment,
		field_option.Unique:        do.UniqueFragment,
		field_option.PrimaryKey:    do.PrimaryKeyFragment,
	}

	return do
}

func (do *DialectOption) CreateTableTemplate() []byte {
	buf := bytes.Buffer{}
	buf.Write(do.CreateClause)
	buf.Write(do.TableFragment)
	return buf.Bytes()
}

func (do *DialectOption) AlterColumnTemplate() []byte {
	buf := bytes.Buffer{}
	buf.Write(do.AlterClause)
	buf.Write(do.ColumnFragment)
	return buf.Bytes()
}

func (do *DialectOption) AddColumnTemplate() []byte {
	buf := bytes.Buffer{}
	buf.Write(do.AddFragment)
	buf.Write(do.ColumnFragment)
	return buf.Bytes()
}

func (do *DialectOption) DropColumnTemplate() []byte {
	buf := bytes.Buffer{}
	buf.Write(do.DropClause)
	buf.Write(do.ColumnFragment)
	return buf.Bytes()
}
