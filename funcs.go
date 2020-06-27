package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// article() returns the contents of the Markdown file itself.
// It can only be used from one of the page regions, not inside
// the markdown text, because that would cause a Markdown inception.
func (App *App) article(params ...string) string {
	if len(params) < 1 {
		return string(App.Site.WebPages[App.Page.filePath].html)
	} else {
		return string(App.Site.WebPages[App.Page.filePath].html)
	}
}

// dirNames() returns a directory listing of the specified
// file names in the document's directory
func (App *App) dirNames(params ...string) []string {
	files, err := ioutil.ReadDir(App.Page.dir)
	if err != nil {
		return []string{}
	}
	var ret []string
	for _, file := range files {
		ret = append(ret, file.Name())
	}
	return ret
}

// files() obtains a slice of filenames in the specified
// directory, using a wildcard specified in suffix.
// Example: {{ files "." "*.jpg }}
func (App *App) files(dir, suffix string) []string {
	files, err := filepath.Glob(filepath.Join(dir, suffix))
	if err != nil {
		return []string{}
	} else {
		return files
	}
}

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
"common" for the Site.commonSubDir subdirectory

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
			filename = filepath.Join(App.Site.commonPath, filename)
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

// toc generates a table of contents and includes all headers with a level less
// than or equal level. Level must be 1-6 inclusive.
func (App *App) toc(params ...string) string {
	pcount := len(params)
	var listType string
	var level int
	var err error
	switch pcount {
	case 0:
		{
			level = 6
			listType = "ul"
		}
	case 1:
		{
			level, err = strconv.Atoi(params[0])
			listType = "ul"
		}
	default:
		{
			level, err = strconv.Atoi(params[0])
			listType = params[1]
			if strings.Contains(listType, "ol") {
				listType = "ol"
			} else {
				listType = "ul"
			}
		}
	}

	// Please leave this error code as is
	if err != nil {
		App.QuitError(errCode("1205", err.Error()))
	}
	// Ditto
	if level <= 0 || level > 6 {
		App.QuitError(errCode("1206", params[0]))
	}
	tocs := App.generateTOC(level)
	var s string
	//TODO: I know, I know, string buffers
	for _, t := range tocs {
		s += strings.Repeat("<"+listType+">", t.Level)
		s += "<li>"
		s += fmt.Sprintf(`<a href="#%s">%s</a></li>`, t.ID,t.Header)
		s += strings.Repeat("</"+listType+">", t.Level)
	}
	return s

	b := new(bytes.Buffer)
	b.Grow(256)
	b.WriteString("<ul>")
	for _, t := range tocs {
		b.WriteString("<li>")
		_, _ = fmt.Fprintf(b, `<a href="#%s">`, t.ID)
		b.WriteString(t.Header)
		b.WriteString("</a>")
		b.WriteString("</li>")
	}
	b.WriteString("</ul>")
	return b.String()
}

func (App *App) addTemplateFunctions() {
	App.funcs = template.FuncMap{
		"article":  App.article,
		"dirnames": App.dirNames,
		"files":    App.files,
		"ftime":    App.ftime,
		"hostname": App.hostname,
		"inc":      App.inc,
		"path":     App.path,
		"scode":    App.scode,
		"toc":      App.toc,
	}
}
