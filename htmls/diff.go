// package htmls contains utilities for comparing and constructing HTML.
package htmls

import (
	"bytes"
	"io"
	"strings"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/net/html"
)

// DiffStrings returns the diff between the two normalized HTML fragments.
func DiffStrings(x, y string) (string, error) {
	return Diff(strings.NewReader(x), strings.NewReader(y))
}

// Diff returns the diff between the two normalized HTML fragments.
func Diff(got, want io.Reader) (string, error) {
	frag1, err := parseFragment(want)
	if err != nil {
		return "", err
	}

	frag2, err := parseFragment(got)
	if err != nil {
		return "", err
	}

	dump1 := renderNodes(frag1)
	dump2 := renderNodes(frag2)
	return cmp.Diff(dump1, dump2), nil
}

// parseFragment parses a normalized version of an HTML node.
func parseFragment(r io.Reader) ([]*html.Node, error) {
	testCtx := &html.Node{Type: html.ElementNode, Data: "normalizedFragment"}
	nodes, err := html.ParseFragment(r, testCtx)
	if err != nil {
		return nil, err
	}
	return normalizeNodes(nodes), nil
}

func normalizeNodes(nodes []*html.Node) []*html.Node {
	ns := make([]*html.Node, 0, len(nodes))
	for _, node := range nodes {
		if isEmptyNode(node) {
			continue
		}
		normalizeNode(node)
		ns = append(ns, node)
	}
	return ns
}

func normalizeNode(node *html.Node) {
	switch node.Type {
	case html.TextNode:
		node.Data = strings.TrimSpace(node.Data)
		if isEmptyNode(node) {
			if node.Parent == nil {
				return
			}

			p := node.Parent
			p.RemoveChild(node)
		}
	case html.ElementNode:
		cur := node.FirstChild
		for cur != nil {
			next := cur.NextSibling
			normalizeNode(cur)
			cur = next
		}
	}
}

func isEmptyNode(node *html.Node) bool {
	return strings.TrimSpace(node.Data) == ""
}

// renderNodes prints a string representation of HTML nodes.
func renderNodes(nodes []*html.Node) string {
	b := new(bytes.Buffer)
	for i, node := range nodes {
		_ = html.Render(b, node)
		if i < len(nodes)-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}
