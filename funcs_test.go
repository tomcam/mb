package main

import (
	"github.com/tomcam/mb/htmls"
	"github.com/tomcam/mb/pkg/texts"
	"strconv"
	"testing"
)

func TestApp_toc(t *testing.T) {
	tests := []struct {
		mdSrc string
		level int
		want  string
	}{
		{
			texts.Dedent(`
		     # h1.1
		     body
		     ## h2.1
		     ## h3.3
		`),
			2,
			texts.Dedent(`
          <ul>
            <li><a href="#h11">h1.1</a></li>
            <li><a href="#h21">h2.1</a></li>
            <li><a href="#h33">h3.3</a></li>
          </ul>
			`)},
	}
	for _, tt := range tests {
		t.Run(tt.mdSrc, func(t *testing.T) {
			app := newApp(tt.mdSrc)
			app.Site.MarkdownOptions.headingIDs = true
			level := strconv.Itoa(tt.level)
			tocHTML := app.toc(level)
			if diff, err := htmls.DiffStrings(tt.want, tocHTML); err != nil {
				t.Error(err)
			} else if diff != "" {
				t.Errorf("tocs() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
