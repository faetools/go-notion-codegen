package bar

import (
	"context"

	"github.com/faetools/go-notion/pkg/notion"
)

var emptyConfig = &map[string]interface{}{}

var noOptions = &notion.PropertyOptionsWrapper{
	Options: []notion.PropertyOption{},
}

// Properties returns the property meta map for bar databases.
var Properties = notion.PropertyMetaMap{
	"Name": notion.PropertyMeta{
		Id:    "title",
		Name:  "Name",
		Type:  notion.PropertyTypeTitle,
		Title: emptyConfig,
	},
	"Draft": notion.PropertyMeta{
		Type:     notion.PropertyTypeCheckbox,
		Checkbox: emptyConfig,
	},
	"Expires": notion.PropertyMeta{
		Type: notion.PropertyTypeDate,
		Date: emptyConfig,
	},
	"Resources": notion.PropertyMeta{
		Type:  notion.PropertyTypeFiles,
		Files: emptyConfig,
	},
	"Tags": notion.PropertyMeta{
		Type:        notion.PropertyTypeMultiSelect,
		MultiSelect: noOptions,
	},
	"Number of People": notion.PropertyMeta{
		Type: notion.PropertyTypeNumber,
		Number: &notion.NumberConfig{
			Format: notion.NumberConfigFormatNumber,
		},
	},
	"Rating": notion.PropertyMeta{
		Type: notion.PropertyTypeNumber,
		Number: &notion.NumberConfig{
			Format: notion.NumberConfigFormatNumberWithCommas,
		},
	},
	"Related To": notion.PropertyMeta{
		Type: notion.PropertyTypeRelation,
		Relation: &notion.RelationConfiguration{
			DatabaseId: "some id",
		},
	},
	"Description": notion.PropertyMeta{
		Type:     notion.PropertyTypeRichText,
		RichText: emptyConfig,
	},
	"Category": notion.PropertyMeta{
		Type:   notion.PropertyTypeSelect,
		Select: noOptions,
	},
}

// CreateDatabase creates a bar database with the properties we want.
func CreateDatabase(ctx context.Context, cli *notion.Client, parentID notion.UUID) (*notion.Database, error) {
	return cli.CreateNotionDatabase(ctx, notion.Database{
		Title:      notion.NewRichTexts("My Bar Database"),
		Parent:     &notion.Parent{PageId: parentID},
		Properties: Properties,
	})
}

// UpdateDatabase updates a bar database so that it has the properties we want.
func UpdateDatabase(ctx context.Context, cli *notion.Client, id notion.UUID) (*notion.Database, error) {
	return cli.UpdateNotionDatabase(ctx, notion.Database{
		Id:         id,
		Properties: Properties,
	})
}
