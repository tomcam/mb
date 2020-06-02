package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
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

func (App *App) newTheme(from, to string) error {
	return App.copyTheme(from, to, false)
}


// copyThemeDirectory() copies the directory specified by the fully qualified directory name
// from, to the fully qualified  directory name to.
func (App *App) copyThemeDirectory(from, to string) error {
	// Create the target directory
	if err := os.MkdirAll(to, PROJECT_FILE_PERMISSIONS); err != nil {
		return errCode("0905", "Unable to create target theme directory "+to)
	}
	// Copy only 1 level deep.
	// There should be nothing interesting or tricky about this directory. Just
	// markdown files, HTML files, and assets.
	if err := copyDirOnly(from, to); err != nil {
		msg := fmt.Sprintf("Unable to copy from pageType directory %s to new pageType directory %s", from, to)
		return errCode("0906", msg)
	}

	// Success
	return nil
}

// createPageType() is very similar to copyTheme() but
// it creates a new pagetype from an existing one and
// puts it one subdirectory down from the original.
// Given the name of the name of a theme, say, "default", copy it to
// its own subdirectory named pageType (only 1 level deep) and rename some files.
// So if you make a pageType called Blog for the Default theme, it copies all the CSS files,
// for example, but renames default.css to blog.css.
func (App *App) createPageType(theme, pageType string) error {
	// Obtain the fully qualified name of the source
	// theme directory to copy
	source := App.themePath(theme)
	// Generate name of TOML file expected to be there
	tomlFile := App.themeTOMLFilename(theme, source)
	// Check for both these elements.
	if err := App.isTheme(source, tomlFile); err != nil {
		return errCode("1010", source+"  doesn't seem to be a theme")
	}
	// Destination directory is a subdirectory of
	// theme
	dest := filepath.Join(source, pageType)
	/*
		if dirExists(dest) {
			// TODO: Original error code needed
			return errCode("0907", "directory "+dest+" already exists")
		}
	*/
	err := App.copyTheme(theme, dest, true)
	if err != nil {
		return errCode("PREVIOUS", err.Error())
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
func (App *App) copyTheme(from, to string, isChild bool) error {
	//fmt.Printf("copyTheme()\nfrom: %v\nto: %v\nisChild: %v\n",from, to, isChild)
	// Obtain the fully qualified name of the source
	// theme directory to copy
	source := App.themePath(from)
	// Generate name of TOML file expected to be there
	tomlFile := App.themeTOMLFilename(from, source)
	// Check for both these elements.
	if err := App.isTheme(source, tomlFile); err != nil {
		return errCode("1008", "Missing directory ", source, " or TOML file ", tomlFile)
	}

	var dest string
	if !isChild {
		dest = App.themePath(to)
	} else {
		dest = to
	}
	if dirExists(dest) {
		return errCode("0904", "directory "+dest+" already exists")
	}

	if err := App.copyThemeDirectory(source, dest); err != nil {
		return errCode("PREVIOUS", "")
	}
	err := App.updateThemeDirectory(from, dest, to, tomlFile, isChild)
	if err != nil {
		return errCode("PREVIOUS", err.Error())
	}
	// Success
	// copyPageType
	App.Verbose("Created theme " + filepath.Base(dest))
	return nil
}

// themePath() returns the fully qualified pathname of the
// named theme's directory.
func (App *App) themePath(theme string) string {
	return filepath.Join(App.Site.themesPath, theme)
}

// themeTOMLFilename() returns the fully qualified pathname
// of the named theme's expected TOML filename.
func (App *App) themeTOMLFilename(theme, themePath string) string {
	return filepath.Join(themePath, theme+"."+CONFIG_FILE_DEFAULT_EXT)
}

// isTheme() returns true if the fully qualified
// directory pathname exists, and if it contains
// a TOML file by the specified name
func (App *App) isTheme(dir, tomlFile string) error {
	// See if there's a directory of that name.
	if !dirExists(dir) {
		return errCode("0701", dir+" directory not found")
	}

	if !fileExists(tomlFile) {
		return errCode("0115", dir+" theme TOML file "+tomlFile+" is missing")
	}
	// Success
	return nil
}



// updateThemeDirectory() takes a theme directory freshly copied from
// another theme directory, renames the .css file with the theme's name,
// then creates a TOML file.
// from is the bare name of the theme, say, "default".
// if isChild is true:
//   -  to is the fully qualified name of the new theme directory,
//      say, "/Users/tom/html/mysite/themes/mytheme".
// tomlFile is the fully qualified name for the theme named from.
func (App *App) updateThemeDirectory(from, dest, to, tomlFile string, isChild bool) error {
	// Create a toml file for the new theme
	fmt.Printf("UpdateThemeDirectory(%s,%s,%s,%s,%v)\n", from, dest, to, tomlFile, isChild)

	// Parse the original toml file to get its list of stylesheets.
	// Goal is to replace the original theme stylesheet name, say, default.css,
	// with the new theme's style sheet name, say, mytheme.css.
	var p PageType
	if _, err := toml.DecodeFile(tomlFile, &p); err != nil {
		return errCode("0116", fmt.Errorf("Problem reading TOML file %s\n", tomlFile).Error(), err.Error())
	}

	// Get the plain name of the target stylesheet, say, "mynewtheme"
	destFilename := filepath.Base(to)
	var targetTomlFile string
	var targetDir string
	// Figure out the name and location of the toml that describes
	// the theme. If it's a new theme, it would be in something
	// like /themes/mynewtheme/mynewtheme.toml. If it's a pagetype for an existing
	// theme, it would be in something like /themes/mynewtheme/blog/blog.toml
	tomlFilename := destFilename + "." + CONFIG_FILE_DEFAULT_EXT
	if !isChild {
		// It's a new theme
		targetDir = filepath.Join(App.Site.themesPath, to)
	} else {
		// It's a pagetype of an existing theme
		targetDir = filepath.Join(App.Site.themesPath, from, destFilename)
	}
	targetTomlFile = filepath.Join(targetDir, tomlFilename)

	// Obtain the contents of the original TOML file.
	if _, err := toml.DecodeFile(tomlFile, &p); err != nil {
		return errCode("0116",
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
		return errCode("PREVIOUS", "")
	}

	// Now get rid of the previous .toml and .css files
	delToml := filepath.Join(targetDir, from+"."+CONFIG_FILE_DEFAULT_EXT)
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


