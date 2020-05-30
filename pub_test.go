package main

import (
	"fmt"
	"strings"
	"testing"
)

//func TestMdFileToHTMLBuffer(t *testing.T) {
func TestConvert(t *testing.T) {
	var tests = []struct {
		have string
		want string
	}{
    /* ==== */
		{`# hi`,
			`<h1 id="hi">hi</h1>`},

    /* ==== */
    {`===
theme="foo"
===
hi`,
			`<p>hi</p>`},
    /* ==== */

	}
	App := newDefaultApp()
	for each, tt := range tests {
		//t.Run("TestMdFileToHTMLBuffer", func(t *testing.T) {
		t.Run("TestConvert", func(t *testing.T) {
      // TODO: Catch that error code
			ans, _ := App.Convert("unitTest", []byte(tt.have))
			ans = []byte(strings.Trim(string(ans), "\n"))
			fmt.Printf("test %v\n", each)
			if string(ans) != tt.want {
				t.Errorf("got %s, want %s", ans, tt.want)
			}
		})
	}
}
