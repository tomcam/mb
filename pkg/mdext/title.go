package mdext

import "github.com/yuin/goldmark/ast"

// InferTitle finds the most-likely title for the markdown document.
// Returns the empty string if no headers exist in the document.
// This is useful when a title is not explicitly given in the YAML markdown.
// The title is the text of the first header where the headers are ordered by
// level then by order of appearance in the document. Meaning, we'll use the
// first H2 header even if an H3 header appears before the H2 header.
func InferTitle(root ast.Node, mdSrc []byte) string {
	headers := make([]*ast.Heading, 6) // there's a max of 6 heading levels
	// Ignoring errors because our visit function doesn't error.
	_ = ast.Walk(root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering || n.Kind() != ast.KindHeading {
			return ast.WalkContinue, nil
		}

		h := n.(*ast.Heading)
		if headers[h.Level] == nil {
			headers[h.Level] = h
		}
		// Skip children because headings can't contain other headings.
		return ast.WalkSkipChildren, nil
	})

	for _, header := range headers {
		if header == nil {
			continue
		}
		return string(header.Text(mdSrc))
	}
	return ""
}
