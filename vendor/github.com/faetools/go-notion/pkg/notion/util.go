package notion

import (
	"fmt"
)

// TitleProperty represents a Title property.
var TitleProperty = PropertyMeta{
	Id:    "title", // must be this
	Type:  PropertyTypeTitle,
	Title: &map[string]interface{}{},
}

// NewRichTexts creates a RichTexts object with the desired content.
func NewRichTexts(content string) RichTexts {
	return RichTexts{NewRichText(content)}
}

// NewRichText creates a RichText object with the desired content.
func NewRichText(content string) RichText {
	return RichText{
		Type:        RichTextTypeText,
		PlainText:   content,
		Text:        &Text{Content: content},
		Annotations: Annotations{Color: ColorDefault},
	}
}

// GetNames returns names of all selected options.
func (opts PropertyOptions) GetNames() []string {
	names := make([]string, len(opts))

	for i, sel := range opts {
		names[i] = sel.Name
	}

	return names
}

// GetIDs returns the UUIDs of all references.
func (refs References) GetIDs() []UUID {
	ids := make([]UUID, len(refs))

	for i, ref := range refs {
		ids[i] = ref.Id
	}

	return ids
}

// URL return the URL of the file
func (f File) URL() string {
	switch f.Type {
	case FileTypeExternal:
		return f.External.Url
	case FileTypeFile:
		return f.File.Url
	default:
		panic(fmt.Errorf("invalid File of type %q", f.Type))
	}
}
