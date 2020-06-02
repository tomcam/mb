package main

import (
	//"os"
	//"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"html/template"
	"log"
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
	Prefs       *Prefs
	Site        *Site
	Page        *Page
	FrontMatter *FrontMatter
	// INTERNAL
	infoLog  *log.Logger
	errorLog *log.Logger

	// Currently loaded theme
	Theme Theme

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
	// Tell viper where to look for config file.
	// It can look in as many places as you want.
	viper.AddConfigPath(".")
	viper.AddConfigPath(filepath.Join(homeDir(), GLOBAL_CONFIG_DIRNAME))
	viper.SetConfigName(PRODUCT_NAME)
	// viper likes to apply its own file extensions
	viper.SetConfigType("toml")
	// TODO: Get this right when I've nailed the other Viper stuff
	viper.AutomaticEnv()
	// Read in command line options, and get the
	// location of the configuration directory that
	// itself points to metabuzz.toml
	if err := viper.ReadInConfig(); err != nil {
		// TODO: Give this a standard error code and display it
		//fmt.Println("error reading in config file:", err.Error())
		if err, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// TODO: Handle error properly
			App.QuitError(err)
		} else {
			// Ignore case where there simply wasn't a config file,
			// since it's not a requirement.
			App.Verbose("No configuration file found")
		}
	}
	// Are we going to look in the local directory for
	// site assets, themes, etc., or are we going to
	// use the standard application configuration directory?
	// This determines its location.
	App.Prefs.configDir = cfgString("configdir")
	if App.Prefs.configDir == "" {
		App.Prefs.configDir = configDir()
	}
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
			assets:  []string{},
			Article: []byte{},
		},
		Prefs: &Prefs{
			configDir: ".",
		},
		Site: &Site{
			// Assets just go into the publish directory
			AssetDir:     ".",
			CommonSubDir: commonSubDirName,
			//configFile: filepath.Join(siteConfigSubDir, siteConfigFilename),
			dirs:     map[string]mdOptions{},
			Language: "en",
			MarkdownOptions: MarkdownOptions{
				hardWraps:      false,
				HighlightStyle: "github",
				headingIDs:     true,
			},
		},
		FrontMatter: &FrontMatter{
			// Name of default theme can overridden in a config file
			Theme: DEFAULT_THEME_NAME,
		},
	}
	// Add config/env support from cobra and viper
	App.addCommands()

	App.themesPath = filepath.Join(configDir(), themeSubDirName)
	App.sCodePath = filepath.Join(configDir(), sCodeSubDirName)
	App.funcs = template.FuncMap{
		"ftime":    App.ftime,
		"hostname": App.hostname,
		"inc":      App.inc,
		"path":     App.path,
		"scode":    App.scode,
	}
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
