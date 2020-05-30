package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// Site contains configuration specific to each site, such as
// its title, publish directory, and branding string.
type Site struct {
	// Full path of the directory where the site source is.
	path string

	// Subdirectry + filename site config file
	configFilePath string

	// Full path of shortcode dir for this project. It's computed
	// at runtime using SCodeDir, also in this struct.
	sCodePath string

	// List of all directories in the site
	dirs map[string]mdOptions

	// Site's project name, so it's a filename.
	// It's an identifier so it should be in slug format:
	// Preferably just alphanumerics, underline or hyphen, and
	// no spaces, for example, 'my-project'
	Name string

	// Site's branding, any string. So, for example,
	// if the Name is 'my-project' this might be
	// 'My Insanely Cool Project"
	Branding string

	// Full pathname of common directory. Derived from CommonSubDir
	commonDir string

	// Directory to share reusable text.
	// Use the computed value Site.commonDir for hte full path.
	CommonSubDir string

	// Language tag for html lang=
	Language string

	// Mode ("dark" or "light") used by this site unless overridden in front matter
	Mode string

	// Name (not path) of Theme used by this site unless overridden in front matter.
	Theme string

	// Directory this site uses to copy themes from. If the -d option was
	// set, use the global factory themes directory. Otherwise, use local copy
	themesPath string

	// List of directories in the source project directory that should be
	// excluded, things like ".git" and "node_modules".
	// Note that direcotory names starting with a "." are excluded too.
	// DO NOT ACCESS DIRECTLY:
	// Use excludedDirs() because it applies other information such as publishDir()
	ExcludeDirs []string

	// List of file extensions to exclude. For example. [ ".css" ".go" ".php" ]
	ExcludeExtensions []string

	// Target subdirectory for assets such as CSS and images.
	// It's expected to be a child of the Publish directory.
	AssetDir string

	// Base directory for URL root, which may be diffferent
	// from its actual root. For example, GitHub Pages prefers
	// the blog to start in /docs instead of root, but
	// a URL would omit it.
	BaseDir string

	// Directory for finished site--rendered HTML & asset output
	Publish string

	// All these files are copied into the HTML header.
	// Example: favicon links.
	Headers string

	// ThemesPath is where all the themes are stored.
	// It is computed at startup based on configuration values.
	// Either it was copied to the site directory or you're
	// using the global theme directory
	themePath string

	// Google Analytics tracking ID
	Ganalytics string

	MarkdownOptions MarkdownOptions

	Authors []author
	Company companyConfig

	// Social media URLs
	Social socialConfig

	// THIS ALWAYS GOES AT THE END OF THE FILE/DATA STRUCTURE
	// User data.
	List interface{}
}

// Indicates whether it's directory, a directory containing
// markdown files, or file, or a Markdown file.
// Used for bit flags
type mdOptions uint8

const (
	// Known to be a directory with at least 1 Markdown file
	markdownDir mdOptions = 1 << iota

	// Known to be a filename with a Markdown extension
	markdownFile

	// Directory. Don't know yet if it contains Markdown files.
	normalDir

	// File. Don't know if it's a markdown file.
	normalFile

	// Set if directory has a file named "index.md", forced to lowercase
	hasIndexMd

	// Set if directory has a file named "README.md", case sensitive
	hasReadmeMd
)

type socialConfig struct {
	DeviantArt string
	Facebook   string
	Github     string
	Gitlab     string
	Instagram  string
	LinkedIn   string
	Pinterest  string
	Reddit     string
	Tumblr     string
	Twitter    string
	Weibo      string
	YouTube    string
}

type companyConfig struct {
	// Company name, like "Metabuzz" or "Example Inc."
	Name string
	URL  string

	// Logo file for the header
	HeaderLogo string
}
type author struct {
	FullName string
	URL      string
	Role     string
}

type authors struct {
	Authors []author
}

// readSiteConfig() opens the expected site.toml file, reads, and
// parses it.
// TODO: Remove? Replaced by Viper?
func (App *App) readSiteConfig() error {
	App.Warning("%s", "readSiteConfig() xxx")
	return readTomlFile(App.Site.configFilePath, &App.Site)
}

// MarkdownOptions consists of goldmark
// options used when converting markdown to HTML.
type MarkdownOptions struct {
	// If true, line breaks are signficant
	hardWraps bool

	// Name of color scheme used for code highlighting,
	// for example, "monokai"
	// List of supported languages:
	// https://github.com/alecthomas/chroma/blob/master/README.md
	// I believe the list of themes is here:
	// https://github.com/alecthomas/chroma/tree/master/styles
	HighlightStyle string

	// Create id= attributes for all headers automatically
	headingIDs bool
}

// writeSiteConfig() writes the contents of App.Site
// and creates or replaces a TOML file in the
// project's site subdirectory.
// Assumes you're in the project directory.
func (App *App) writeSiteConfig() error {
	App.Warning("writeSiteConfig(%s) xxx", App.Site.configFilePath)
	return writeTomlFile(App.Site.configFilePath, App.Site)
}

// newSite() attempts to create an empty
// project site using the
// supplied directory name.
func (App *App) newSite(sitename string) error {
	fmt.Println("App.Args.NewSite: " + App.Args.NewSiteName)
	if sitename == "" {
		return errCode("1013", "")
	}
	// Do a simplistic, fallible check to see if there's
	// already a site present and quit if so.
	if isProject(sitename) {
		return errCode("0951", sitename)
	}

	// Create the site subdirectory.
	err := os.MkdirAll(sitename, PROJECT_FILE_PERMISSIONS)
	if err != nil {
		return errCode("401", sitename)
	}

	// Make it the current directory.
	if err := os.Chdir(sitename); err != nil {
		return errCode("1106", sitename)
	}

	// Create minimal directory structure: Publish directory
	// .site directory, .themes, etc.
	if err := createDirStructure(&siteDirs); err != nil {
		return err
	}

	// Create its site.toml file
  App.siteDefaults()
  if err := App.writeSiteConfig(); err != nil {
    // Custom error message already generated
    return errCode("PREVIOUS", err.Error(), App.Site.configFilePath)
  }

	// Copy all themes
	// from the user application data directory to the site.
	App.Warning("App.themesPath: %s", App.themesPath)
	App.Warning("App.Site.themesPath: %s", App.Site.themesPath)
	promptString("Ready to continue? ")
	App.Warning("Pretending to copy theme directory from %s to %s",
		App.themesPath, App.Site.themesPath)
	/*
		err = copyDirAll(App.themesPath,App.Site.themesPath)
		if err != nil {
			// TODO: Make more specific? and why did this cause a runtime error?
			//QuitError(errCode("PREVIOUS", err.Error()))
			QuitError(errCode("0911", "from '"+App.themesPath+"' to '"+App.Site.themesPath+"'"))
		}
	*/
	// Create a little home page
	indexMd = fmt.Sprintf(indexMd, sitename, sitename)
	return writeTextFile("index.md", indexMd)

}


// siteDefaults() computes values for location of site
// theme files, publish directory, etc. 
// Most of them are relative to the site directory.
// It must be called after config files are read.
func (App *App) siteDefaults() {
  App.Warning("%s","siteDefaults() xxx")
	App.Site.path = currDir()
  App.Site.configFilePath = filepath.Join(App.Site.path, siteConfigSubDir, siteConfigFilename)
	App.themesPath = filepath.Join(cfgString("configdir"), THEME_SUBDIRNAME)
	App.Site.Publish = filepath.Join(App.Site.path, PublishSubDirName)
	App.Site.Headers = filepath.Join(App.Site.path, headersDir)
	App.Site.commonDir = filepath.Join(App.Site.path, commonSubDirName)
	if App.Flags.DontCopy {
		App.Site.themesPath = App.themesPath // xxx
	} else {
		App.Site.themesPath = filepath.Join(App.Site.path, siteThemeDir)
	}
	///xxx
	fmt.Printf("\tApp.Site.path: %+v\n", App.Site.path)
	fmt.Printf("\tApp.Site.name: %+v\n", App.Site.Name)
	fmt.Printf("\tApp.Page.filename: %+v\n", App.Page.filename)
	fmt.Printf("\tApp.Site.themesPath: %+v\n", App.Site.themesPath)
	fmt.Printf("\tApp.themesPath: %+v\n", App.themesPath)
}



