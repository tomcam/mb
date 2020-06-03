package main

import (
	//"github.com/spf13/viper"
	"fmt"
	"os"
)

// info() displays relevant site configuration info for
// debug purposes. If -v (verbose mode), also print data structures
func (App *App) info() {
  App.siteDefaults()
	//fmt.Println("*** foo.bar: " + viper.GetString("foo.bar"))
	//fmt.Println("*** configDir: " + viper.GetString("configDir"))
	//fmt.Println("*** App.Prefs.configDir: " + App.Prefs.configDir)
	fmt.Println("Home dir: " + homeDir())
	fmt.Println("Reported current dir: " + App.Site.path)
	fmt.Println("Actual current dir: " + currDir())
	exists("scode path", App.Site.sCodePath)
	fmt.Println("App.Flags.Verbose", App.Flags.Verbose)
	exists("Default config directory", configDir())
	exists("Actual config directory", App.Prefs.configDir)
  exists("Config file: ",App.Site.configFilePath)
	exists("Theme directory", App.themesPath)
	fmt.Println("Code highlighting style: ", App.Site.MarkdownOptions.HighlightStyle)
	fmt.Println("Default theme: ", App.defaultTheme())
	fmt.Println("Highlight:", cfgString("highlight"))
	if isProject(".") {
		fmt.Println("This appears to be a project/site source directory")
		exists("Site directory: ", App.Site.path)
		exists("Publish directory", App.Site.Publish)
		exists("Theme directory", App.Site.themesPath)
		exists("Headers directory", App.Site.Headers)
		//:exists("Asset directory", App.assetDir())
		exists("Shortcode directory: ", App.Site.sCodePath)
	}
	if App.Flags.Verbose {
		fmt.Fprintf(os.Stdout, "\nPrefs\n-----\n%#v\n", App.Prefs)
		fmt.Fprintf(os.Stdout, "\nFront matter\n----- ------\n%#v\n", App.FrontMatter)
		fmt.Fprintf(os.Stdout, "\nSite\n----\n%#v\n", App.Site)
		fmt.Fprintf(os.Stdout, "\nDirectory structure for site\n----\n%#v\n", siteDirs)
	}
}

// exists() is a helper utility that simply displays a filename and
// shows if it's actually present
func exists(description, filename string) {
	found := false
	if isDirectory(filename) {
		found = true
	}
	fmt.Print(description, " ", filename)
	if fileExists(filename) {
		found = true
	}

	if found {
		fmt.Println(": (present)")
	} else {
		fmt.Println(": (Not present)")
	}
}
