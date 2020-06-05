package main

var (
	Version = productName + " version " +
		"0.45.0"

	// Tiny starter file for index.md
	indexMd = `
# %s

Welcome to %s
`
	// Directory configuration for a project--a new site.
	siteDirs = [][]string{
		{publishDir},
		{globalConfigurationDirName, commonDir},
		{globalConfigurationDirName, headersDir},
		{globalConfigurationDirName, sCodeDir},
		{globalConfigurationDirName, siteConfigDir},
		{globalConfigurationDirName, themeDir},
	}
	// Markdown file extensions
	markdownExtensions = searchInfo{list: []string{
		".Rmd",
		".md",
		".mdown",
		".mdtext",
		".mdtxt",
		".markdown",
		".mdwn",
		".mkd",
		".mkdn",
		".text"}, sorted: false}

	// excludedAssetExtensions are the extensions of files in the asset
	// directory that should NOT be copied to the publish directory.
	// The contents of a theme directory mix both things you want copied,
	// like CSS files, and things you don't, like TOML or Markdown files.
	excludedAssetExtensions = searchInfo{list: []string{
		".html",
		".toml",
		".bak",
	}, sorted: false}
)

const (
	// Name of the subdirectory that holds shared text.
	// Excluded from publishing.
	commonDir = "common"

	// Name of the subdirectory the rendered files get rendered
	// to. It can't be changed because it's used to determine
	// whether a site is contained within its parent directory.
	// Excluded from publishing.
	publishDir = ".pub"

	// Name of the subdirectory containing files that get copied
	// into the header of each HTML file rendered by the site
	// Excluded from publishing.
	headersDir = "headers"

	// Name of subdirectory containing shortcode files
	// Excluded from publishing.
	sCodeDir = "scodes"

	// Name of subdirectory within the theme that holds help & sample files
	// for that theme.
	THEME_HELP_SUBDIRNAME = ".help"

	// Name of subdirectory under the publish directory for style sheet assets
	// (Currently not well thought out nor in use, though assets directory is
	// being used)
	DEFAULT_PUBLISH_CSS_SUBDIRNAME = "css"

	// Name of subdirectory under the publish directory for image assets
	// (Currently not well thought out nor in use, though assets directory is
	// being used)
	DEFAULT_PUBLISH_IMG_SUBDIRNAME = "img"

	// Name of theme used to create other themes, or theme to be
	// used if for some reason no others are present
	defaultThemeName = "wide"

	// Name of the directory that holds items used by projects, such
	// as themes and shortcodes.
	// TODO: Change this when I settle on a product name
	globalConfigurationDirName = ".mb"

	// Default file extension used by configuration files.
	configFileDefaultExt = "toml"

	// A configuration file passed to the command line.
	CONFIG_FILENAME = productName + "." + configFileDefaultExt

	// The configuration file in the user's application
	// data directory, without the path.
	APP_DATA_CONFIG_FILENAME = productName + "." + configFileDefaultExt

	// The local configuration file name without the path.
	LOCAL_CONFIG_FILENAME = productName + "." + configFileDefaultExt

	// By default, the published site gets its theme from a local copy
	// within the site directory. It then copies from that copy to
	// generate pages in the Publish directory. Helps prevent unintended changes
	// from being made to the originals, and makes it much easier to
	// make theme changes, especially if you're a noob or just want to
	// type less.
	themeDir = "themes"

	// Configuration file found in the current site source directory
	SOURCE_DIR_CONFIG_FILENAME = productName + "." + configFileDefaultExt

	// Actual colloquial name for this product
	// but used in directories & other
	// purposes, like storing config files.
	// Make it in lowercase. One word,
	// like docset or metabuzz.
	// TODO: If this changes update GLOBAL_CONFIG_DIRNAME
	// TODO: Change this when I settle on a product name, and also change PRODUCT_SHORT_NAME
	productName = "metabuzz"

	// Abbreviation, used for name command line program.
	ProductShortName = "mb"

	// Values set through the environment as opposed to config files
	// or command line options.
	// A short version of the product name
	// used as a prefix for environment variables.
	// TODO: Change this when I settle on a product name
	PRODUCT_ENV_PREFIX = "MBZ_"
	// Examples:
	//PRODUCT_ENV_PREFIX+"DEFAULT_THEME"
	//PRODUCT_ENV_PREFIX+"SC_DIR"

	// The permissions given to output files, and also to
	// configuration files.
	// 0755 means the owner can read, write and execute (first 7)
	// Other people can only read (5 and 5). That makes sense
	// for a webserver
	PUBLIC_FILE_PERMISSIONS = 0755

	// Objects that must be visible to the project, but not the public
	PROJECT_FILE_PERMISSIONS = 0700

	// Name of the subdirectory in the project where the site info is held.
	// That incldes the site.toml file and also the publish directory.
	siteConfigDir = "site"

	// Name of the file that holds site configuration information
	siteConfigFilename = "site" + "." + configFileDefaultExt

	// String that precedes error codes
	errorCodePrefix = "mbz"
)

// optionSet() returns whether a particular bit is set in bitwise command parsing.
// That's a bad description but I'm tired.
func optionSet(b, options mdOptions) bool {
	return b&options != 0
}
