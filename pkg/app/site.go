package app

import (
	"fmt"
	"github.com/tomcam/mb/pkg/defaults"
	"github.com/tomcam/mb/pkg/errs"
	"os"
	"path/filepath"
)

// Site contains configuration specific to each site, such as
// its title, publish directory, and branding string.
type Site struct {

	// Target subdirectory for assets such as CSS and images.
	// It's expected to be a child of the Publish directory.
	// The function App.assetDir() computes the full
	// path of that directory, based on the app path,
	// the current theme, etc.
	// See also its subdirectories, CSSDir and ImageDir
	AssetDir string

	// Make it easy if you just have 1 author.
	Author author

	// List of authors with roles and websites in site.toml
	Authors []author

	// Base directory for URL root, which may be different
	// from its actual root. For example, GitHub Pages prefers
	// the blog to start in /docs instead of root, but
	// a URL would omit it.
	BaseDir string

	// Site's branding, any string, that user specifies in site.toml.
	// So, for example, if the Name is 'my-project' this might be
	// 'My Insanely Cool Project"
	Branding string

	// Full pathname of common directory. Derived from CommonSubDir
	commonPath string

	// Company name & other info user specifies in site.toml
	Company companyConfig

	// Subdirectory under the AssetDir where CSS files go
	CSSDir string

	// List of all directories in the site
	dirs map[string]dirInfo

	// List of directories in the source project directory that should be
	// excluded, things like ".git" and "node_modules".
	// Note that directory names starting with a "." are excluded too.
	// DO NOT ACCESS DIRECTLY:
	// Use excludedDirs() because it applies other information such as PublishDir()
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

	// Subdirectory under the AssetDir where image files go
	ImageDir string

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
	PublishDir string

	// Full path of shortcode dir for this project. It's computed
	// at runtime using SCodeDir, also in this struct.
	sCodePath string

	// Full path of file containing JSON version of site text
	// to be indexed
	// xxx
	SearchJSONFilePath string

	// Full path to site config file
	siteFilePath string

	// Social media URLs
	Social socialConfig

	// Name (not path) of Theme used by this site unless overridden in front matter.
	Theme string

	// Target subdirectory where themes get copied for publication.
	// It's expected to be a child of the Publish directory.
  ThemesDir string

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
// potentially also other stuff like file create date.
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
type MdOptions uint8

type dirInfo struct {
	mdOptions MdOptions
}

const (
	// Known to be a directory with at least 1 Markdown file
	MarkdownDir MdOptions = 1 << iota

	// Known to be a filename with a Markdown extension
	MarkdownFile

	// Directory. Don't know yet if it contains Markdown files.
	NormalDir

	// File. Don't know if it's a markdown file.
	NormalFile

	// Set if directory has a file named "index.md", forced to lowercase
	HasIndexMd

	// Set if directory has a file named "README.md", case sensitive
	HasReadmeMd
)

// IsOptionSet returns true if the opt bit is set.
func (m MdOptions) IsOptionSet(opt MdOptions) bool {
	return m&opt != 0
}

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
	// If true, line breaks are significant
	hardWraps bool

	// Name of color scheme used for code highlighting,
	// for example, "monokai"
	// List of supported languages:
	// https://github.com/alecthomas/chroma/blob/master/README.md
	// I believe the list of themes is here:
	// https://github.com/alecthomas/chroma/tree/master/styles
	HighlightStyle string

	// Create id= attributes for all headers automatically
	HeadingIDs bool
}

// writeSiteConfig() writes the contents of App.Site
// to .site/site.toml
// and creates or replaces a TOML file in the
// project's site subdirectory.
// Assumes you're in the project directory.
func (a *App) writeSiteConfig() error {
	return writeTomlFile(a.Site.siteFilePath, a.Site)
}

// NewSite() attempts to create an empty
// project site using the supplied directory name.
func (a *App) NewSite(sitename string) error {
	if sitename == "" {
		return errs.ErrCode("1013", "")
	}
  var err error
  inDirAlready := false
  // Lets you turn the current directory into a site.
  // Kind of annoying in the source though.
  // It's just so you can do this:
  //   $ mb new site .
  if sitename == "." {
    inDirAlready = true
    sitename = filepath.Base(currDir())
  }
	// Do a simplistic, fallible check to see if there's
	// already a site present and quit if so.
	// EXCEPTION: You get to assign one site name to
	// testsite= in metabuzz.toml, and that site
	// gets destroyed and replaced with the new one.
	if isProject(sitename) && sitename != cfgString("testsite") {
		return errs.ErrCode("0951", sitename)
	}

  a.Verbose("Creating site named %s", sitename)
	// Create the site subdirectory.
  // Don't do it if already in the directory
  if !inDirAlready {
    err = os.MkdirAll(sitename, defaults.ProjectFilePermissions)
    if err != nil {
      return errs.ErrCode("401", sitename)
    }
  }
	// Make it the current directory
  // unless it was invoked as $ mb new site .
  if !inDirAlready {
    if err = os.Chdir(sitename); err != nil {
      return errs.ErrCode("1106", sitename)
    }
  }
	// Create minimal directory structure: Publish directory
	// .site directory, .themes, etc.
	if err = createDirStructure(&defaults.SiteDirs); err != nil {
		return errs.ErrCode("PREVIOUS", err.Error())
	}
	a.SiteDefaults()
	// Create its site.toml file
	if err = a.writeSiteConfig(); err != nil {
		// Custom error message already generated
		return errs.ErrCode("PREVIOUS", err.Error(), a.Site.siteFilePath)
	}

	// Copy all themes from the application data directory
	// to the site directory.
	err = copyDirAll(a.themesPath, a.Site.themesPath)
	if err != nil {
		a.QuitError(errs.ErrCode("0911", "from '"+a.themesPath+"' to '"+a.Site.themesPath+"'"))
	}

	// Copy all scodes from the user application data directory
	// to the project directory.
	err = copyDirAll(a.sCodePath, a.Site.sCodePath)
	if err != nil {
		a.QuitError(errs.ErrCode("0915", "from '"+a.sCodePath+"' to '"+a.Site.sCodePath+"'"))
	}

	// Copy all header tags.
	err = copyDirAll(a.headTagsPath, a.Site.headTagsPath)
	if err != nil {
		a.QuitError(errs.ErrCode("0923", "from '"+a.sCodePath+"' to '"+a.Site.sCodePath+"'"))
	}

	a.Site.AssetDir = filepath.Join(a.Site.PublishDir, a.Site.AssetDir)

	// Create a little home page
	// The home page is based on the site name.
	// Remove its path, leaving just the directory name.
	sitename = filepath.Base(sitename)
	starterMd := fmt.Sprintf(defaults.IndexMd, sitename, sitename)
	return writeTextFile("index.md", starterMd)

}

// SiteDefaults() computes values for location of site
// theme files, publish directory, etc.
// Most of them are relative to the site directory.
// It must be called after command line flags, env
// variables, and other application configuration has been done.
func (a *App) SiteDefaults() {
	// Initial defaults. Some values may immediately be overridden
	// when a.readSiteConfig() is called.
	a.Site.path = currDir()
	// Next read in the site configuration file, which may override things
	// like AssetDir and Publish.
	a.Site.siteFilePath = filepath.Join(a.Site.path, defaults.GlobalConfigurationDirName,
		defaults.SiteConfigDir, defaults.SiteConfigFilename)
	if err := a.readSiteConfig(); err != nil {
		displayError(errs.ErrCode("PREVIOUS", ""))
	}

	// Unlike the other dot directories, Publish is only
	// 1 level deep. It is not nested inside the .mb directory
	a.Site.PublishDir = filepath.Join(a.Site.path, defaults.PublishDir)

	// xxx
	a.Site.SearchJSONFilePath = filepath.Join(a.Site.PublishDir, defaults.SearchJSONFilename)

	a.commonPath = filepath.Join(a.configDir, defaults.CommonDir)
	a.headTagsPath = filepath.Join(a.configDir, defaults.HeadTagsDir)
	a.sCodePath = filepath.Join(a.configDir, defaults.SCodeDir)
	a.themesPath = filepath.Join(a.configDir, defaults.ThemeDir)

	a.Site.commonPath = filepath.Join(
		a.Site.path, defaults.GlobalConfigurationDirName, defaults.CommonDir)
	a.Site.headTagsPath = filepath.Join(
		a.Site.path, defaults.GlobalConfigurationDirName, defaults.HeadTagsDir)
	a.Site.sCodePath = filepath.Join(
		a.Site.path, defaults.GlobalConfigurationDirName, defaults.SCodeDir)
	a.Site.themesPath = filepath.Join(
		a.Site.path, defaults.GlobalConfigurationDirName, defaults.ThemeDir)

	if a.Flags.DontCopy {
		fmt.Println("TODO: Finish the -d option here in site.go")
		a.Site.themesPath = a.themesPath
		a.Site.sCodePath = a.sCodePath
	}
}

// readSiteConfig() opens the expected site.toml file, reads,
// parses it, and assigns its values to App.Site.
func (a *App) readSiteConfig() error {
	return readTomlFile(a.Site.siteFilePath, &a.Site)
}
