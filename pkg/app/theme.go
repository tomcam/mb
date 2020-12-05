package app

import (
	"fmt"
	"github.com/tomcam/mb/pkg/slices"
	"io/ioutil"
	"github.com/BurntSushi/toml"
	"github.com/tomcam/mb/pkg/defaults"
	"github.com/tomcam/mb/pkg/errs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// TODO: Not needed?
type Supported struct {
	Left_sidebar  string
	right_sidebar string
	dark          string
	light         string
}

type Theme struct {
	// Parent root stylesheets get copied to the child automatically
	//RootStylesheets []string
	PageType  PageType
  // If you base a theme on Textual it can be either
  // wide style (the default) or pillar style. Wide
  // means the site stretches across the full with of
  // the window. Pillar means it has space on the right
  // and left sides. This determines whether a 
  // Textual-based theme gets the ol' Pillar Conversion Kit

	Supported Supported

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

	// "Root" stylesheets available to all pagetypes. In the
	// default/root theme directory this list is used for "inheritance"
	// to child pagetypes.
	// All child pagetype themes copy these over to their
	// theme directories.
	// In the child directories these will be
	// copied by default unless this value is nonempty, in which
	// case only the named stylesheets will be copied over.
	RootStylesheets []string

	// List of stylesheets used by this theme. There may be others
	// found in the PageType directory; they aren't included unless
	// listed here.
	// Example [ "foo.css", "bar.css" ]
	Stylesheets []string

	// List of assets found in the PageType directory that are not
	// stylesheets
	otherAssets []string

	// Full pathname of the containing directory
  // TODO: should probably be lowercase. No reason to export
	PathName string

	// List of files to exclude from copying to the publish directory
	// One use: when you don't want one or more of the parent style sheets
	// to be copied to the child PageType
	// Don't use wildcards or other Unix patterns normally expanded by the shell.
	Exclude []string

	// All parts of the page
	Nav     layoutElement
	Header  layoutElement
	Article layoutElement
	Footer  layoutElement
	Sidebar layoutElement

  // Should be Semver format (https://semver.org)
  // so major-minor-patch + additional labels.
  // Not enforced.
  Version string
}

// layoutElement could be, say, a header:
// html is inline html. filename would be a pathname containing the HTML.
// It defaults to the name of the component, so if it's a nav and
// no filename is specified it assumes nav.html
// Inline HTML would override File if both are specified.
// These correspond directly to the entries in themename.toml
// that look like this:
//  [Header]
//  HTML=""
//  File="header.md"
//
type layoutElement struct {
	// Inline HTML
	HTML string

	// Filename specifying HTML or Markdown
	File string
}

// Return the fuly qualified directory name of the parent theme
func (a *App) parentThemeFullDirectory() string {
	return filepath.Join(a.themesPath, a.FrontMatter.Theme)
}

// Return the fuly qualified filename of the
// parent theme, including the .toml extension.
func (a *App) parentThemeFullPath() string {
	return filepath.Join(a.parentThemeFullDirectory(), a.FrontMatter.Theme+"."+defaults.ConfigFileDefaultExt)
}

// childThemeFullDirectory() Returns the fuly qualified pathname of the
// child theme directory, if any, or "" if none is present.
func (a *App) childThemeFullDirectory() string {
	if a.FrontMatter.PageType!= "" {
    return filepath.Join(a.themesPath, a.FrontMatter.Theme, a.FrontMatter.PageType)
  }
	return  ""
}

// Return the fuly qualified filename of the
// child theme, including the .toml extension.
func (a *App) childThemeFullPath() string {
	return filepath.Join(a.childThemeFullDirectory(), a.FrontMatter.PageType+"."+defaults.ConfigFileDefaultExt)
}


// loadDefaltTheme(): If the site has a preset theme, load it.
// If no theme was specified, use the default.
func (a *App) loadDefaultTheme() {
	if !dirExists(a.themesPath) {
		a.QuitError(errs.ErrCode("1004","theme directory " + a.themesPath))
	}
	// If no theme was specified in the front matter,
	// but one was specified in the site config site.toml,
	// make the one specified in site.toml the theme.
	if a.Site.Theme != "" && a.FrontMatter.Theme == "" {
		a.FrontMatter.Theme = a.Site.Theme
	}
	// If no theme was specified at all, use the Metabuzz default.
	if a.FrontMatter.Theme == "" {
		a.FrontMatter.Theme = strings.ToLower(defaults.DefaultThemeName)
	}

}


func (a *App) loadTheme(parent bool) {
  a.loadDefaultTheme()
  if !parent && a.FrontMatter.PageType== "" {
    // There is no child theme
    return
  }
  var themePath string
	if parent {
    themePath = a.parentThemeFullPath()
  } else {
    themePath = a.childThemeFullPath()
  }
	var p PageType
	if err := readTomlFile(themePath, &p); err != nil {
		a.QuitError(errs.ErrCode("0105", themePath, a.FrontMatter.PageType))
	}
	a.Page.Theme.PageType = p
  if parent {
    a.Page.Theme.PageType.PathName = a.parentThemeFullDirectory()
	  a.FrontMatter.isChild = false
  } else {
    a.Page.Theme.PageType.PathName = a.childThemeFullDirectory()
	  // TODO: Should probably force filename to lowercase
	  a.FrontMatter.isChild = true
  }
}


// pageTypePath() is a utility function to generate the full pathname  of a PageType file
// from a subdirectory name.
func pageTypePath(subDir, themeName string) string {
	path := filepath.Join(subDir, themeName+"."+defaults.ConfigFileDefaultExt)
	return path
}



// file: Name of styleshet in source directory
// from: Fully qualified pathname of source directory
// to: Fully qualified directory name of target directory
func (a *App) copyStyleSheet(file, from, to string) {
	// Pass through if not a local file
	if strings.HasPrefix(strings.ToLower(file), "http") {
		a.appendStr(stylesheetTag(file))
		return
	}

	// Get fully qualified source filename to copy.
	from = filepath.Join(from, file)
  if a.FrontMatter.isChild {
	  //from = filepath.Join(a.childThemeFullDirectory(), file)
    a.Quit("\t\tcopyStyleSheet(): handle case of child theme")
  }

	// Relative path to the publish directory for themes
	to = filepath.Join(to, file)

	if from == to {
		a.QuitError(errs.ErrCode("0922", "from '"+from+"' to '"+to+"'"))
	}

	// Actually copy the style sheet to its destination
	if err := Copy(from, to); err != nil {
		//a.QuitError(errs.ErrCode("0916", "from '"+from+"' to '"+to+"'"))
		a.QuitError(errs.ErrCode("PREVIOUS", err.Error()))

	}
}




func (a *App) updateCopiedThemeDirectory(from, dest string, isChild bool) error {
	// Create a toml file for the new theme

	// Parse the original toml file to get its list of stylesheets.
	// Goal is to replace the original theme stylesheet name, say, default.css,
	// with the new theme's style sheet name, say, mytheme.css.
	var p PageType
  tomlFile := a.themeTOMLFilename(filepath.Base(from))
	if _, err := toml.DecodeFile(tomlFile, &p); err != nil {
		return errs.ErrCode("0128", fmt.Errorf("Problem reading TOML file %s\n", tomlFile).Error(), err.Error())
	}

	// Get the plain name of the target stylesheet, say, "mynewtheme"
	destFilename := filepath.Base(dest)
	var targetTomlFile string
  //var targetDir string
  //targetDir = dest
	// Figure out the name and location of the toml that describes
	// the theme. If it's a new theme, it would be in something
	// like /themes/mynewtheme/mynewtheme.toml. If it's a pagetype for an existing
	// theme, it would be in something like /themes/mynewtheme/blog/blog.toml
	tomlFilename := destFilename + "." + defaults.ConfigFileDefaultExt
	if !isChild {
		// It's a new theme
		//targetDir = filepath.Join(a.Site.themesPath, dest)
		targetTomlFile = filepath.Join(a.themesPath, destFilename, tomlFilename)
	} else {
		// It's a pagetype of an existing theme
		//targetDir = filepath.Join(a.Site.themesPath, dest, from)
		//targetDir = dest
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
	sourceCSSFile := path.Base(from) + ".css"
	// Get the new name to replace it with.
  //targetCSSFile = filepath.Join(dest, destFilename + ".css")
  targetCSSFile = path.Base(dest + ".css")

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
	delToml := filepath.Join(dest, path.Base(tomlFile))
	delCSS := filepath.Join(dest, sourceCSSFile)
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

// copyTheme() copies from the theme directory to, from
// the theme directory from. "from" is specifed only as a file/theme
// name, not a fully qualified pathame, so "wide" for example.
// It copies everything in from, and
// renames the from.toml file in the new theme directory to
// to.toml. to is a fully qualified pathname.
// If isChild is true, then to is actually a child pageType of from,
// so there's different handling.
// It is different from publishTheme(), which is a bit selective.
// For example, it only copies sidebar-left.css or sidebar-right.css
// as needed.
func (a *App) copyTheme(from, to string, isChild bool)  {
  a.Verbose("copyTheme(%v,%v)",from,to)
	// Obtain the fully qualified name of the source
	// theme directory to copy
	source := a.themePath(from)
	// Generate name of TOML file expected to be there
	tomlFile := a.themeTOMLFilename(path.Base(from))
	// Check for both these elements.
	if !a.isTheme(source, tomlFile) {
		a.QuitError(errs.ErrCode("1008", source))
	}

  a.Verbose("\tsource: %v", source)
	var dest string
	if !isChild {
		dest = a.themePath(to)
	} else {
		dest = a.themePath(filepath.Join(from,to))
	}
	if dirExists(dest) {
		a.QuitError(errs.ErrCode("0904", "directory "+dest+" already exists"))
	}
  a.Verbose("\tdest: %v", dest)
  copyDirOnly(source, dest)
	a.updateThemeDirectory(from, dest, tomlFile, isChild)
	// Success
	//a.Verbose("Created theme " + filepath.Base(dest))
	a.Verbose("Created theme " + to + " from " + from + " in " + dest)
}

// themeVersion() returns the Semver version number
// of the theme. 
// tomlFile should be the fully qualified filename
// of the theme file
func (a *App) themeVersion(tomlFile string) string {
	var p PageType
	if _, err := toml.DecodeFile(tomlFile, &p); err != nil {
		a.QuitError(errs.ErrCode("0131", tomlFile))
	}
  return p.Version
}









// newPageType() creates a new pagetype from an existing one, 
// placing it one subdirectory down from the original.
func (a *App) newPageType(theme, pageType string)  {
	// Obtain the fully qualified name of the source
	// theme directory to copy
	source := a.themePath(theme)
	// Generate name of TOML file expected to be there
	tomlFile := a.themeTOMLFilename(theme)
	// Check for both these elements.
	if !a.isTheme(source, tomlFile) {
		a.QuitError(errs.ErrCode("1010", source+"  doesn't seem to be a theme"))
	}
	// Destination directory is a subdirectory of
	// theme
	dest := filepath.Join(source, pageType)
	if dirExists(dest) {
		// TODO: Original error code needed
		a.QuitError(errs.ErrCode("0919", "directory "+dest+" already exists"))
	}
	a.copyTheme(theme, pageType,true)

}


// newTheme() generates a new theme from an old one.
// Equivalent of mb new theme
// This is a parent theme.
func (a *App) newTheme(from, to string) {
  a.Verbose("newTheme(%v,%v)",from,to)
	if from == to {
		a.QuitError(errs.ErrCode("0918", ""))
	}
	if from == "" {
		a.QuitError(errs.ErrCode("1020", ""))
	}
	if to == "" {
		a.QuitError(errs.ErrCode("1017", ""))
	}
  fromThemePath := a.themePath(from)
  toThemePath := a.themePath(to)
  if !dirExists(fromThemePath) {
    a.QuitError(errs.ErrCode("1021", fromThemePath))
  }
  if dirExists(toThemePath) {
    a.QuitError(errs.ErrCode("0952", to))
  }

  a.Verbose("\tcopyDirOnly(%v,%v)", fromThemePath, toThemePath)
  if err := copyDirOnly(fromThemePath, toThemePath); err != nil {
		a.QuitError(errs.ErrCode("PREVIOUS", err.Error()))
  }

  a.Verbose("\tupdateCopiedThemeDirectory(%v,%v,false)", fromThemePath, toThemePath)
	err := a.updateCopiedThemeDirectory(fromThemePath, toThemePath, false)
	if err != nil {
		a.QuitError(errs.ErrCode("PREVIOUS", err.Error()))
	}

  return


  // SAVE: This could refactored and used for publishing, I think
	// Generate name of TOML file expected to be there
	tomlFile := a.themeTOMLFilename(filepath.Base(fromThemePath))
	// Check for both these elements.
	if !fileExists(tomlFile) {
    a.QuitError(errs.ErrCode("1008", filepath.Base(fromThemePath)))
	}

  // Create the destination directory.
	if err := os.MkdirAll(toThemePath, defaults.PublicFilePermissions); err != nil {
		a.QuitError(errs.ErrCode("0409", toThemePath))
	}

  var p PageType

  // Load the TOML file for the source theme
	if err := readTomlFile(tomlFile, &p); err != nil {
		a.QuitError(errs.ErrCode("0127", tomlFile, ))
	}

	for _, file := range p.RootStylesheets {
		// If user has requested a dark theme, then don't copy skin.css
		// to the target. Copy theme-dark.css instead.
		// TODO: Decide whether this should be in root stylesheets and/or regular.
		file = a.getMode(file)
    a.Verbose("\tcopyStyleSheet(%v,%v,%v)\n", file, fromThemePath, toThemePath)
		a.copyStyleSheet(file, fromThemePath, toThemePath)
	}

	for _, file := range p.Stylesheets {
		file = a.getMode(file)
		a.copyStyleSheet(file, fromThemePath, toThemePath)
	}

	err = a.updateCopiedThemeDirectory(fromThemePath, toThemePath, false)
	if err != nil {
		a.QuitError(errs.ErrCode("PREVIOUS", err.Error()))
	}

}


func (a *App) copyRootStylesheets(from, to string) {
  a.Verbose("copyRootStylesheets(%v,%v)",from,to)
  toThemePath := to
  fromThemePath := from
	// Generate name of TOML file expected to be there
	tomlFile := a.themeTOMLFilename(filepath.Base(fromThemePath))

  var p PageType

  a.Verbose("\tTOML file:%v",tomlFile)
  // Load the TOML file for the source theme
	if err := readTomlFile(tomlFile, &p); err != nil {
		a.QuitError(errs.ErrCode("0127", tomlFile, ))
	}

	for _, file := range p.RootStylesheets {
		// If user has requested a dark theme, then don't copy skin.css
		// to the target. Copy theme-dark.css instead.
		// TODO: Decide whether this should be in root stylesheets and/or regular.
		file = a.getMode(file)
    a.Verbose("\tcopyStyleSheet(%v,%v,%v)\n", file, fromThemePath, toThemePath)
		a.copyStyleSheet(file, fromThemePath, toThemePath)
	}

}


// publishThemeDirectory() copies the directory specified 
// by the fully qualified directory name from, 
// to the fully qualified  directory name to.
func (a *App) publishThemeDirectory(from, to string) {
  a.Verbose("\tpublishThemeDirectory(%v,%v)",from,to)
  // Create the destination directory.
	if err := os.MkdirAll(to, defaults.PublicFilePermissions); err != nil {
		a.QuitError(errs.ErrCode("0402", to))
	}
  tomlFile := filepath.Join(from,filepath.Base(from) +".toml")
  a.Verbose("\tversion number of theme: %s", a.themeVersion(tomlFile))

  // xxxx

	// Get the directory listing.
	candidates, err := ioutil.ReadDir(from)
	if err != nil {
		a.QuitError(errs.ErrCode("1023", from, err.Error()))
	}

	// Get list of files in the local directory to exclude from copy
	excludeFromDir := slices.NewSearchInfo(a.FrontMatter.ExcludeFilenames)

  // Copy all files that aren't stylesheets
	for _, file := range candidates {
		filename := file.Name()
		// Don't copy if it's a directory.
		if !file.IsDir() {
			// Don't copy if its extension is on one of the excluded lists.
			if !hasExtension(filename, ".css") &&
				!hasExtensionFrom(filename, defaults.MarkdownExtensions) &&
        !hasExtensionFrom(filename, defaults.ExcludedAssetExtensions) &&
				!excludeFromDir.Contains(filename) &&
				!strings.HasPrefix(filename, ".") {
				// Got the file. Get its fully qualified name.
				copyFrom := filepath.Join(from, filename)
				// Figure out the target directory.
				relDir := relDirFile(a.Site.path, copyFrom)
				// Get the target file's fully qualified filename.
				copyTo := filepath.Join(a.Site.PublishDir, relDir, filename)
        copyTo = filepath.Join(to, filename)
				a.Verbose("\t\tCopy(%s,%s", copyFrom, copyTo)
				if err := Copy(copyFrom, copyTo); err != nil {
					a.QuitError(errs.ErrCode("PREVIOUS", err.Error()))
				}
				// TODO: Get rid of fileExists() when possible
				if !fileExists(copyTo) {
					a.QuitError(errs.ErrCode("PREVIOUS", copyTo))
				}
			}
		}
  }

  // Publish appropriate style sheets

	p := a.Page.Theme.PageType
	a.publishRootStylesheets(to)

	for _, file := range p.Stylesheets {
		file = a.getMode(file)
		a.publishStyleSheet(file, filepath.Join(to, path.Base(file)))
	}
  sidebar := strings.ToLower(a.FrontMatter.Sidebar)
  switch sidebar {
  case "left":
    a.publishStyleSheet("sidebar-left.css", filepath.Join(to, "sidebar-left.css"))
  case "right":
    a.publishStyleSheet("sidebar-right.css", filepath.Join(to, "sidebar-right.css"))
  }
	// responsive.css is always last
	a.publishStyleSheet("responsive.css",filepath.Join(to, "responsive.css"))

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
	tomlFile := a.themeTOMLFilename(theme)
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
 	err := a.publishTheme(theme, dest, true)
	if err != nil {
		return errs.ErrCode("PREVIOUS", err.Error())
	}

	// success
	return nil
}

// publishTheme() copies from the theme directory to, from
// the theme directory from. "from" is specifed only as a file/theme
// name, not a fully qualified pathame, so "wide" for example.
// It copies everything in from, and
// renames the from.toml file in the new theme directory to
// to.toml. to is a fully qualified pathname.
// If isChild is true, then to is actually a child pageType of from,
// so there's different handling.
func (a *App) publishTheme(from, to string, isChild bool) error {
  // xxx
  a.Verbose("publishTheme(%v,%v)",from,to)
	// Obtain the fully qualified name of the source
	// theme directory to copy
	source := a.themePath(from)
	// Generate name of TOML file expected to be there
	tomlFile := a.themeTOMLFilename(path.Base(from))
	// Check for both these elements.
	if !a.isTheme(source, tomlFile) {
		return errs.ErrCode("1008", source)
	}
  // xxxx
  a.Verbose("\tsource: %v", source)
	var dest string
	if !isChild {
		dest = a.themePath(to)
	} else {
		dest = a.themePath(filepath.Join(from,to))
	}
	if dirExists(dest) {
		return errs.ErrCode("0904", "directory "+dest+" already exists")
	}
  a.Verbose("\tdest: %v", dest)
  a.publishThemeDirectory(source, dest)
	a.updateThemeDirectory(from, dest, tomlFile, isChild)
	// Success
	//a.Verbose("Created theme " + filepath.Base(dest))
	a.Verbose("Created theme " + to + " from " + from + " in " + dest)
	return nil
}

// themePath() returns the fully qualified pathname of the
// named theme's directory.
func (a *App) themePath(theme string) string {
	return filepath.Join(a.themesPath, theme)
}

// themeTOMLFilename() returns the fully qualified pathname
// of the named theme's expected TOML filename.
func (a *App) themeTOMLFilename(theme string) string {
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

// updateThemeDirectory() cleans up after a theme has been copied to
// another name (and therefore location). It:
// - Changes the old name to the new name in the new theme TOML file
// - Changes the old theme name CSS file to the new theme name CSS file
// - Deletes the old theme TOML file
// - Deletes the old theme CSS file
func (a *App) updateThemeDirectory(from, dest, tomlFile string, isChild bool) error {
	// Create a toml file for the new theme
  a.Verbose("updateThemeDirectorry(%v,%v,%v,%v)",from,dest,tomlFile,isChild)
  promptString("")
	// Parse the original toml file to get its list of stylesheets.
	// Goal is to replace the original theme stylesheet name, say, default.css,
	// with the new theme's style sheet name, say, mytheme.css.
	var p PageType
	if _, err := toml.DecodeFile(tomlFile, &p); err != nil {
		return errs.ErrCode("0116", fmt.Errorf("Problem reading TOML file %s\n", tomlFile).Error(), err.Error())
	}

	// Get the plain name of the target stylesheet, say, "mynewtheme"
	destFilename := filepath.Base(dest)
	var targetTomlFile string
	// Figure out the name and location of the toml that describes
	// the theme. If it's a new theme, it would be in something
	// like /themes/mynewtheme/mynewtheme.toml. If it's a pagetype for an existing
	// theme, it would be in something like /themes/mynewtheme/blog/blog.toml
	tomlFilename := destFilename + "." + defaults.ConfigFileDefaultExt
	if !isChild {
		// It's a new theme
		targetTomlFile = filepath.Join(a.themesPath, destFilename, tomlFilename)
	} else {
		// It's a pagetype of an existing theme
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
	var targetDir string
	// Figure out the name and location of the toml that describes
	// the theme. If it's a new theme, it would be in something
	// like /themes/mynewtheme/mynewtheme.toml. If it's a pagetype for an existing
	// theme, it would be in something like /themes/mynewtheme/blog/blog.toml
	tomlFilename = destFilename + "." + defaults.ConfigFileDefaultExt
	if !isChild {
		// It's a new theme
		targetDir = filepath.Join(a.Site.themesPath, dest)
		targetTomlFile = filepath.Join(a.themesPath, destFilename, tomlFilename)
	} else {
		// It's a pagetype of an existing theme
		targetDir = dest
		targetTomlFile = filepath.Join(dest, filepath.Base(dest)+"."+defaults.ConfigFileDefaultExt)
	}

	// Now get rid of the previous .toml and .css files
	delToml := filepath.Join(targetDir, from+"."+defaults.ConfigFileDefaultExt)
	delCSS := filepath.Join(targetDir, from+".css")
	// Delete them if they exist. No error is returned if there's a problem.
	// Because I live on the edge, baby.
  promptString("Attemping to delete " + delToml + " and " + delCSS)
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
