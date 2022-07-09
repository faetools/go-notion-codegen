package gen_test

import (
	"os"
	"testing"

	"github.com/faetools/go-notion-codegen/gen"
	"github.com/faetools/go-notion/pkg/notion"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var noOptions = &notion.PropertyOptionsWrapper{
	Options: []notion.PropertyOption{},
}

var emptyConfig = &map[string]interface{}{}

func TestPropertyValues(t *testing.T) {
	t.Parallel()
	os.Stdout = nil

	memFs := afero.NewMemMapFs()

	require.NoError(t, gen.PropertyValues(memFs, "mypackage",
		notion.PropertyMetaMap{
			"My Title": notion.TitleProperty,
			"check": notion.PropertyMeta{
				Type:     notion.PropertyTypeCheckbox,
				Checkbox: emptyConfig,
			},
			"my date": notion.PropertyMeta{
				Type: notion.PropertyTypeDate,
				Date: emptyConfig,
			},
			"my files": notion.PropertyMeta{
				Type:  notion.PropertyTypeFiles,
				Files: emptyConfig,
			},
			"my multi select": notion.PropertyMeta{
				Type:        notion.PropertyTypeMultiSelect,
				MultiSelect: noOptions,
			},
			"my number": notion.PropertyMeta{
				Type: notion.PropertyTypeNumber,
				Number: &notion.NumberConfig{
					Format: notion.NumberConfigFormatNumber,
				},
			},
			"my float": notion.PropertyMeta{
				Type: notion.PropertyTypeNumber,
				Number: &notion.NumberConfig{
					Format: notion.NumberConfigFormatNumberWithCommas,
				},
			},
			"my relation": notion.PropertyMeta{
				Type:     notion.PropertyTypeRelation,
				Relation: nil,
			},
			"my richtext": notion.PropertyMeta{
				Type:     notion.PropertyTypeRichText,
				RichText: emptyConfig,
			},
			"my select": notion.PropertyMeta{
				Type:   notion.PropertyTypeSelect,
				Select: noOptions,
			},
		}))

	b, err := afero.ReadFile(memFs, "mypackage/mypackage.gen.go")
	assert.NoError(t, err)

	assert.Equal(t, `package mypackage

import "github.com/faetools/go-notion/pkg/notion"

type PropertyValues struct {
	Check         bool
	MyDate        notion.Date
	MyFiles       notion.Files
	MyFloat       float32
	MyMultiSelect notion.PropertyOptions
	MyNumber      int
	MyRelation    notion.References
	MyRichtext    notion.RichTexts
	MySelect      notion.SelectValue
	MyTitle       notion.RichTexts
}

func GetPropertyValues(props notion.PropertyValueMap) PropertyValues {
	return PropertyValues{
		Check:         props["Check"].GetCheckbox(),
		MyDate:        props["My Date"].GetDate(),
		MyFiles:       props["My Files"].GetFiles(),
		MyFloat:       props["My Float"].GetNumber(),
		MyMultiSelect: props["My Multi Select"].GetMultiSelect(),
		MyNumber:      int(props["My Number"].GetNumber()),
		MyRelation:    props["My Relation"].GetRelation(),
		MyRichtext:    props["My Richtext"].GetRichText(),
		MySelect:      props["My Select"].GetSelect(),
		MyTitle:       props["My Title"].GetTitle(),
	}
}
`, string(b))
}
