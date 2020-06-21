package main

import (
	"github.com/google/go-cmp/cmp"
	"github.com/tomcam/mb/pkg/texts"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
	"testing"
)

func Test_extractTOCs(t *testing.T) {
	gm := goldmark.New()
	tests := []struct {
		mdSrc string
		want  []TOCEntry
	}{
		{texts.Dedent(`
		     # h1.1
		     body
		     ## h2.1
		`), []TOCEntry{{"h1.1", 1}, {"h2.1", 2}}},
		{texts.Dedent(`
		     # h1.1 *foo* bar
		     body
		     ## h2.1
		     ## h2.2
		     ### h3.1
		     # h1.2
		`), []TOCEntry{{"h1.1 <em>foo</em> bar", 1}, {"h2.1", 2},
			{"h2.2", 2}, {"h3.1", 3}, {"h1.2", 1}}},
	}
	for _, tt := range tests {
		t.Run(tt.mdSrc, func(t *testing.T) {
			node := gm.Parser().Parse(text.NewReader([]byte(tt.mdSrc)))
			got, err := extractTOCs(gm.Renderer(), node, []byte(tt.mdSrc))
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("extractTOCs() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
