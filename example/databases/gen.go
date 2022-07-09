package main

import (
	"log"

	"github.com/faetools/go-notion-codegen/example/databases/bar"
	"github.com/faetools/go-notion-codegen/example/databases/blub"
	"github.com/faetools/go-notion-codegen/example/databases/foo"
	"github.com/faetools/go-notion-codegen/gen"
	"github.com/faetools/go-notion/pkg/notion"
	"github.com/spf13/afero"
)

//go:generate go run gen.go

func main() {
	fs := afero.NewOsFs()

	for pkgName, props := range map[string]notion.PropertyMetaMap{
		"bar":  bar.Properties,
		"blub": blub.Properties,
		"foo":  foo.Properties(true),
	} {
		if err := gen.PropertyValues(fs, pkgName, props); err != nil {
			log.Fatal(err)
		}
	}
}
