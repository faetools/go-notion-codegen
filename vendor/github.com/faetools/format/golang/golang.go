package golang

import (
	"github.com/faetools/format/format"
	"github.com/pkg/errors"
	"golang.org/x/tools/imports"
	gofumpt "mvdan.cc/gofumpt/format"
)

// FormatOptions defines the format options we're using.
var FormatOptions = gofumpt.Options{
	LangVersion: format.GoVersion.String(),
	ExtraRules:  true,
}

// Format formats golang code.
func Format(filepath string, src []byte) ([]byte, error) {
	res, err := imports.Process(filepath, src, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "running 'imports'")
	}

	return gofumpt.Source(res, FormatOptions)
}
