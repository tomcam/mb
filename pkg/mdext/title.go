package mdext

import (
	"fmt"
	"github.com/yuin/goldmark/ast"
)

// InferTitle finds the most-likely title for the markdown document.
// Returns the empty string if no headers exist in the document.
// This is useful when a title is not explicitly given in the YAML markdown.
// The title is the text of the first header where the headers are ordered by
// level then by order of appearance in the document. Meaning, we'll use the
// first H2 header even if an H3 header appears before the H2 header.
func InferTitle(root ast.Node, mdSrc []byte) (string, error) {
	headers := make([]*ast.Heading, 6) // there's a max of 6 heading levels
	err := ast.Walk(root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering || n.Kind() != ast.KindHeading {
			return ast.WalkContinue, nil
		}

		h := n.(*ast.Heading)
		if h.Level < 1 || 6 < h.Level {
			return ast.WalkStop, fmt.Errorf("invalid heading level: %d", h.Level)
		}
		// Subtract 1 to convert from 1-based headers to 0-based slices.
		if headers[h.Level-1] == nil {
			headers[h.Level-1] = h
		}
		// Skip children because headings can't contain other headings.
		return ast.WalkSkipChildren, nil
	})
	if err != nil {
		return "", err
	}

	// Pick the first header that exists by order.
	for _, header := range headers {
		if header == nil {
			continue
		}
		return string(header.Text(mdSrc)), nil
	}
	return "", nil
}
