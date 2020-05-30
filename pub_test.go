package main

import (
	"testing"
  "fmt"
  "strings"
)

func TestMdFileToHTMLBuffer(t *testing.T) {
	var tests = []struct {
		have string
    want string
	}{
		{`# hi`, `<h1 id="hi">hi</h1>`},
		{`# foo`, `<h1 id="hi">hi</h1>`},
	}
	App := newDefaultApp()
	for _, tt := range tests {
		t.Run("TestMdFileToHTMLBuffer", func(t *testing.T) {
			ans := App.MdFileToHTMLBuffer("foo", []byte(tt.have))
      ans = []byte(strings.Trim(string(ans),"\n"))
      fmt.Println("******** html: " + string(ans))
			if string(ans) != tt.want {
				t.Errorf("got %s, want %s", ans, tt.want)
			}
		})
	}
}
