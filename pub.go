package main

import (
	"bytes"
	"fmt"
	//"github.com/gohugoio/hugo/markup/tableofcontents"
	"io/ioutil"
	"os"
	"path/filepath"
	//"regexp"
	//"strings"
	//"text/template"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// publishFile() is the heart of this program. It converts
// a Markdown document (with optional TOML at the beginning)
// to HTML.
func (App *App) publishFile(filename string) error {
	var input []byte
	var err error
	// Note the filename
	/// xxx
	App.Page.filePath = filename
	App.Page.filename = filepath.Base(filename)
	// Read the whole Markdown file into memory as a byte slice.
	input, err = ioutil.ReadFile(filename)
	if err != nil {
		return errCode("0102", filename)
	}

	// Output filename
	outfile := replaceExtension(filename, "html")
	// Strip out everything but the filename.
	//base := filepath.Base(outfile)

	// Write everything to a temp file so in case there was an error, the
	// previous HTML file is preserved.
	tmpFile, err := ioutil.TempFile(App.Site.Publish, PRODUCT_NAME+"-tmp-")

	// Translate from Markdown to HTML!
	App.html = App.MdFileToHTMLBuffer(filename, input)
	writeTextFile(tmpFile.Name(), string(App.html))

	// If the write succeeded, rename it to the output file
	// This way if there was an existing HTML file but there was
	// an error in output this time, it doesn't get clobbered.
	if err = os.Rename(tmpFile.Name(), outfile); err != nil {
		return err
	}

	if !fileExists(outfile) {
		QuitError(errCode("0910", outfile))
	}
	App.Verbose("\tCreated file %s", outfile)
	App.fileCount++
	//
	// Success
	return nil
}

// Takes a buffer containing Markdown
// and converts to HTML.
// Doesn't know about front matter.
func (App *App) markdownBufferToBytes(input []byte) []byte {
	fmt.Println("Converting " + string(input))
	autoHeadings := parser.WithAttribute()
	if App.Site.MarkdownOptions.headingIDs == true {
		autoHeadings = parser.WithAutoHeadingID()
	}

	if App.Site.MarkdownOptions.hardWraps == true {
		// TODO: Figure out how to get this damn thing in as an option
		//hardWraps := html.WithHardWraps()
	}

	// TODO: Make the following option: Footnote,
	markdown := goldmark.New(
		goldmark.WithExtensions(extension.GFM, extension.DefinitionList, extension.Footnote,
			highlighting.NewHighlighting(
				highlighting.WithStyle(App.Site.MarkdownOptions.HighlightStyle),
				highlighting.WithFormatOptions(
				//highlighting.WithLineNumbers(),
				),
			)),
		goldmark.WithParserOptions(
			autoHeadings, // parser.WithAutoHeadingID(),
			parser.WithAttribute(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
			/* html.WithHardWraps(), */
			html.WithXHTML(),
		),
	)

	var buf bytes.Buffer
	if err := markdown.Convert(input, &buf); err != nil {
		return []byte{}
	}
	return buf.Bytes()
}

// appendBytes() Appends a byte slice to the byte slice containing the rendered output
func (App *App) appendBytes(b []byte) {
	App.html = append(App.html, b...)
}

// appendStr() Appends a string to the byte slice containing the rendered output
func (App *App) appendStr(s string) {
	App.html = append(App.html, s...)
}

// MdFileToHTMLBuffer() takes a byte slice buffer containing
// a pure Markdown file as input, and returns
// a byte slice containing the file converted to HTML.
// It doesn't know about front matter.
// So it should be preceded by a call to App.parseFrontMatter()
// if there's any possibility that the file contains front matter.
// In the spirit of a browser, it simply returns an empty buffer on error.
func (App *App) MdFileToHTMLBuffer(filename string, input []byte) []byte {
	// Resolve any Go template variables before conversion to HTML.
	fmt.Println("Skipping template execution xxx")
	interp := App.interps(filename, string(input))
	return App.markdownBufferToBytes([]byte(interp))
	//return App.markdownBufferToBytes(input)
}
