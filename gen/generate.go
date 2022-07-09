package gen

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	_ "embed" // template

	"github.com/ettle/strcase"
	"github.com/faetools/cgtools"
	"github.com/faetools/go-notion/pkg/notion"
	"github.com/spf13/afero"
)

var (
	//go:embed property-values.tpl
	tplPropertyValuesRaw string

	tplPropertyValues = template.Must(template.New("property-values.tpl").Parse(tplPropertyValuesRaw))
)

type property struct {
	Key  string
	meta notion.PropertyMeta
}

func (p property) Name() string {
	return strcase.ToPascal(p.Key)
}

func (p property) GoType() string {
	switch p.meta.Type {
	case notion.PropertyTypeTitle,
		notion.PropertyTypeRichText:
		return "notion.RichTexts"
	case notion.PropertyTypeSelect:
		return "notion.SelectValue"
	case notion.PropertyTypeCheckbox:
		return "bool"
	case notion.PropertyTypeMultiSelect:
		return "notion.PropertyOptions"
	case notion.PropertyTypeNumber:
		if p.meta.Number.Format == notion.NumberConfigFormatNumber {
			return "int"
		}

		return "float32"
	case notion.PropertyTypeRelation:
		return "notion.References"
	default:
		return fmt.Sprintf("notion.%s", strcase.ToPascal(string(p.meta.Type)))
	}
}

func (p property) IsInt() bool {
	num := p.meta.Number
	return num != nil && num.Format == notion.NumberConfigFormatNumber
}

func (p property) GetFunc() string {
	return fmt.Sprintf("Get%s()", strcase.ToPascal(string(p.meta.Type)))
}

type ctxPropertyValues struct {
	PkgName    string
	Properties []property
}

// PropertyValues generates the go file associated with the property values of a database.
func PropertyValues(fs afero.Fs, pkgName string, m notion.PropertyMetaMap) error {
	g := cgtools.NewGenerator(fs)

	props := make([]property, 0, len(m))

	for key, val := range m {
		props = append(props, property{
			// we get only title keys from notion
			Key:  strings.Title(key),
			meta: val,
		})
	}

	// we want every run to have the same result
	sort.Slice(props, func(i, j int) bool {
		return props[i].Key < props[j].Key
	})

	return g.WriteTemplate(filepath.Join(pkgName, pkgName+".gen.go"),
		tplPropertyValues, ctxPropertyValues{
			PkgName:    pkgName,
			Properties: props,
		})
}
