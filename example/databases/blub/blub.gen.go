package blub

import "github.com/faetools/go-notion/pkg/notion"

type PropertyValues struct {
	Category       notion.SelectValue
	Description    notion.RichTexts
	Draft          bool
	Expires        notion.Date
	Labels         notion.PropertyOptions
	Name           notion.RichTexts
	NumberOfPeople int
	RelatedTo      notion.References
	Resources      notion.Files
}

func GetPropertyValues(props notion.PropertyValueMap) PropertyValues {
	return PropertyValues{
		Category:       props["Category"].GetSelect(),
		Description:    props["Description"].GetRichText(),
		Draft:          props["Draft"].GetCheckbox(),
		Expires:        props["Expires"].GetDate(),
		Labels:         props["Labels"].GetMultiSelect(),
		Name:           props["Name"].GetTitle(),
		NumberOfPeople: int(props["Number Of People"].GetNumber()),
		RelatedTo:      props["Related To"].GetRelation(),
		Resources:      props["Resources"].GetFiles(),
	}
}
