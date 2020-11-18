package app

import (
	"bytes"

	"github.com/BurntSushi/toml"
	"github.com/tomcam/mb/pkg/errs"
)

// Starts the (optional)front matter of  a markdown file
var frontMatterDelimiter = []byte("===\n")
var frontMatterDelimiterWin = []byte("===\r\n")

// FrontMatter contains information about each separate page, such as
// its HTML Description tag or pageType.
type FrontMatter struct {
	// isChild is false if it's a root/default theme, and true if it's a child/pagetype declaration
	isChild bool

	// Current theme in use.
	Theme string

	// If "dark", use dark them when available.
	// If "light", use light theme, which is the default.
	// If true, use dark theme if available
	Mode string

	// Current page template in use
	PageTypeName string

	// Generates a Description metatag on output
	Description string

	// Generates a Title tag on output
	Title string

	// Determine whether aside is on the
	// right, left, or none
	Sidebar string

	// List of filenames to exclude
	ExcludeFilenames []string

	// If set to "off", don't execute templates on this page.
	// Used for documentation purposes.
	Templates string

	// User data.
	List interface{}
}

// parseFrontMatter() obtains the TOML front matter,
// if any, from the input buffer.
// Return the starting byte of the markdown
// that follows the front matter.
// A page could just be this:
// # Dude
// hello, world.
//
// But it could also have TOML front matter:
// ===
// description = "Go tutorial"
// ===
// # Dude
// hello, world
//

// The front matter is sandwiched between delimiters
// as shown above. So first see if there's any front
// matter. If there is, return positions needed to
// make it a slice.
func (a *App) parseFrontMatter(filename string, input []byte) (markupStart []byte, err error) {
	// Obtain its front matter
	// The front matter has a list of style sheets needed by
	// the page. Their base directory is the page template's
	// base directory. Obtaining that list means finding
	// the stylesheets for this page's template.

	// First see if there's any front
	// matter. If there is, return positions needed to
	// make it a slice.
	start, stop, markup := getFrontMatter(input, frontMatterDelimiter)
	if start <= 0 && markup >= 0 {
		// No front matter.
		// There is markup.
		return input[markup:], nil
	}
	if start <= 0 && markup < 0 {
		// Not even sure how this happens.
		// No markdown AND no front matter.
		return input, errs.ErrCode("1001", filename)
	}
	decode := string(input[start:stop])
	if _, err := toml.Decode(decode, &a.FrontMatter); err != nil {
		return input[markup:], errs.ErrCode("0101", filename)
	}
	return input[markup:], nil
}

// getFrontMatter() detects whether the input has front matter,
// which looks something like this:
// ===
// description="foo"
// ===
// Where the "===" along with a newline is passed in as the
// delimiter. (It's not hardcoded and could be anything else,
// like "++" or whatever).
//
// Given a byte slice containing an entire markup file, which
// may or may not contain front matter, return the start and
// end positions of the front matter, if any (fstart and fend), and the
// starting position of markdown text, if any.
//
// Front matter must consist of at least 2 lines of delimiters,
// each followed by a newline. It must appear at the very beginning
// of the file. Normally of course the delimiters surround
// actual TOML assignments. Returns -1, -1, 0 if no front matter.
// Returns 0, 0, and markdown start position if front matter is
// empty-that is, both delimiter lines
// exist but there's nothing in between.
//
// Otherwise returns with fstart and fend as the position of the front
// matter text, then the start of markdown in mstart. mstart == -1 if there
// isn't any markdown text.
//
// Assumes Unix-style newlines, but also accounts for Windows.
// This code is fugly, but it seems to work and it's faster than
// reading the first lines of the file as text.
func getFrontMatter(body []byte, delimiter []byte) (fstart, fend, mstart int) {
	max := len(body)
	dl := len(delimiter)
	wdl := len(frontMatterDelimiterWin)
	// Shortest valid delimiter set would be 2 delimiters + a newline
	// after each. If the body is smaller we already know there's
	// no valid front matter.
	if max < dl*2 {
		//fmt.Println("File is too short to have front matter")
		return -1, -1, 0
	}

	// Look for the first delimiter.
	firstDelim := bytes.Index(body, delimiter)
	// Returns == -1 if no delimiter found, >= 0 if found.
	// But also it must start the file, so any value > 0
	// is not valid.
	if firstDelim != 0 {
		// Unix-style delimiter not found.
		// See if the Windows is present.
		firstDelim = bytes.Index(body, frontMatterDelimiterWin)
		if firstDelim != 0 {
			//fmt.Println("Couldn't find front matter delimiter at beginning of file")
			return -1, -1, 0
		} else {
			// fmt.Println("Found first Windows delimiter")
			// Change the length of the delimiter once it's known
			// this is a Windows-style file.
			dl = wdl
		}
	}

	// File is long enough to have both delimiters.
	// It does start with a delimiter.
	// Now look for 2nd one
	secondDelim := bytes.Index(body[dl:], delimiter)
	// fmt.Println("Found first delim. Looking for second start at pos ", dl)
	if secondDelim < 0 {
		// Not found. Maybe it's a windows-style file.
		secondDelim = bytes.Index(body[dl:], frontMatterDelimiterWin)
		if secondDelim < 0 {
			return -1, -1, 0
		}
	}

	// Second delimiter found.
	if secondDelim == 0 {
		// But position 0 means there's nothing between
		// the two delimiters. the front matter is empty.
		//fmt.Println("Front matter delimiters both found, but no content there")
		return 0, 0, max
	}

	start := dl
	end := start + secondDelim
	mstart = start + end
	// Handle rare case of no markdown body at all
	if mstart >= max {
		mstart = -1
	}
	return start, end, mstart
}
