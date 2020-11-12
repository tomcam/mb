package app

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/tomcam/mb/pkg/defaults"
	"github.com/tomcam/mb/pkg/errs"
	"os"
	"path/filepath"
	"strings"
)

type Supported struct {
	Left_sidebar  string
	right_sidebar string
	dark          string
	light         string
}

type Theme struct {
	// Parent root stylesheets get copied to the child automatically
	RootStylesheets []string
	PageType        PageType
	Supported       Supported
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
	Nav     pageRegion
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

// themeName() returns the simple name of the current theme.
func (a *App) currThemeName() string {
	theme := strings.ToLower(strings.TrimSpace(a.FrontMatter.Theme))
	if theme == "" {
		theme = a.defaultTheme()
	}
	return theme
}

// themePath() returns the fully qualified path name of the curren theme
func (a *App) currThemePath() string {
	path := filepath.Join(a.themesPath, a.currThemeName())
	return path
}

// loadTheme() copies the theme and pageType, if any, to the Publish directory.
// They're found in App.FrontMatter.Theme and App.FrontMatter.PageType.
func (a *App) loadTheme() {
	themeName := a.currThemeName()
	themeDir := a.currThemePath()
	if !dirExists(themeDir) {
		a.QuitError(errs.ErrCode("1004",
			fmt.Errorf("theme \"%v\" was specified, but couldn't find a theme directory named %v", themeName, themeDir).Error()))
	}

	// Generate the fully qualified name of the TOML file for this theme.
	// TODO: a.themePath()?
	themePath := pageTypePath(themeDir, themeName)
	// First get the parent theme shared assets
	// Temp var because the goal is simply to get the
	// shared assets.
	var p PageType
	if err := a.PageType(themeName, themeDir, themePath, &p); err != nil {
		a.QuitError(errs.ErrCode("0117", themePath, err.Error()))
	}
	a.Page.Theme.RootStylesheets = p.RootStylesheets
	// See if a pagetype has been requested.
	if a.FrontMatter.PageType != "" {
		//if a.FrontMatter.isChild {
		// This is a child theme/page type, not a default/root theme
		a.FrontMatter.isChild = true
		themeDir = filepath.Join(themeDir, a.FrontMatter.PageType)
		themePath = pageTypePath(themeDir, a.FrontMatter.PageType)
		//}
	} else {
		// This is a default/root theme, not a child theme/page type
		a.FrontMatter.isChild = false
		// Try to load the .toml file named after the theme directory.
		themePath = pageTypePath(themeDir, themeName)
	}
	if err := a.PageType(themeName, themeDir, themePath, &a.Page.Theme.PageType); err != nil {
		a.QuitError(errs.ErrCode("0108", fmt.Errorf("Error loading %s", themePath).Error(), err.Error()))
	}
}

// pageTypePath() is a utility function to generate the full pathname  of a PageType file
// from a subdirectory name.
func pageTypePath(subDir, themeName string) string {
	path := filepath.Join(subDir, themeName+"."+defaults.ConfigFileDefaultExt)
	return path
}

// PageType() reads in either the default/anonymous pageType (root of the
// theme directory) or a pageType, named by directory, one level in.
// ThemeDir is the fully qualified path name of the theme directory.
// fullpathName is the fully qualified path name of the .toml file.
func (a *App) PageType(themeName, themeDir, fullPathName string, PageType *PageType) error {
	if err := readTomlFile(fullPathName, PageType); err != nil {
		return errs.ErrCode("0104", fmt.Errorf("Problem reading TOML file %s for theme %s\n", fullPathName, a.FrontMatter.Theme).Error(), err.Error())
	}
	PageType.name = themeName
	PageType.PathName = themeDir
	a.Page.Theme.PageType = *PageType
	// Success
	return nil
}

// newTheme() generates a new theme from an old one.
// Equivalent of mb new theme
func (a *App) newTheme(from, to string) error {
	if from == to {
		a.QuitError(errs.ErrCode("0918", ""))
	}
	if from == "" {
		from = a.defaultTheme()
	}
	if to == "" {
		a.QuitError(errs.ErrCode("1017", ""))
	}
  return a.copyTheme(from, to, false)
}

// copyThemeDirectory() copies the directory specified by the fully qualified directory name
// from, to the fully qualified  directory name to.
func (a *App) copyThemeDirectory(from, to string) error {
	// Create the target directory
	if err := os.MkdirAll(to, defaults.ProjectFilePermissions); err != nil {
		return errs.ErrCode("0905", "Unable to create target theme directory "+to)
	}
	// Copy only 1 level deep.
	// There should be nothing interesting or tricky about this directory. Just
	// markdown files, HTML files, and assets.
	if err := copyDirOnly(from, to); err != nil {
		msg := fmt.Sprintf("Unable to copy from pageType directory %s to new pageType directory %s", from, to)
		return errs.ErrCode("0906", msg)
	}

	// Success
	return nil
}

func (a *App) newPageType(theme, pageType string) error {
	return a.createPageType(theme, pageType)
}

// createPageType() is very similar to copyTheme() but
// it creates a new pagetype from an existing one and
// puts it one subdirectory down from the original.
// Given the name of the name of a theme, say, "default", copy it to
// its own subdirectory named pageType (only 1 level deep) and rename some files.
// So if you make a pageType called Blog for the Default theme, it copies all the CSS files,
// for example, but renames default.css to blog.css.
func (a *App) createPageType(theme, pageType string) error {
	// Obtain the fully qualified name of the source
	// theme directory to copy
	source := a.themePath(theme)
	// Generate name of TOML file expected to be there
	tomlFile := a.themeTOMLFilename(theme, source)
	// Check for both these elements.
	if !a.isTheme(source, tomlFile) {
		return errs.ErrCode("1010", source+"  doesn't seem to be a theme")
	}
	// Destination directory is a subdirectory of
	// theme
	dest := filepath.Join(source, pageType)
	if dirExists(dest) {
		// TODO: Original error code needed
		return errs.ErrCode("0919", "directory "+dest+" already exists")
	}
	err := a.copyTheme(theme, dest, true)
	if err != nil {
		return errs.ErrCode("PREVIOUS", err.Error())
	}

	// success
	return nil
}

// copyTheme() creates a new theme in the theme directory to, from
// the theme directory from. "from" is specifed only as a file/theme
// name, not a fully qulaified pathame, so "wide" for example.
// It copies everything in from, and
// renames the from.toml file in the new theme directory to
// to.toml. to is a fully qualified pathname.
// If isChild is true, then to is actually a child pageType of from,
// so there's different handling.
func (a *App) copyTheme(from, to string, isChild bool) error {
	// Obtain the fully qualified name of the source
	// theme directory to copy
  //fmt.Println("Create theme " + from)
  //fmt.Fprintf(os.Stdout, "Created theme %s from %s",to, from); // xxx
  // xxx
	source := a.themePath(from)
	// Generate name of TOML file expected to be there
	tomlFile := a.themeTOMLFilename(from, source)
	// Check for both these elements.
	if !a.isTheme(source, tomlFile) {
		return errs.ErrCode("1008", from)
	}

	var dest string
	if !isChild {
		dest = a.themePath(to)
	} else {
		dest = to
	}
	if dirExists(dest) {
		return errs.ErrCode("0904", "directory "+dest+" already exists")
	}
	a.Verbose("Creating theme " + dest + to)

	if err := a.copyThemeDirectory(source, dest); err != nil {
		return errs.ErrCode("PREVIOUS", "")
	}
	err := a.updateThemeDirectory(from, dest, to, tomlFile, isChild)
	if err != nil {
		return errs.ErrCode("PREVIOUS", err.Error())
	}
	// Success
	//a.Verbose("Created theme " + filepath.Base(dest))
  fmt.Println("Created theme " + to + " from " + from + " in " + dest)
	return nil
}

// themePath() returns the fully qualified pathname of the
// named theme's directory.
func (a *App) themePath(theme string) string {
	return filepath.Join(a.themesPath, theme)
}

// themeTOMLFilename() returns the fully qualified pathname
// of the named theme's expected TOML filename.
func (a *App) themeTOMLFilename(theme, themePath string) string {
	return filepath.Join(a.themesPath, theme, theme+"."+defaults.ConfigFileDefaultExt)
}

// isTheme() returns true if the fully qualified
// directory pathname exists, and if it contains
// a TOML file by the specified name
func (a *App) isTheme(dir, tomlFile string) bool {
	// See if there's a directory of that name.

	if !dirExists(dir) {
		return false
	}

	if !fileExists(tomlFile) {
		a.QuitError(errs.ErrCode("0115", dir+" theme TOML file "+tomlFile+" is missing"))
	}
	// Success
	return true
}

// updateThemeDirectory() takes a theme directory freshly copied from
// another theme directory, renames the .css file with the theme's name,
// then creates a TOML file.
// from is the bare name of the theme, say, "default".
// if isChild is true, meaning it's a new pagetype, so from a parent theme:
//   -  to is the fully qualified name of the new theme directory,
//      say, "/Users/tom/html/mysite/themes/mytheme".
//  -   dest is the same as to.
//      Yes, this needs to be refactored.
// if isChild is false, meaning it's for a new theme:
//   -   to is a bare name, such as "home"
// tomlFile is the fully qualified name for the theme named from.
func (a *App) updateThemeDirectory(from, dest, to, tomlFile string, isChild bool) error {
	// Create a toml file for the new theme

	// Parse the original toml file to get its list of stylesheets.
	// Goal is to replace the original theme stylesheet name, say, default.css,
	// with the new theme's style sheet name, say, mytheme.css.
	var p PageType
	if _, err := toml.DecodeFile(tomlFile, &p); err != nil {
		return errs.ErrCode("0116", fmt.Errorf("Problem reading TOML file %s\n", tomlFile).Error(), err.Error())
	}

	// Get the plain name of the target stylesheet, say, "mynewtheme"
	destFilename := filepath.Base(to)
	var targetTomlFile string
	var targetDir string
	// Figure out the name and location of the toml that describes
	// the theme. If it's a new theme, it would be in something
	// like /themes/mynewtheme/mynewtheme.toml. If it's a pagetype for an existing
	// theme, it would be in something like /themes/mynewtheme/blog/blog.toml
	tomlFilename := destFilename + "." + defaults.ConfigFileDefaultExt
	if !isChild {
		// It's a new theme
		targetDir = filepath.Join(a.Site.themesPath, to)
		targetTomlFile = filepath.Join(a.themesPath, destFilename, tomlFilename)
	} else {
		// It's a pagetype of an existing theme
		targetDir = filepath.Join(a.Site.themesPath, to, from)
		targetDir = dest
		targetTomlFile = filepath.Join(dest, filepath.Base(dest)+"."+defaults.ConfigFileDefaultExt)
	}
	// Obtain the contents of the original TOML file.
	if _, err := toml.DecodeFile(tomlFile, &p); err != nil {
		return errs.ErrCode("0116",
			fmt.Errorf("Problem reading TOML file %s\n",
				tomlFile).Error(), err.Error())
	}

	var targetCSSFile string

	// Search its list of stylesheets for the old name.
	sourceCSSFile := from + ".css"
	// Get the new name to replace it with.
	targetCSSFile = destFilename + ".css"

	// The TOML file has a declaration along the lines of
	//stylesheets = [ "sizes.css", "theme-light.css", "myoldtheme.css"  ]
	// Look through the old list of stylesheets from the TOML file. Replace the old stylesheet
	// name ("myoldtheme.css" in this example)with the new one.
	newStylesheets := []string{}
	for _, cssFile := range p.Stylesheets {
		if cssFile == sourceCSSFile {
			// Found a matching stylesheet filename. Replace
			// it with the new stylesheet name.
			newStylesheets = append(newStylesheets, targetCSSFile)
		} else {
			// It's a generic file like sizes.css or fonts.css,
			// so copy over unchanged
			newStylesheets = append(newStylesheets, cssFile)
		}

	}
	// Search and replace completed.
	// Replace the old list of stylesheets in the PageType struct.
	p.Stylesheets = newStylesheets

	// Write out the new TOML file, with the search/replaced stylesheet name in the
	// Stylesheets list.
	if err := writeTomlFile(targetTomlFile, &p); err != nil {
		a.QuitError(errs.ErrCode("PREVIOUS", err.Error()))
	}

	// Now get rid of the previous .toml and .css files
	delToml := filepath.Join(targetDir, from+"."+defaults.ConfigFileDefaultExt)
	delCSS := filepath.Join(targetDir, from+".css")
	// Delete them if they exist. No error is returned if there's a problem.
	// Because I live on the edge, baby.
	deleteFileMust(delToml)
	deleteFileMust(delCSS)
	// Create copy of css file
	sourceCSSFile = replaceExtension(tomlFile, "css")
	if isChild {
		d := filepath.Dir(dest)
		targetCSSFile = filepath.Join(d, destFilename, targetCSSFile)
		return Copy(sourceCSSFile, targetCSSFile)

	}
	// It's not a child pageType, it's peer of the original.
	targetCSSFile = replaceExtension(targetTomlFile, "css")
	return Copy(sourceCSSFile, targetCSSFile)
}

// defaultTheme() returns the simple name of
// the theme used to create new pages
// if no theme is specified and to create new themes if no
// source theme is specified.
func (a *App) defaultTheme() string {
	theme := defaults.DefaultThemeName
	if a.Site.Theme != "" {
		theme = a.Site.Theme
	}
	if cfgString("defaulttheme") != "" {
		theme = cfgString("defaulttheme")
	}
	return strings.ToLower(theme)
}
