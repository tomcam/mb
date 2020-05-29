package main

/*
import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
	"strings"
)
*/

type Theme struct {
	// Parent root stylesheets get copied to the child automatically
	RootStylesheets []string
	PageType        PageType
}

type PageType struct {
	// Name of the page layout
	// is read-only, and derived from the name of the directory.
	name string

	// The normal theme name is taken from the directory name.
	// This offers more flexbility, since it can be any string.
	Branding string

	// A sentence or two describing this theme, who'd like to use it,
	// and why it will make their lives better.
	Description string

	// List of stylesheets used by this theme. There may be others
	// found in the PageType directory; they aren't included unless
	// listed here.
	// Example [ "foo.css", "bar.css" ]
	Stylesheets []string

	// "Root" stylesheets available to all pagetypes. In the
	// default/root theme directory this list is used for "inheritance"
	// to child pagetypes.
	// It's unlikely
	// you'd want different theme-light.css files for each pagetype,
	// for example. All child pagetype themes copy these over to their
	// theme directorices
	// In the child directories these will be
	// copied by default unless this value is nonempty, in which
	// case only the named stylesheets will be copied over.
	RootStylesheets []string

	// List of assets found in the PageType directory that are not
	// stylesheets
	otherAssets []string

	// Full pathname of the containing directory
	PathName string

	// List of files to exclude from copying to the publish directory
	// One use: when you don't want one or more of the parent style sheets
	// to be copied to the child PageType
	// Don't use wildcards or other Unix patterns normally expanded by the shell.
	Exclude []string

	// List of all areas used by this page
	Nav     area
	Header  area
	Article area
	Footer  area
	Sidebar area
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
