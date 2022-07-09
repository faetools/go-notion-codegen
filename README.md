# go-notion-codegen

![Devtool version](https://img.shields.io/badge/Devtool-0.0.18-brightgreen.svg)
![Maintainer](https://img.shields.io/badge/team-firestarters-blue)

## About

This repository expands [go-notion](https://github.com/faetools/go-notion) so that you can generate code that helps you use your particular databases.

## Usage

### Generate Database Property Values

You can generate code that will transform a `notion.PropertyValueMap` into a struct so that you can more **easily get the values of a database entry**.

To do this, create a package for each database and define the properties of each database.

For example, you could have three databases, `foo`, `bar`, and `blub`. For each, you create a package in the folder `databases`. Each package has a public variable called `Properties` of type `notion.PropertyMetaMap`

Then you just need to create a file with the following content in `databases`:

```go
package main

import (
	"log"

	"github.com/faetools/go-notion-codegen/gen"
	"github.com/faetools/go-notion/pkg/notion"
	"github.com/spf13/afero"
	"github.com/user/myrepo/databases/bar"
	"github.com/user/myrepo/databases/blub"
	"github.com/user/myrepo/databases/foo"
)

//go:generate go run gen.go

func main() {
	fs := afero.NewOsFs()

	for pkgName, props := range map[string]notion.PropertyMetaMap{
		"foo":  foo.Properties,
		"bar":  bar.Properties,
		"blub": blub.Properties,
	} {
		if err := gen.PropertyValues(fs, pkgName, props); err != nil {
			log.Fatal(err)
		}
	}
}
```

Run `go generate ./...` and your code will get generated.

See also [the example](example/databases/).
