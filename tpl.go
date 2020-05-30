package main

import (
	"bytes"
	"text/template"
)

// Resolve template variables
// input is an HTML file that includes entities like {{.FrontMatter.Description}}
// Replace with the appropriate values in generated output.
// The filename is passed in because it
// produces an accurate location of any
// source file parsing errors that occur.
func (App *App) interps(filename string, input string) string {
	var s string
	var err error
	if s, err = App.execute(filename, input, App.funcs); err != nil {
		return ""
	}
	return s
}

// execute() parses a Go template, then executes it against HTML/template source.
// Returns a string containing the result.
func (App *App) execute(templateName string, tpl string, funcMap template.FuncMap) (buf string, err error) {
	var t *template.Template
	if t, err = template.New(templateName).Funcs(funcMap).Parse(tpl); err != nil {
		return "", err
	}
	var b bytes.Buffer
	err = t.ExecuteTemplate(&b, templateName, App)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}
