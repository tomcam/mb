package main

import (
	"fmt"
	//"github.com/BurntSushi/toml"
	//"os"
	"path/filepath"
	"strings"
)

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

  // All parts of the page
	Nav    pageRegion
	Header  pageRegion
	Article pageRegion
	Footer  pageRegion
	Sidebar pageRegion
}

// pageRegion could be, say, a header:
// html is inline html. filename would be a pathname containing the HTML.
// It defaults to the name of the component, so if it's a nav and
// no filename is specified it assumes nav.html
// Inline HTML would override File if both are specified.
type pageRegion struct {
	// Inline HTML
	HTML string

	// Filename specifying HTML or Markdown
	File string
}



// loadTheme() copies the theme and pageType, if any, to the Publish directory.
// They're found in App.FrontMatter.Theme and App.FrontMatter.PageType.
func (App *App) loadTheme() {
	themeName := strings.ToLower(strings.TrimSpace(App.FrontMatter.Theme))
	themeDir := filepath.Join(App.Site.themesPath, themeName)
	if !dirExists(themeDir) {
		  App.QuitError(errCode("1004",
			fmt.Errorf("theme \"%v\" was specified, but couldn't find a directory named %v", App.FrontMatter.Theme, themeDir).Error()))
	}

	// Generate the fully qualified name of the TOML file for this theme.
	// TODO: App.themePath()?
	themePath := pageTypePath(themeDir, themeName)

	// First get the parent theme shared assets
	// Temp var because the goal is simply to get the
	// shared assets.
	var p PageType
	if err := App.PageType(themeName, themeDir, themePath, &p); err != nil {
		App.QuitError(errCode("0117", themePath, err.Error()))
	}
	App.Page.Theme.RootStylesheets = p.RootStylesheets
	// See if a pagetype has been requested.
	if App.FrontMatter.PageType != "" {
		//if App.FrontMatter.isChild {
		//fmt.Println("loadTheme(), PageType", App.FrontMatter.PageType)
		// This is a child theme/page type, not a default/root theme
		App.FrontMatter.isChild = true
		themeDir = filepath.Join(themeDir, App.FrontMatter.PageType)
		themePath = pageTypePath(themeDir, App.FrontMatter.PageType)

	} else {
		//fmt.Println("loadTheme(), root theme")
		// This is a default/root theme, not a child theme/page type
		App.FrontMatter.isChild = false
		// Try to load the .toml file named after the theme directory.
		themePath = pageTypePath(themeDir, themeName)
		//fmt.Println("Hope we inherited", App.Page.Theme.RootStylesheets)
	}
	if err := App.PageType(themeName, themeDir, themePath, &App.Page.Theme.PageType); err != nil {
		App.QuitError(errCode("0108", fmt.Errorf("Error loading %s", themePath).Error(), err.Error()))
	}
}

// pageTypePath() is a utility function to generate the full pathname  of a PageType file
// from a subdirectory name.
func pageTypePath(subDir, themeName string) string {
	return filepath.Join(subDir, themeName+"."+CONFIG_FILE_DEFAULT_EXT)
}

// PageType() reads in either the default/anonymous pageType (root of the
// theme directory) or a pageType, named by directory, one level in.
// themeDir is the fully qualified path name of the theme directory.
// fullpathName is the fully qualified path name of the .toml file.
func (App *App) PageType(themeName, themeDir, fullPathName string, PageType *PageType) error {
	if err := readTomlFile(fullPathName, PageType); err != nil {
		return errCode("0104", fmt.Errorf("Problem reading TOML file %s for theme %s\n", fullPathName, App.FrontMatter.Theme).Error(), err.Error())
	}
	PageType.name = themeName
	PageType.PathName = themeDir
	App.Page.Theme.PageType = *PageType
	// Success
	return nil
}



