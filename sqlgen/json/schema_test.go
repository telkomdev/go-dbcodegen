package json_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/wartek-id/core/tools/dbgen/config"
	"gitlab.com/wartek-id/core/tools/dbgen/sqlgen/json"
	mock_schema "gitlab.com/wartek-id/core/tools/dbgen/sqlgen/mocks/schema"
	"gitlab.com/wartek-id/core/tools/dbgen/types/field_option"
)

func TestGenerateBySchemas(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	outputDir := t.TempDir()
	mockCrawler := mock_schema.NewMockSchema(ctrl)
	mockCrawler.EXPECT().GetSchemas().Return([]*config.Schema{
		{
			Name: "example",
			Fields: []*config.Field{
				{
					Name: "id",
					Type: "bigserial",
					Options: []field_option.FieldOption{
						field_option.PrimaryKey,
					},
				},
				{
					Name:  "name",
					Type:  "varchar",
					Limit: 100,
				},
			},
		},
	}, nil).AnyTimes()
	gen := json.NewSchemasGenerator(mockCrawler)
	err := gen.GenerateBySchemas(outputDir)

	assert.NoError(t, err)
}
