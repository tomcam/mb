package main

import (
	"bytes"
	"fmt"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
)

// generateTOC reads the Markdown source and returns a slice of TOC entries
// corresponding to each header less than or equal to level.
func (App *App) generateTOC(level int) []TOCEntry {
	node := App.markdownAST(App.Page.markdownStart)
	tocs, err := extractTOCs(App.newGoldmark().Renderer(), node, App.Page.markdownStart)
	if err != nil {
		App.QuitError(errCode("0901", err.Error()))
	}
	leveledTocs := make([]TOCEntry, 0, len(tocs))
	for _, toc := range tocs {
		if toc.Level <= level {
			leveledTocs = append(leveledTocs, toc)
		}
	}
	return leveledTocs
}

// extractTOCs returns TOC entries corresponding to each header in the markdown
// document ordered by appearance.
func extractTOCs(rn renderer.Renderer, node ast.Node, mdSrc []byte) ([]TOCEntry, error) {
	const numTocs = 8 // assume a reasonable amount of headers per md document
	tocNodes := make([]*ast.Heading, 0, numTocs)
	// Safe to ignore error because only our walk function can error and we don't
	// error.
	_ = ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering || n.Kind() != ast.KindHeading {
			return ast.WalkContinue, nil
		}
		h := n.(*ast.Heading)
		tocNodes = append(tocNodes, h)
		return ast.WalkContinue, nil
	})

	// Render each TOC into an HTML string in a TOCEntry.
	tocs := make([]TOCEntry, len(tocNodes))
	buf := new(bytes.Buffer)
	const headerSize = 64
	buf.Grow(headerSize)
	for i, node := range tocNodes {
		// Render the children of the heading because we don't want the enclosing
		// <h1> tags.
		for c := node.FirstChild(); c != nil; c = c.NextSibling() {
			if err := rn.Render(buf, mdSrc, c); err != nil {
				return nil, fmt.Errorf("extract TOC render: %w", err)
			}
		}
		id := ""
		if attr, ok := node.AttributeString("id"); ok {
			id = string(attr.([]byte))
		}
		tocs[i] = TOCEntry{
			ID:     id,
			Header: buf.String(),
			Level:  node.Level,
		}
		buf.Reset()
	}
	return tocs, nil
}
