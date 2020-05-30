package main

import (
	//"os"
	"fmt"
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

	// Page being rendered
	html []byte

	// Custom functions used in the template language.
	// All built-in functions must appear here to be publicly available
	funcs map[string]interface{}

	// Identical to funcs except fewerFuncs cannot have shortcode in it
	fewerFuncs map[string]interface{}
}

// initConfig() determines where configuration file (and other
// forms of configuration info) can be found, then reads in
// all that info.
func (App *App) initConfig() {
	// Tell viper where to look for global config file.
	viper.AddConfigPath(App.Prefs.configDir)
	// It can look in as many places as you want.
	viper.AddConfigPath(filepath.Join(homeDir(), GLOBAL_CONFIG_DIRNAME))
	viper.AddConfigPath(".")
	viper.SetConfigName(PRODUCT_NAME)
	// viper likes to apply its own file extensions
	viper.SetConfigType("toml")
	// TODO: Get this right when I've nailed the other Viper stuff
	viper.AutomaticEnv()
	// Read in metabuzz.toml
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("error reading in config file:", err.Error())
		if err, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// TODO: Handle error properly
			QuitError(err)
		} else {
			// Ignore case where there simply wasn't a config file,
			// since it's not a requirement.
			App.Verbose("No configuration file found")
		}
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
			assets: []string{},
		},
		Prefs: &Prefs{
			configDir: ".",
		},
		Site: &Site{
			configFile: filepath.Join(siteConfigSubDir, siteConfigFilename),
			dirs:       map[string]mdOptions{},
			// Assets just go into the publish directory
			AssetDir: ".",
			//SiteConfigFile: filepath.Join(siteConfigSubDir, SITE_CONFIG_FILENAME),
			CommonSubDir: commonSubDirName,
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

	App.funcs = template.FuncMap{/* "scode": App.scode, */
		"ftime": App.ftime, 
    /*"hostname": App.hostname, "path": App.path, "inc": App.inc */
  }



	// CONFIG HAS NOT BEEN READ   YET
	return &App
}


