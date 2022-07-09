package foo

import "github.com/faetools/go-notion/pkg/notion"

var emptyConfig = &map[string]interface{}{}

// Properties returns the property meta map for foo databases.
func Properties(expanded bool) notion.PropertyMetaMap {
	props := notion.PropertyMetaMap{
		"Title": notion.TitleProperty,
		"Summary": notion.PropertyMeta{
			Type:     notion.PropertyTypeRichText,
			RichText: emptyConfig,
		},
	}

	if expanded {
		props["Important"] = notion.PropertyMeta{
			Type:     notion.PropertyTypeCheckbox,
			Checkbox: emptyConfig,
		}
	}

	return props
}
