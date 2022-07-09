package markdown

import (
	"bytes"
	"fmt"

	"github.com/faetools/format/writers"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark"

	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
	yamlV2 "gopkg.in/yaml.v2"
	yamlV3 "gopkg.in/yaml.v3"
)

const (
	startYAMLFrontmatter = "---\n"
	endYAMLFrontmatter   = "---\n\n"
)

var myParser = goldmark.New(
	goldmark.WithExtensions(meta.Meta),
	goldmark.WithExtensions(extension.GFM)).Parser()

// Render renders a given parsed document.
func Render(metaData interface{}, src []byte, doc ast.Node,
	renderFuncsOverrides NodeRendererFuncs, options ...renderer.Option,
) ([]byte, error) {
	b := &bytes.Buffer{}

	if err := renderFrontMatter(b, metaData); err != nil {
		return nil, err
	}

	w := writers.NewTrimWriter(b, sNewLine)

	nr := NewNodeRenderer(renderFuncsOverrides)
	nr.additionalOptions = options

	opts := append([]renderer.Option{
		renderer.WithNodeRenderers(util.Prioritized(nr, 0)),
	}, options...)

	if err := renderer.NewRenderer(opts...).Render(w, src, doc); err != nil {
		return nil, errors.Wrap(err, "rendering markdown")
	}

	// We trimmed all new lines but we still want one at the end.
	if _, err := b.Write(newLine); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func renderFrontMatter(w writers.Writer, fm interface{}) error {
	if ctx, ok := fm.(parser.Context); ok {
		fm = meta.GetItems(ctx)
	}

	switch m := fm.(type) {
	case nil:
		return nil
	case yamlV2.MapSlice:
		if len(m) == 0 {
			return nil
		}

		if _, err := w.WriteString(startYAMLFrontmatter); err != nil {
			return err
		}

		if err := yamlV2.NewEncoder(w).Encode(m); err != nil {
			return fmt.Errorf("encoding %T: %w", m, err)
		}
	default:
		if _, err := w.WriteString(startYAMLFrontmatter); err != nil {
			return err
		}

		if err := yamlV3.NewEncoder(w).Encode(m); err != nil {
			return fmt.Errorf("encoding %T: %w", m, err)
		}
	}

	if _, err := w.WriteString(endYAMLFrontmatter); err != nil {
		return err
	}

	return nil
}
