package main

// type Page contains read-only information about the Markdown page currently
// being processed.
type Page struct {
	// Name of theme used for code highlighting
	// Currently using Chroma:
	// https://github.com/alecthomas/chroma/tree/master/styles
	CodeTheme string

  // Page being rendered
	html []byte

LHeadTag[]byte

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
