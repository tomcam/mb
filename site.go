package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// Site contains configuration specific to each site, such as
// its title, publish directory, and branding string.
type Site struct {

	// Target subdirectory for assets such as CSS and images.
	// It's expected to be a child of the Publish directory.
	AssetDir string

	// Make it easy if you just have 1 author.
	Author author

	// List of authors with roles and websites in site.toml
	Authors []author

	// Base directory for URL root, which may be diffferent
	// from its actual root. For example, GitHub Pages prefers
	// the blog to start in /docs instead of root, but
	// a URL would omit it.
	BaseDir string

	// Site's branding, any string, that user specifies in site.toml
	//. So, for example,
	// if the Name is 'my-project' this might be
	// 'My Insanely Cool Project"
	Branding string

	// Full pathname of common directory. Derived from CommonSubDir
	commonPath string

	// Company name & other info user specifies in site.toml
	Company companyConfig

	// List of all directories in the site
	dirs map[string]mdOptions

	// List of directories in the source project directory that should be
	// excluded, things like ".git" and "node_modules".
	// Note that direcotory names starting with a "." are excluded too.
	// DO NOT ACCESS DIRECTLY:
	// Use excludedDirs() because it applies other information such as publishDir()
	ExcludeDirs []string

	// List of file extensions to exclude. For example. [ ".css" ".go" ".php" ]
	ExcludeExtensions []string

	// Google Analytics tracking ID specified in site.toml
	Ganalytics string

	// All these files are copied into the HTML header.
	// Example: favicon links.
	HeadTags []string

	// Full path of header tags for "code injection"
	headTagsPath string

	// for HTML header, as in "en" or "fr"
	Language string

	// Flags indicating which non-CommonMark Markdown extensions to use
	MarkdownOptions MarkdownOptions

	// Mode ("dark" or "light") used by this site unless overridden in front matter
	Mode string

	// Site's project name, so it's a filename.
	// It's an identifier so it should be in slug format:
	// Preferably just alphanumerics, underline or hyphen, and
	// no spaces, for example, 'my-project'
	Name string

	// Full path of the directory where the site source is.
	path string

	// Directory for finished site--rendered HTML & asset output
	Publish string

	// Fully qualified directory name of the location themes get copied
	// to in the published site
	pubThemesPath string

	// Full path of shortcode dir for this project. It's computed
	// at runtime using SCodeDir, also in this struct.
	sCodePath string

	// Full path to site config file
	siteFilePath string

	// Social media URLs
	Social socialConfig

	// Language tag for html lang=
	// Name (not path) of Theme used by this site unless overridden in front matter.
	Theme string

	// Directory this site uses to copy themes from. If the -d option was
	// set, use the global factory themes directory. Otherwise, use local copy
	themesPath string

	// All the rendered pages on the site, plus meta information.
	// Index by the fully qualified path name of the source .md file.
	WebPages map[string]WebPage

	// THIS ALWAYS GOES AT THE END OF THE FILE/DATA STRUCTURE
	// User data.
	List interface{}
}

// Everything relevant about the page to be published,
// namely its rendered text and what's in the front matter, but
// potentinally also other stuff like file create date.
type WebPage struct {
	// Rendered text, the HTML after going through templates
	html []byte
}

// OpenGraph tags
// https://ogp.me/
type OG struct {
	Title              string
	Type               string
	Image              string
	Url                string
	Audio              string
	Description        string
	Determiner         string
	Locale             string
	Locale_alternative string
	Site_name          string
	Video              string
	// Structured properties
	Image_url        string
	Image_secure_url string
	Image_type       string
	Image_width      string
	Image_height     string
	Image_alt        string
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
// to .site/site.toml
// and creates or replaces a TOML file in the
// project's site subdirectory.
// Assumes you're in the project directory.
func (App *App) writeSiteConfig() error {
	return writeTomlFile(App.Site.siteFilePath, App.Site)
}

// newSite() attempts to create an empty
// project site using the
// supplied directory name.
func (App *App) newSite(sitename string) error {
	if sitename == "" {
		return errCode("1013", "")
	}
	// Do a simplistic, fallible check to see if there's
	// already a site present and quit if so.
	// EXCEPTION: You get to assign one site name to
	// testsite= in metabuzz.toml, and that site
	// gets destroyed and replaced with the new one.
	if isProject(sitename) && sitename != cfgString("testsite") {
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
		return errCode("PREVIOUS", err.Error())
	}
	App.siteDefaults()
	// Create its site.toml file
	if err := App.writeSiteConfig(); err != nil {
		// Custom error message already generated
		return errCode("PREVIOUS", err.Error(), App.Site.siteFilePath)
	}

	// Copy all themes from the application data directory
	// to the site directory.
	err = copyDirAll(App.themesPath, App.Site.themesPath)
	if err != nil {
		App.QuitError(errCode("0911", "from '"+App.themesPath+"' to '"+App.Site.themesPath+"'"))
	}

	// Copy all scodes from the user application data directory
	// to the project directory.
	err = copyDirAll(App.sCodePath, App.Site.sCodePath)
	if err != nil {
		App.QuitError(errCode("0915", "from '"+App.sCodePath+"' to '"+App.Site.sCodePath+"'"))
	}
	App.Site.AssetDir = filepath.Join(App.Site.Publish, App.Site.AssetDir)

	// Create a little home page
	indexMd = fmt.Sprintf(indexMd, sitename, sitename)
	return writeTextFile("index.md", indexMd)

}

// siteDefaults() computes values for location of site
// theme files, publish directory, etc.
// Most of them are relative to the site directory.
// It must be called after command line flags, env
// variables, and other application configuration has been done.
func (App *App) siteDefaults() {
	App.Site.path = currDir()
	// Next read in the site configuration file, which may override things
	// like AssetDir and Publish.
	App.Site.siteFilePath = filepath.Join(App.Site.path, globalConfigurationDirName, siteConfigDir, siteConfigFilename)
	if err := App.readSiteConfig(); err != nil {
		displayError(errCode("PREVIOUS", ""))
	}

	// Unlike the other dot directories, Publish is only
	// 1 level deep. It is not nested inside the .mb directory
	App.Site.Publish = filepath.Join(App.Site.path, publishDir)
	App.Site.pubThemesPath = filepath.Join(App.Site.Publish, themeDir)

	// TODO: Move to  NewDefaultApp() and change to AppDefaults I think
	// scode path?
	App.commonPath = filepath.Join(App.configDir, commonDir)
	App.headTagsPath = filepath.Join(App.configDir, headTagsDir)
	App.sCodePath = filepath.Join(App.configDir, sCodeDir)
	App.themesPath = filepath.Join(App.configDir, themeDir)

	App.Site.commonPath = filepath.Join(App.Site.path, globalConfigurationDirName, commonDir)
	App.Site.headTagsPath = filepath.Join(App.Site.path, globalConfigurationDirName, headTagsDir)
	App.Site.sCodePath = filepath.Join(App.Site.path, globalConfigurationDirName, sCodeDir)
	App.Site.themesPath = filepath.Join(App.Site.path, globalConfigurationDirName, themeDir)

	if App.Flags.DontCopy {
		fmt.Println("TODO: Finish the -d option here in site.go")
		App.Site.themesPath = App.themesPath
		App.Site.sCodePath = App.sCodePath
	}
}

// readSiteConfig() opens the expected site.toml file, reads,
// parses it, and assigns its values to App.Site.
func (App *App) readSiteConfig() error {
	return readTomlFile(App.Site.siteFilePath, &App.Site)
}
