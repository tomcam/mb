package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestMdFileToHTMLBuffer(t *testing.T) {
	var tests = []struct {
		have string
		want string
	}{
		{`# hi`,
			`<h1 id="hi">hi</h1>`},

		{`# h1
hello, world.`,

			`<h1 id="h1">h1</h1>
<p>hello, world.</p>`},

	}
	App := newDefaultApp()
	for each, tt := range tests {
		t.Run("TestMdFileToHTMLBuffer", func(t *testing.T) {
			ans := App.MdFileToHTMLBuffer("unitTest", []byte(tt.have))
			ans = []byte(strings.Trim(string(ans), "\n"))
			fmt.Printf("test %v\n", each)
			if string(ans) != tt.want {
				t.Errorf("got %s, want %s", ans, tt.want)
			}
		})
	}
}
