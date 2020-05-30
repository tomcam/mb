package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	//"html/template"
	"log"
	"os"
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
	// viper likes to apply its on file extensions
	viper.SetConfigType("toml")
	// TODO: Get this right when I've nailed the other Viper stuff
	viper.AutomaticEnv()
	// Read in metabuzz.toml
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("error reading in config file:", err.Error())
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// TODO: Handle error properly
			fmt.Println("config file read error:", err.Error())
			os.Exit(1)
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

	// CONFIG HAS NOT BEEN READ   YET
	return &App
}

// App.Verbose() displays a message followed
// by a newline to stdout
// if the verbose flag was used. Formats it like Fprintf.
func (App *App) Verbose(format string, a ...interface{}) {
	if App.Flags.Verbose {
		fmt.Println(App.fmtMsg(format, a...))
	}
}

// App.Warn() displays a message followed by a newline
// to stdout, preceded by the text "Warning: "
// Overrides the verbose flag. Formats it like Fprintf.
func (App *App) Warning(format string, a ...interface{}) {
	fmt.Println("Warning: " + App.fmtMsg(format, a...))
}

// fmtMsg() formats string like Fprintf and writes to a string
func (App *App) fmtMsg(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

// QuitError() displays the error passed to it and exits
// to the operating system, returning a 1 (any nonzero
// return means an error ocurred).
// Normally functions that can generate a runtime error
// do so by returning an error. But sometimes there's a
// constraint, for example, fulfilling an interface method
// that doesn't support this practice.
func QuitError(e error) {
	if e == nil {
		os.Exit(0)
	} else {
		fmt.Println(e)
		os.Exit(1)
	}
}
