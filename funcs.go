package main

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ftime() returns the current, local, formatted time.
// Can pass in a formatting string
// https://golang.org/pkg/time/#Time.Format
func (App *App) ftime(param ...string) string {
	var ref = "Mon Jan 2 15:04:05 -0700 MST 2006"
	var format string

	if len(param) < 1 {
		format = ref
	} else {
		format = param[0]
	}
	t := time.Now()
	return t.Format(format)
}

// hostname() returns the name of the machine
// this code is running on
func (App *App) hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	} else {
		return hostname
	}
}

// inc inserts the named file into the current Markdown file.
/* Treats it as a Go template, so either HTML or Markdown
work fine.
The location of the file appears first, before a pipe character.
It can be one of:

"article" for the current markdown file's directory
"common" for the Site.CommonSubDir subdirectory

So it might look like :

  inc "articles|kitchen.md"

*/
func (App *App) inc(filename string) template.HTML {

	// Read the HTML file into a byte slice.
	var input []byte
	var err error
	if filename == "" {
		return template.HTML("")
	}
	parsed := strings.Split(filename, "|")
	// If nothing specified, look in article directory
	if len(parsed) < 2 {
		filename = filepath.Join(App.Page.dir, parsed[0])
	} else {
		location := parsed[0]
		filename = parsed[1]

		switch strings.ToLower(location) {
		case "article":
			filename = filepath.Join(App.Page.dir, filename)
		case "common":
			filename = filepath.Join(App.Site.commonDir, filename)
		default:
			App.QuitError(errCode("0119", location))
		}
	}
	if !fileExists(filename) {
		App.QuitError(errCode("0120", filename))
	}

	input, err = ioutil.ReadFile(filename)
	if err != nil {
		App.QuitError(errCode("0121", filename))
	}

	// Apply the template to it.
	// The one function missing from fewerFuncs is shortcode() itself.
	s := App.execute(filename, string(input), App.fewerFuncs)
	return template.HTML(s)
}

// path() returns the current markdown document's directory
func (App *App) path() string {
	return App.Page.Path
}

// scode() provides a not-very-good shortcode feature. Can't figure
// out how to do a better job considering a Go template function
// can take only a map, but you can't pass map literals.
// You need to pass it a map with a key named "filename"
// that matches a file in ".scodes". Currently the
// only thing that works with Javascript is youtube.html
/*
   Example:
   ====
   [List]
   youtube = { filename="youtube.html", id = "dQw4w9WgXcQ" }
   ===
   ## Youtube?
   {{ scode .FrontMatter.List.youtube }}


*/
func (App *App) scode(params map[string]interface{}) template.HTML {
	filename, ok := params["filename"].(string)
	if !ok {
		return template.HTML("filename missing")
	}
	var input []byte
	var err error
	if len(params) < 1 {
		return ("ERROR0")
	}

	// If no extension specified assume HTML
	if filepath.Ext(filename) == "" {
		filename = replaceExtension(filename, "html")
	}

	// Find that file in the shortcode file directory
	filename = filepath.Join(App.Site.sCodePath, filename)

	if !fileExists(filename) {
		App.QuitError(errCode("0122", filename))
	}

	input, err = ioutil.ReadFile(filename)
	if err != nil {
		App.QuitError(errCode("0123", filename))
	}

	// Apply the template to it.
	// The one function missing from fewerFuncs is shortcode() itself.
	s := App.execute(filename, string(input), App.fewerFuncs)
	return template.HTML(s)
}
