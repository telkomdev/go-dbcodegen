package sb

import "bytes"

type SQLBuilder interface {
	Write([]byte) SQLBuilder
	WriteString(...string) SQLBuilder
	WriteRunes(...rune) SQLBuilder
	WriteNewLine() SQLBuilder
	ToSQL() (string, error)
	Bytes() []byte
	String() string
}

type sqlBuilder struct {
	buf *bytes.Buffer
	err error
}

func NewSQLBuilder() SQLBuilder {
	return &sqlBuilder{
		buf: &bytes.Buffer{},
	}
}

func (b *sqlBuilder) Write(bs []byte) SQLBuilder {
	if b.err == nil {
		_, _ = b.buf.Write(bs)
	}

	return b
}

func (b *sqlBuilder) WriteString(ss ...string) SQLBuilder {
	if b.err == nil {
		for _, s := range ss {
			_, _ = b.buf.WriteString(s)
		}
	}
	return b
}

func (b *sqlBuilder) WriteRunes(rs ...rune) SQLBuilder {
	if b.err == nil {
		for _, r := range rs {
			_, _ = b.buf.WriteRune(r)
		}
	}
	return b
}

func (b *sqlBuilder) ToSQL() (string, error) {
	if b.err != nil {
		return "", b.err
	}

	return b.buf.String(), nil
}

func (b *sqlBuilder) Bytes() []byte {
	return b.buf.Bytes()
}

func (b *sqlBuilder) String() string {
	return b.buf.String()
}

func (b *sqlBuilder) WriteNewLine() SQLBuilder {
	return b.WriteString("\n")
}
