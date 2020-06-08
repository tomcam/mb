package main

import (
	"bytes"
	"strings"
	"text/template"
)

// Resolve template variables
// input is an HTML file that includes entities like {{.FrontMatter.Description}}
// Replace with the appropriate values in generated output.
// The filename is passed in because it
// produces an accurate location of any
// source file parsing errors that occur.
func (App *App) interps(filename string, input string) string {
	if strings.ToLower(App.FrontMatter.Templates) != "off" {
		return App.execute(filename, input, App.funcs)
	}
	return  input
}

// execute() parses a Go template, then executes it against HTML/template source.
// Returns a string containing the result.
func (App *App) execute(templateName string, tpl string, funcMap template.FuncMap) string {
	var t *template.Template
	var err error
	if t, err = template.New(templateName).Funcs(funcMap).Parse(tpl); err != nil {
		App.QuitError(errCode("0917", err.Error()))
	}
	var b bytes.Buffer
	err = t.ExecuteTemplate(&b, templateName, App)
	if err != nil {
		App.QuitError(errCode("1204", err.Error()))
	}
	return b.String()
}
