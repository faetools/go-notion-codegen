package {{ .PkgName }}

import "github.com/faetools/go-notion/pkg/notion"

type PropertyValues struct {
{{- range .Properties }}
	{{ .Name }} {{ .GoType -}}
{{ end }}
}

func GetPropertyValues(props notion.PropertyValueMap) PropertyValues {
	return PropertyValues{
	{{- range .Properties }}
		{{ .Name }}:
			{{- if .IsInt }}int({{ end -}}
			props["{{ .Key }}"].{{ .GetFunc }}
			{{- if .IsInt -}}){{ end }},
	{{- end }}
	}
}
