package main

import (
	"bytes"
	"fmt"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
)

// generateTOC reads the Markdown source and returns a string
// array containing each header and its level in the document.
func (App *App) generateTOC(level int) {
	tocs, err := extractTOCs(App.goldmark.Renderer(), App.Page.mdNode, App.Page.markdownStart)
	if err != nil {
		App.QuitError(errCode("0901", err.Error()))
	}
	for _, toc := range tocs {
		if toc.Level <= level {
			App.Page.TOC = append(App.Page.TOC, toc)
		}
	}
}

// extractTOCs parses returns TOC entries corresponding to each header in the
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
		tocs[i] = TOCEntry{
			Header: buf.String(),
			Level:  node.Level,
		}
		buf.Reset()
	}
	return tocs, nil
}
