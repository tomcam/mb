package app

import (
	"bytes"
	"github.com/tomcam/mb/pkg/errs"
	"strings"
	"text/template"
)

// Resolve template variables
// input is an HTML file that includes entities like {{.FrontMatter.Description}}
// Replace with the appropriate values in generated output.
// The filename is passed in because it
// produces an accurate location of any
// source file parsing errors that occur.
func (a *App) interps(filename string, input string) string {
	if strings.ToLower(a.FrontMatter.Templates) != "off" {
		return a.execute(filename, input, a.funcs)
	}
	return input
}

// execute() parses a Go template, then executes it against HTML/template source.
// Returns a string containing the result.
func (a *App) execute(templateName string, tpl string, funcMap template.FuncMap) string {
	var t *template.Template
	var err error
	if t, err = template.New(templateName).Funcs(funcMap).Parse(tpl); err != nil {
		a.QuitError(errs.ErrCode("0917", err.Error()))
	}
	var b bytes.Buffer
	err = t.ExecuteTemplate(&b, templateName, a)
	if err != nil {
		a.QuitError(errs.ErrCode("1204", err.Error()))
		//a.QuitError(ErrCode("PREVIOUS",""))
	}
	return b.String()
}
