package markdown

import (
	"bytes"
	"fmt"

	"github.com/faetools/format/writers"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// A NodeRenderer struct is an implementation of renderer.NodeRenderer that renders
// nodes as markdown.
type NodeRenderer struct {
	Config
	additionalOptions    []renderer.Option
	renderFuncsOverrides NodeRendererFuncs
}

// NodeRendererFuncs is a mapping of node rendering functions that should overwrite the default.
type NodeRendererFuncs map[ast.NodeKind]renderer.NodeRendererFunc

// NewNodeRenderer returns a new Renderer with given options.
func NewNodeRenderer(renderFuncsOverrides NodeRendererFuncs, opts ...Option) *NodeRenderer {
	r := &NodeRenderer{
		Config:               NewConfig(),
		renderFuncsOverrides: renderFuncsOverrides,
	}

	for _, opt := range opts {
		opt.SetMarkdownOption(&r.Config)
	}

	return r
}

// RegisterFuncs implements NodeRenderer.RegisterFuncs .
func (r *NodeRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	funcs := NodeRendererFuncs{
		// Blocks.
		ast.KindHeading:          r.renderHeading,
		ast.KindBlockquote:       r.renderBlockquote,
		ast.KindCodeBlock:        r.renderCodeBlock,
		ast.KindFencedCodeBlock:  r.renderFencedCodeBlock,
		ast.KindHTMLBlock:        r.renderHTMLBlock,
		ast.KindList:             r.renderList,
		ast.KindListItem:         r.renderListItem,
		ast.KindParagraph:        r.renderParagraph,
		ast.KindTextBlock:        r.renderTextBlock,
		ast.KindThematicBreak:    r.renderThematicBreak,
		extast.KindStrikethrough: r.renderStrikethrough,

		// Inlines.
		ast.KindAutoLink: r.renderAutoLink,
		ast.KindCodeSpan: r.renderCodeSpan,
		ast.KindEmphasis: r.renderEmphasis,
		ast.KindImage:    r.renderImage,
		ast.KindLink:     r.renderLink,
		ast.KindRawHTML:  r.renderRawHTML,
		ast.KindText:     r.renderText,
		ast.KindString:   r.renderString,
	}

	for kind, f := range r.renderFuncsOverrides {
		funcs[kind] = f
	}

	for kind, f := range funcs {
		reg.Register(kind, f)
	}
}

func (r *NodeRenderer) renderHeading(w util.BufWriter,
	source []byte, node ast.Node, entering bool,
) (ast.WalkStatus, error) {
	if !entering {
		_, _ = w.Write(twoNewLines)
		return ast.WalkContinue, nil
	}

	//nolint:forcetypeassert // function is only called with that
	repeat(w, bHash, node.(*ast.Heading).Level)
	_ = w.WriteByte(bSpace)

	return ast.WalkContinue, nil
}

func (r *NodeRenderer) renderBlockquote(w util.BufWriter,
	source []byte, n ast.Node, entering bool,
) (ast.WalkStatus, error) {
	if !entering {
		_, _ = w.Write(twoNewLines)
		return ast.WalkContinue, nil
	}

	indent := getIndentOfListItem(n)
	if indent > 0 {
		// Separate from the list with a new line.
		_ = w.WriteByte(bNewLine)
	}

	// Write all indents.
	iw := writers.NewIndentWriter(w, indent)

	// Write as blockquote.
	bw := newBlockquoteWriter(iw)

	err := r.renderChildren(bw, source, n)
	if err != nil {
		return 0, err
	}

	return ast.WalkSkipChildren, nil
}

func (r *NodeRenderer) renderCodeBlock(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	indent := getIndentOfListItem(n)
	if indent == 0 {
		// No indent, so we can transform into a fenced code block.
		return r.renderFencedCodeBlock(w, source, n, entering)
	}

	if !entering {
		return ast.WalkContinue, nil
	}

	// Separate from the rest.
	_ = w.WriteByte(bNewLine)

	l := n.Lines().Len()
	lines := make([][]byte, l)

	for i := 0; i < l; i++ {
		line := n.Lines().At(i)
		lines[i] = line.Value(source)
	}

	// Shift all lines to the left.
	for hasSpacePrefix(lines) {
		lines = trimSpacePrefix(lines)
	}

	// Write lines with indent.
	iw := writers.NewIndentWriter(w, indent+1)
	for _, line := range lines {
		_, _ = iw.Write(line)
	}

	return ast.WalkContinue, nil
}

func (r *NodeRenderer) renderFencedCodeBlock(w util.BufWriter,
	source []byte, node ast.Node, entering bool,
) (ast.WalkStatus, error) {
	_, _ = w.Write(codeBlockFence)

	if !entering {
		_, _ = w.Write(twoNewLines)
		return ast.WalkContinue, nil
	}

	lang := getLanguage(source, node)

	// E.g. "```go\n".
	_, _ = w.Write(lang)
	_ = w.WriteByte(bNewLine)

	writeFormattedLines(w, source, node, string(lang))

	return ast.WalkContinue, nil
}

func (r *NodeRenderer) renderHTMLBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		// Write the closure line, if any.
		//nolint:forcetypeassert // function is only called with that
		if cl := node.(*ast.HTMLBlock).ClosureLine; cl.Start > -1 {
			_, _ = w.Write(cl.Value(source))
		}

		_ = w.WriteByte(bNewLine)

		return ast.WalkContinue, nil
	}

	writeFormattedLines(w, source, node, langHTML)

	return ast.WalkContinue, nil
}

func (r *NodeRenderer) renderList(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering && node.Parent().Kind() != ast.KindListItem {
		// Two new lines at the end if not inside another list.
		_, _ = w.Write(twoNewLines)
	}

	return ast.WalkContinue, nil
}

func (r *NodeRenderer) renderListItem(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	indent := getIndentOfListItem(node)

	if !entering {
		if // There is no next list item.
		node.NextSibling() != nil &&
			// No content.
			(node.LastChild() == nil ||
				// We have not rendered a blockquote.
				node.LastChild().Kind() != ast.KindBlockquote) {

			_ = w.WriteByte(bNewLine)
		}

		return ast.WalkSkipChildren, nil
	}

	indentWriter := writers.NewIndentWriter(w, indent)

	p := node.Parent()
	if !p.(*ast.List).IsOrdered() { //nolint:forcetypeassert // function is only called with that
		_, _ = indentWriter.Write(preUListItem)
	} else {
		_, _ = indentWriter.WriteString(fmt.Sprintf(preOListItemFormat, numElement(node)))
	}

	if err := r.renderChildren(indentWriter, source, node); err != nil {
		return 0, err
	}

	return ast.WalkSkipChildren, nil
}

func (r *NodeRenderer) renderParagraph(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	indent := getIndentOfListItem(n)
	partOfList := indent > 0

	if !entering {
		if partOfList && n.NextSibling() != nil {
			_ = w.WriteByte(bNewLine)
		} else if !partOfList {
			_, _ = w.Write(twoNewLines)
		}

		return ast.WalkContinue, nil
	}

	if partOfList && n.PreviousSibling() != nil {
		_ = w.WriteByte(bNewLine)
		iw := writers.NewIndentWriter(w, indent)

		if err := r.renderChildren(iw, source, n); err != nil {
			return 0, err
		}
		_ = w.WriteByte(bNewLine)
	} else {
		if err := r.renderChildren(w, source, n); err != nil {
			return 0, err
		}
	}

	return ast.WalkSkipChildren, nil
}

func (r *NodeRenderer) renderTextBlock(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering && n.NextSibling() != nil {
		_ = w.WriteByte(bNewLine)
	}

	return ast.WalkContinue, nil
}

func (r *NodeRenderer) renderThematicBreak(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		_, _ = w.Write(twoNewLines)
		return ast.WalkContinue, nil
	}

	_, _ = w.Write(thematicBreak)

	return ast.WalkContinue, nil
}

func (r *NodeRenderer) renderAutoLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		_ = w.WriteByte(bGreaterThan)
		return ast.WalkContinue, nil
	}

	_ = w.WriteByte(bLessThan)

	n := node.(*ast.AutoLink)
	_, _ = w.Write(n.URL(source))

	return ast.WalkContinue, nil
}

func (r *NodeRenderer) renderCodeSpan(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	// Check if backticks need to be escaped.
	switch txt := node.FirstChild().(type) {
	case *ast.Text:
		if bytes.Contains(txt.Segment.Value(source), []byte{bGraveAccent}) {
			_ = w.WriteByte(bGraveAccent)
		}
	case *ast.String:
		if bytes.Contains(txt.Value, []byte{bGraveAccent}) {
			_ = w.WriteByte(bGraveAccent)
		}
	}

	if !entering {
		_ = w.WriteByte(bGraveAccent)
		return ast.WalkContinue, nil
	}

	_ = w.WriteByte(bGraveAccent)

	for c := node.FirstChild(); c != nil; c = c.NextSibling() {
		switch txt := c.(type) {
		case *ast.Text:
			_, _ = w.Write(txt.Segment.Value(source))
		case *ast.String:
			_, _ = w.Write(txt.Value)
		}
	}

	return ast.WalkSkipChildren, nil
}

const (
	levelItalics    = 1
	levelBold       = 2
	levelUnderlined = 3
)

var errInvalidEmphLevel = errors.New("invalid emphasis level")

func ancestorHas(n ast.Node, criterion func(ast.Node) bool) bool {
	for p := n.Parent(); p != nil; p = p.Parent() {
		if criterion(p) {
			return true
		}
	}

	return false
}

func selfOrOnlyChildHas(n ast.Node, criterion func(ast.Node) bool) bool {
	if n == nil {
		return false
	}

	return criterion(n) || onlyChildHas(n, criterion)
}

func onlyChildHas(n ast.Node, criterion func(ast.Node) bool) bool {
	if n == nil {
		return false
	}

	c := n.FirstChild()

	if c == nil || c.NextSibling() != nil {
		return false
	}

	return criterion(c)
}

func (r *NodeRenderer) renderEmphasis(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Emphasis)

	sameEmphasis := func(other ast.Node) bool {
		o, ok := other.(*ast.Emphasis)
		return ok && o.Level == n.Level
	}

	// don't write the same emphasis again
	if (entering && onlyChildHas(n.PreviousSibling(), sameEmphasis)) ||
		(!entering && onlyChildHas(n.NextSibling(), sameEmphasis)) ||
		ancestorHas(n, func(ancestor ast.Node) bool {
			return sameEmphasis(ancestor) ||
				selfOrOnlyChildHas(ancestor.PreviousSibling(), sameEmphasis) ||
				selfOrOnlyChildHas(ancestor.NextSibling(), sameEmphasis)
		}) {
		return ast.WalkContinue, nil
	}

	// fmt.Println(n.Level)
	// n.Dump([]byte("foofoofoofoofoofoofoofoofoofoofoofoofoofoofoo"), 0)

	// if s, ok := n.PreviousSibling().(*ast.Emphasis); ok && s.Level == n.Level {
	// 	return ast.WalkContinue, nil
	// }

	// if s, ok := n.NextSibling().(*ast.Emphasis); ok && s.Level == n.Level {
	// 	return ast.WalkContinue, nil
	// }

	if r.Config.Terminal {
		if entering {
			switch n.Level {
			case levelItalics:
				_, _ = w.Write(tItalic)
			case levelBold:
				_, _ = w.Write(tBold)
			case levelUnderlined:
				_, _ = w.Write(tUnderline)
			}
		} else {
			_, _ = w.Write(tReset)
		}

		return ast.WalkContinue, nil
	}

	switch n.Level {
	case levelItalics:
		_ = w.WriteByte(bAsterisk)
	case levelBold:
		_, _ = w.Write(twoAsterisks)
	default:
		return 0, fmt.Errorf("%w (level %d)", errInvalidEmphLevel, n.Level)
	}

	return ast.WalkContinue, nil
}

func (r *NodeRenderer) renderLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_ = w.WriteByte(bLeftSquareBracket)
		return ast.WalkContinue, nil
	}

	n := node.(*ast.Link)
	_, _ = w.Write(linkTransition)
	_, _ = w.Write(n.Destination) // Link.

	if len(n.Title) != 0 {
		_, _ = w.Write(linkTitleStart)
		_, _ = w.Write(n.Title)
		_ = w.WriteByte(bQuotationMark)
	}

	_ = w.WriteByte(bRightParenthesis)

	return ast.WalkContinue, nil
}

func (r *NodeRenderer) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*ast.Image)
	_, _ = w.Write(imageStart)
	_, _ = w.Write(n.Text(source)) // Alt.
	_, _ = w.Write(linkTransition)
	_, _ = w.Write(n.Destination) // Link.

	if len(n.Title) != 0 {
		_, _ = w.Write(linkTitleStart)
		_, _ = w.Write(n.Title)
		_ = w.WriteByte(bQuotationMark)
	}

	_ = w.WriteByte(bRightParenthesis)

	// if part of a document, we want to leave some space before the next element
	if p := n.Parent(); p == nil || p.Kind() == ast.KindDocument {
		_, _ = w.Write(twoNewLines)
	}

	return ast.WalkSkipChildren, nil
}

func (r *NodeRenderer) renderRawHTML(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkSkipChildren, nil
	}

	_, _ = w.Write(node.Text(source))

	//nolint:forcetypeassert // function is only called with that
	segs := node.(*ast.RawHTML).Segments
	for i, l := 0, segs.Len(); i < l; i++ {
		seg := segs.At(i)
		_, _ = w.Write(seg.Value(source))
	}

	return ast.WalkSkipChildren, nil
}

func (r *NodeRenderer) renderString(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.Write(n.(*ast.String).Value)
	} else if s := n.NextSibling(); s != nil && s.Kind() == ast.KindList {
		w.WriteString("\n")
	}

	return ast.WalkContinue, nil
}

func (r *NodeRenderer) renderText(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*ast.Text)
	_, _ = w.Write(n.Segment.Value(source))

	if n.IsRaw() {
		return ast.WalkContinue, nil
	}

	if n.SoftLineBreak() {
		switch {
		case r.HardWraps:
			_, _ = w.Write(twoNewLines)
		default:
			_ = w.WriteByte(bNewLine)
		}

		return ast.WalkContinue, nil
	}

	if n.HardLineBreak() {
		_, _ = w.Write(twoNewLines)
	}

	return ast.WalkContinue, nil
}

func (r *NodeRenderer) renderStrikethrough(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	return ast.WalkContinue, nil
}
