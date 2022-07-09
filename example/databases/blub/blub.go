package blub

import (
	"github.com/faetools/go-notion-codegen/example/databases/bar"
	"github.com/faetools/go-notion/pkg/notion"
)

// Properties returns the property meta map for blub databases.
var Properties = func() notion.PropertyMetaMap {
	props := bar.Properties

	delete(props, "Tags")
	delete(props, "Rating")

	props["Labels"] = notion.PropertyMeta{
		Type: notion.PropertyTypeMultiSelect,
		MultiSelect: &notion.PropertyOptionsWrapper{
			Options: []notion.PropertyOption{},
		},
	}

	return props
}()
