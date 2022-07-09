package foo

import "github.com/faetools/go-notion/pkg/notion"

type PropertyValues struct {
	Important bool
	Summary   notion.RichTexts
	Title     notion.RichTexts
}

func GetPropertyValues(props notion.PropertyValueMap) PropertyValues {
	return PropertyValues{
		Important: props["Important"].GetCheckbox(),
		Summary:   props["Summary"].GetRichText(),
		Title:     props["Title"].GetTitle(),
	}
}
