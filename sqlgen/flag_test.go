package sqlgen_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen"
)

func TestNewFlag(t *testing.T) {
	flag, err := sqlgen.NewFlag("db/migration", "migration", false)
	assert.NoError(t, err)
	assert.NotNil(t, flag)

	flag, err = sqlgen.NewFlag("db/migration", "migration/001", false)
	assert.Error(t, err, "output target cannot contain \"\\\" character")
	assert.Nil(t, flag)
}
