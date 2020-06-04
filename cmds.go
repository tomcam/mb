package main

import (
	"flag"
	"fmt"
	"github.com/spf13/cobra"
	//"github.com/spf13/viper"
)

var (

	// Declare command-line subcommand to display config info
	cmdInfo = flag.NewFlagSet("info", flag.ExitOnError)

	// Declare command-line subcommand to build  a test site
	cmdKitchenSink = flag.NewFlagSet("kitchensink", flag.ExitOnError)

	///cmdK = flag.NewFlagSet("k", flag.ExitOnError)

	// Declare command-line subcommand to build a project
	cmdBuild = flag.NewFlagSet("build", flag.ExitOnError)

	// Declare command-line subcomand for copying theme
	// Example: copytheme -from=default to=newtest
	cmdCopyTheme = flag.NewFlagSet("copytheme", flag.ExitOnError)
	cmdCopyFrom  = cmdCopyTheme.String("from", "", "theme to copy")
	cmdCopyTo    = cmdCopyTheme.String("to", "", "name of new theme")

	// Creates a new pageType for an existing theme
	cmdPageType     = flag.NewFlagSet("pagetype", flag.ExitOnError)
	cmdPageTypeFrom = cmdPageType.String("from", "", "theme to start with")
	cmdPageTypeTo   = cmdPageType.String("to", "", "name of new pagetype for that theme")
)

// Command-line argument values
type Args struct {
	// Name of config file
	config string

	// Name for new site created with new site name=foo
	NewSiteName string
}

// Globally required flags, such as Verbose
type Flags struct {
	// DontCopy means don't copy theme directory to the site directory.
	// Use the global theme set (which means if you change it, it
	// will affect all new sites created using that theme)
	DontCopy bool

	// Global verbose mode
	Verbose bool
}

// addCommands() initializes the set of CLI
// commands, flags, and command-line options,
// then calls initConfig which obotains those options
// from command line, environment, etc.
func (App *App) addCommands() {
	var (

		/*****************************************************
		  TOP LEVEL COMMAND: info
		 *****************************************************/
		cmdInfo = &cobra.Command{
			Use:   "info",
			Short: "Display configuration and debug information about the site",
			Long: `info: TODO: Long version
Show such information as where theme files can be found,
whether the current directory is Metabuzz project, and so on`,
			Run: func(cmd *cobra.Command, args []string) {
				App.info()
			},
		}

		/*****************************************************
		  TOP LEVEL COMMAND:build
		 *****************************************************/
		cmdBuild = &cobra.Command{
			Use:   "build",
			Short: "build: Generates the site HTML and copies to publish directory",
			Long: `"build: Generates the site HTML and copies to publish directory 
      Typical usage:
      : Create the project named mysite in its own directory.
      : (Generates a tiny file named index.md)
      mb new site mysite
      : Make that the current directory. 
      cd mysite
      : Optional step: Write your Markdown here!
      : Find all .md files and convert to HTML
      : Copy them into the publish directory named .pub
      mb build
      : Load the site's home page into a browser.
      : Windows users, omit the open
      open .pub/index.html
`,
			Run: func(cmd *cobra.Command, args []string) {
				err := App.build()
				if err != nil {
					App.QuitError(err)
				}
			},
		}

		/*****************************************************
		  TOP LEVEL COMMAND: kitchensink
		 *****************************************************/
		cmdKitchenSink = &cobra.Command{
			Use:   "kitchensink [sitename]",
			Short: "Generates a test site showing most features",
			Long: `kitchensink:  Builds a disposable site that exercises many Metabuzz features.

      Typical usage:

      : Create a standard test project named mysite in its own directory.
      mb kitchensink mysite

      : Make that the current directory. 
      cd mysite

      : Optional step: Write your Markdown here!

      : Find all .md files and convert to HTML
      : Copy them into the publish directory named .pub
      mb build

      : Load the site's home page into a browser.
      : Windows users, omit the open
      open .pub/index.html
`,
			Run: func(cmd *cobra.Command, args []string) {
				var err error
				if len(args) > 0 {
					err = App.kitchenSink(args[0])
				} else {
					err = App.kitchenSink("")
				}
				if err != nil {
					App.QuitError(err)
				}
			},
		}

		/*****************************************************
		  TOP LEVEL COMMAND: new
		 *****************************************************/
		cmdNew = &cobra.Command{
			Use:   "new",
			Short: "new commands: new site|theme",
			Long: `site: Use new site to start a new project. Use new theme to 
create theme based on an existing one. 

      Typical usage of new site:

      : Create the project named mysite in its own directory.
      : (Generates a tiny file named index.md)
      mb new site mysite

      : Make that the current directory. 
      cd mysite

      : Optional step: Write your Markdown here!

      : Find all .md files and convert to HTML
      : Copy them into the publish directory named .pub
      mb build

      : Load the site's home page into a browser.
      : Windows users, omit the open
      open .pub/index.html
`,
		}

		/*****************************************************
		    Subcommand: new site
		*****************************************************/

		cmdNewSite = &cobra.Command{
			Use:   "site {sitename}",
			Short: "new site mycoolsite",
			Long: `new site {sitename}

      Where {sitename} is a valid directory name. For example, if your site is called basiclaptop.com, you would do this:

      mb new site basiclaptop
`,
			Run: func(cmd *cobra.Command, args []string) {
				// If there are arguments after build, then
				// just convert these files one at at time.
				if len(args) > 0 {
					App.Site.Name = args[0]
				} else {
					// Them more likely case: it's build all by
					// itself, so go through the whole directory
					// tree and build as a complete site.
					App.Site.Name = promptString("Name of site to create?")
				}
				err := App.newSite(App.Site.Name)
				if err != nil {
					App.QuitError(err)
				} else {
					fmt.Println("Created site ", App.Site.Name)
				}
			},
		}

		/*****************************************************
		     Subcommand: new theme
		*****************************************************/

		// the foo part of:
		// new theme foo
		NewThemeName string
		// the bar part of
		// new theme foo from bar
		NewThemeFrom = App.defaultTheme()
		cmdNewTheme  = &cobra.Command{
			Use:   "theme {newtheme} | from {oldtheme} ",
			Short: "new theme mytheme",
			Long: `site: Use new site to start a new project. Use new theme to 
create theme based on an existing one. 

      Typical usage of new theme:

      mb new theme brochure

      : Edit theme files
      vim .mb/.themes/brochure/brochure.css
      vim .mb/.themes/brochure/theme-dark.css
`,
			Run: func(cmd *cobra.Command, args []string) {
				// If there are arguments after build, then
				// just convert these files one at at time.
				if len(args) > 0 {
					NewThemeName = args[0]
				} else {
					// Them more likely case: it's build all by
					// itself, so go through the whole directory
					// tree and build as a complete site.
					NewThemeName = promptString("Name of theme to create?")
				}
				// xxx
				// Create a new theme from the default theme
				App.newTheme(NewThemeName, NewThemeFrom)
				/*
					if err != nil {
						App.QuitError(err)
					} else {
						fmt.Println("Created theme", NewThemeName)
					}
				*/
			},
		}

		cmdNewThemeFrom = &cobra.Command{
			Use:   "theme {newtheme} from {oldtheme} ",
			Short: "new theme mytheme from empty",
			Long: `Create a new theme by copying an existing one. 

      Typical usage of new theme:

      mb new theme brochure from marlow

      : Edit theme files
      vim .mb/.themes/brochure/brochure.css
      vim .mb/.themes/brochure/theme-dark.css
`,
			Run: func(cmd *cobra.Command, args []string) {
				// If there are arguments after build, then
				// just convert these files one at at time.
				if len(args) > 0 {
					promptString("Create theme " + args[0] + "?")
				} else {
					// Them more likely case: it's build all by
					// itself, so go through the whole directory
					// tree and build as a complete site.
					promptString("xxx Name of theme to create?")
				}
				// xxx
				promptString("xxx Prtending to create new theme")
				/*
					err := App.newSite(App.Site.Name)
					if err != nil {
						App.QuitError(err)
					} else {
						fmt.Println("Created site ", App.Site.Name)
					}
				*/
			},
		}
	)

	// Example command line:
	// new theme --from=pillar
	//cmdNewTheme.Flags().StringVarP(&App.Args.NewThemeFrom, "from", "f", DEFAULT_THEME_NAME, "name of theme to copy from")
	// Example command line:
	// new theme --to=mytheme
	//cmdNewTheme.Flags().StringVarP(&App.Args.NewThemeTo, "to", "t", "", "name of theme to create (required)")
	//cmdNewTheme.MarkFlagRequired("to")

	// Example command line:
	// new site
	cmdNew.AddCommand(cmdNewSite)

	// Example command line:
	// new
	App.Cmd.AddCommand(cmdNew)
	cmdNew.AddCommand(cmdNewTheme)
	cmdNewTheme.AddCommand(cmdNewThemeFrom)

	App.Cmd.AddCommand(cmdBuild)
	App.Cmd.AddCommand(cmdKitchenSink)
	App.Cmd.AddCommand(cmdInfo)
	// Handle global flags such as Verbose
	App.Cmd.PersistentFlags().BoolVarP(&App.Flags.Verbose, "verbose", "v", false, "verbose output")
	App.Cmd.PersistentFlags().BoolVarP(&App.Flags.DontCopy, "dontcopy", "d", false, "don't copy theme file; use global theme")
	// Code highlighting style to use
	App.Cmd.PersistentFlags().StringVarP(&App.Site.MarkdownOptions.HighlightStyle, "highlight", "l", "github", "default code highlighting scheme")

	App.Cmd.PersistentFlags().StringVarP(&App.Args.config, "config", "c", APP_DATA_CONFIG_FILENAME, "configuration filename")

	// When cobra is ready to go call initConfig()
	cobra.OnInitialize(App.initConfig)
}
