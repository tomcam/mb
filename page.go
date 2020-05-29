package main

// type Page contains read-only information about the Markdown page currently
// being processed.
type Page struct {
	// Name of theme used for code highlighting
	// Currently using Chroma:
	// https://github.com/alecthomas/chroma/tree/master/styles
	CodeTheme string

	// Content of the article itself, obviously without
	// header, nav, footer, aside etc.
	Content []byte

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
