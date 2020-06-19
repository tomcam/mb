package main

import (
	//"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path/filepath"
)

// App contains all runtime options required to convert a markdown
// file or project to an HTML file or site.
// Compound data structure for config example at
// https://gist.github.com/alexedwards/5cd712192b4831058b21
type App struct {
	Flags Flags
	Args  Args
	// Number of markdown files processed
	fileCount   uint
	Cmd         *cobra.Command
	Site        *Site
	Page        *Page
	FrontMatter *FrontMatter

	// Fully qualified directory name of the common files subdirectory
	commonPath string

	// Fully qualified directory name of application data directory
	configDir string

	// Fully qualified directory name of the header tags directory for "code injection"
	headTagsPath string

	// Location of global themes directory
	themesPath string

	// Location of directory containing shortcode files
	sCodePath string

	// Custom functions used in the template language.
	// All built-in functions must appear here to be publicly available
	funcs map[string]interface{}
	// Copy of funcs but without "scode"
	fewerFuncs map[string]interface{}
}

// initConfig() determines where configuration file (and other
// forms of configuration info) can be found, then reads in
// all that info.
func (App *App) initConfig() {
	// There may or may not be a metabuzz.toml file around redirecting where
	// to look for Metabuzz application data such as themes and shortcodes.
	// So assume it's where the system likes it, under a "metabuzz/.mb" subdirectory.
	App.configDir = configDir()
	// Places to look for a metabuzz.toml pointing to the global application config dir.
	// It can look in as many places as you want.
	// Look in the local directory for a directory named just named ".mb".
	//viper.AddConfigPath(filepath.Join(".", globalConfigurationDirName))
	viper.AddConfigPath(filepath.Join("."))
	// Location to look for metabuzz.toml
	// Look in the ~/ directory for an ".mb" directory.
	viper.AddConfigPath(filepath.Join(homeDir(), globalConfigurationDirName))
	// Name of the config file is metabuzz, dot..
	viper.SetConfigName(productName)
	// toml. viper likes to apply its own file extensions
	viper.SetConfigType(configFileDefaultExt)
	// TODO: Get this right when I've nailed the other Viper stuff
	viper.AutomaticEnv()
	// Read in command line options, and get the
	// location of the configuration directory that
	// itself points to metabuzz.toml
	if err := viper.ReadInConfig(); err != nil {
		// Actually not an error if there's no config file
		// so you have to be in Verbose mode
		App.Verbose(errCode("0126", err.Error()).Error())
	}
	// Are we going to look in the local directory for
	// site assets, themes, etc., or are we going to
	// use the standard application configuration directory?
	// This determines its location.
	if cfgString("configdir") != "" {
		App.configDir = cfgString("configdir")
	}
	App.siteDefaults()
}

// newDefaultApp() allocates an App runtime environment
// No other config info has been read in.
func newDefaultApp() *App {
	App := App{
		Cmd: &cobra.Command{
			// TODO: Don't hardcode this name
			Use:   ProductShortName,
			Short: "Create static sites",
			Long:  `Headless CMS to create static sites`,
		},

		Page: &Page{
			assets:        []string{},
			Article:       []byte{},
			html:          []byte{},
			markdownStart: []byte{},
		},
		Site: &Site{
			// Assets just go into the publish directory
			AssetDir: ".",
			//configFile: filepath.Join(siteConfigDir, siteConfigFilename),
			//dirs:     map[string]mdOptions{},
			dirs:     map[string]dirInfo{},
			WebPages: map[string]WebPage{},
			Language: "en",
			MarkdownOptions: MarkdownOptions{
				hardWraps:      false,
				HighlightStyle: "github",
				headingIDs:     true,
			},
		},
		FrontMatter: &FrontMatter{
			// Name of default theme can overridden in a config file
			Theme: defaultThemeName,
		},
	}
	// Add config/env support from cobra and viper
	App.addCommands()

	App.addTemplateFunctions()

	/*"hostname": App.hostname, "path": App.path, "inc": App.inc */

	// Get a copy of funcs but without
	// scode, because including it would cause a
	// cycle condition for the scode function
	App.fewerFuncs = make(map[string]interface{})
	for key, value := range App.funcs {
		if key != "scode" {
			App.fewerFuncs[key] = value
		}
	}

	// CONFIG HAS NOT BEEN READ   YET
	return &App
}
