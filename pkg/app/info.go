package app

import (
	//"github.com/spf13/viper"
	"fmt"
	"github.com/tomcam/mb/pkg/defaults"
	"os"
)

// info() displays relevant site configuration info for
// debug purposes. If -v (verbose mode), also print data structures
func (a *App) info() {
	a.SiteDefaults()
	//fmt.Println("*** foo.bar: " + viper.GetString("foo.bar"))
	//fmt.Println("*** configFilePath: " + viper.GetString("configFilePath"))
	//fmt.Println("*** a.configFilePath: " + a.Prefs.configFilePath)
	fmt.Println("Home dir: " + homeDir())
	fmt.Println("Reported current dir: " + a.Site.path)
	fmt.Println("Actual current dir: " + currDir())
	exists("scode path", a.Site.sCodePath)
	fmt.Println("a.Flags.Verbose", a.Flags.Verbose)
	exists("Default config directory", configDir())
	exists("Actual config directory", a.configDir)
	exists("Site file: ", a.Site.siteFilePath)
	exists("Theme directory", a.themesPath)
	fmt.Println("Code highlighting style: ", a.Site.MarkdownOptions.HighlightStyle)
	fmt.Println("Default theme: ", a.defaultTheme())
	fmt.Println("Highlight:", cfgString("highlight"))
	if isProject(".") {
		fmt.Println("This appears to be a project/site source directory")
		exists("Site directory: ", a.Site.path)
		exists("Publish directory", a.Site.Publish)
		exists("Theme directory", a.Site.themesPath)
		exists("Headers directory", a.Site.headTagsPath)
		//:exists("Asset directory", a.assetDir())
		exists("Shortcode directory: ", a.Site.sCodePath)
	}
	if a.Flags.Verbose {
		fmt.Fprintf(os.Stdout, "\nPrefs\n-----\n%#v\n", a)
		fmt.Fprintf(os.Stdout, "\nFront matter\n----- ------\n%#v\n", a.FrontMatter)
		fmt.Fprintf(os.Stdout, "\nSite\n----\n%#v\n", a.Site)
		fmt.Fprintf(os.Stdout, "\nDirectory structure for site\n----\n%#v\n",
			defaults.SiteDirs)
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
