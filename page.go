package main

// The table of contents consists of an array of these
// dudes.
type TOCEntry struct {
	// The text of the header
	header string
	// Its level 1-6 for h1-h6
	level int
}

// type Page contains read-only information about the Markdown page currently
// being processed.
type Page struct {
	// Currently loaded theme
	Theme Theme

	// Table of contents for this page
	TOC []TOCEntry

	// Name of theme used for code highlighting
	// Currently using Chroma:
	// https://github.com/alecthomas/chroma/tree/master/styles
	CodeTheme string
	// Begins the HTML document, assigns the language,
	// and opens the head tag
	startHTML string

	// Contennts of <title> tag
	titleTag string

	// Optional <head> tags, like Google Analytics
	headTags string

	// Optional, additional header entries
	headerFiles string

	// Contents of <description> tage
	descriptionTag string

	// Page being rendered
	html []byte

	// Position in buffer where Markdown text begins, after
	// the optional front matter
	markdownStart []byte

	// Content of the article md file itself converted to HTML,
	// obviously without header, nav, footer, aside etc.
	Article []byte

	// Fully qualified filename.
	filePath string

	// Filename
	filename string

	// Current directory
	dir string

	// Relative directory
	Path string

	// List of assets to be published; any graphics files, etc. in
	// the local directory
	assets []string
}

// Area could be, say, a header:
// html is inline html. filename would be a pathname containing the HTML.
// It defaults to the name of the component, so if it's a nav and
// no filename is specified it assumes nav.html
// Inline HTML would override File if both are specified.
type area struct {
	// Inline HTML
	HTML string

	// Filename specifying HTML or Markdown
	File string
}
