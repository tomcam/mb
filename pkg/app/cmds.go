package app

import (
	//"os"
	"flag"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomcam/mb/pkg/defaults"
	//"github.com/tomcam/mb/pkg/errs"
	//"github.com/spf13/viper"
)

var (

	// Temp command to extract config files
	cmdCfgGen = flag.NewFlagSet("cfggen", flag.ExitOnError)

	// Declare command-line subcommand to display config info
	cmdInfo = flag.NewFlagSet("info", flag.ExitOnError)

	// Declare command-line subcommand to build  a test site
	cmdKitchenSink = flag.NewFlagSet("kitchensink", flag.ExitOnError)

	///cmdK = flag.NewFlagSet("k", flag.ExitOnError)

	// Declare command-line subcommand to build a project
	cmdBuild = flag.NewFlagSet("build", flag.ExitOnError)

	// Declare command-line subcomand for copying theme
	cmdCopyTheme = flag.NewFlagSet("copytheme", flag.ExitOnError)
	cmdCopyFrom  = cmdCopyTheme.String("from", "", "theme to copy")
	cmdCopyTo    = cmdCopyTheme.String("to", "", "name of new theme")

	// Creates a new pageType for an existing theme
	// xxx
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
func (a *App) addCommands() {
	var (
		/*****************************************************
		  TOP LEVEL COMMAND:cfggen
		 *****************************************************/
		cmdCfgGen = &cobra.Command{
			Use:   "cfggen",
			Short: "TEMPORARY function to extract dirs",
			Long: `cfggen: TODO: Long version
`,
			Run: func(cmd *cobra.Command, args []string) {
				/*
				   fs := FS(true)
				   _, err := fs.Open(".mb")
				   if (err !=nil) {
				     fmt.Println("ERROR: " + err.Error())
				     os.Exit(1)
				   }
				*/

			},
		}

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
				a.info()
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
				err := a.build()
				if err != nil {
					a.QuitError(err)
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
					err = a.kitchenSink(args[0])
				} else {
					err = a.kitchenSink("")
				}
				if err != nil {
					a.QuitError(err)
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
		    Subcommand: new pagetype
		*****************************************************/

		cmdNewPageType = &cobra.Command{
			Use:   "pagetype TODO: {sitename}",
			Short: "pagetype",
			Long: `pagetype

      Where {sitename} is a valid directory name. For example, if your site is called basiclaptop.com, you would do this:

      mb new site basiclaptop
`,
			Run: func(cmd *cobra.Command, args []string) {
				// If there are arguments after build, then
				// just convert these files one at at time.
				var newPageType, fromTheme string
				if len(args) > 0 {
					newPageType = args[0]
				} else {
					// Them more likely case: it's build all by
					// itself, so go through the whole directory
					// tree and build as a complete site.
					newPageType = promptString("Name of pagetype to create?")
				}
				fromTheme = promptString("Add this pagetype to which theme?")
				//kjjerr := a.newPageType(fromTheme, newPageType)
				a.newPageType(fromTheme, newPageType)
			},
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
					a.Site.Name = args[0]
				} else {
					// Them more likely case: it's build all by
					// itself, so go through the whole directory
					// tree and build as a complete site.
					a.Site.Name = promptString("Name of site to create?")
				}
				err := a.NewSite(a.Site.Name)
				if err != nil {
					a.QuitError(err)
				} else {
					fmt.Println("Created site ", a.Site.Name)
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
		SourceTheme = a.defaultTheme()
		cmdNewTheme = &cobra.Command{
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
				// Create a new theme from the default theme
				SourceTheme = promptStringDefault("Name to copy it from?", SourceTheme)
				a.newTheme(SourceTheme, NewThemeName)
				/*
					if err := a.newTheme(NewThemeFrom, NewThemeName); err != nil {
						a.QuitError(errs.ErrCode("PREVIOUS", err.Error()))
					}
				*/
				// Could put a message that it was created
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
				promptString("xxx Pretending to create new theme")
				/*
					err := a.NewSite(a.Site.Name)
					if err != nil {
						a.QuitError(err)
					} else {
						fmt.Println("Created site ", a.Site.Name)
					}
				*/
			},
		}
	)

	// Example command line:
	// new site
	cmdNew.AddCommand(cmdNewSite)
	cmdNew.AddCommand(cmdNewPageType)

	// Example command line:
	// new
	a.Cmd.AddCommand(cmdNew)
	cmdNew.AddCommand(cmdNewTheme)
	cmdNewTheme.AddCommand(cmdNewThemeFrom)

	a.Cmd.AddCommand(cmdBuild)
	a.Cmd.AddCommand(cmdKitchenSink)
	a.Cmd.AddCommand(cmdInfo)
	a.Cmd.AddCommand(cmdCfgGen)

	// Handle global flags such as Verbose
	a.Cmd.PersistentFlags().BoolVarP(&a.Flags.Verbose, "verbose", "v", false, "verbose output")
	a.Cmd.PersistentFlags().BoolVarP(&a.Flags.DontCopy, "dontcopy", "d", false, "don't copy theme file; use global theme")
	// Code highlighting style to use
	a.Cmd.PersistentFlags().StringVarP(&a.Site.MarkdownOptions.HighlightStyle, "highlight", "l", "github", "default code highlighting scheme")

	a.Cmd.PersistentFlags().StringVarP(&a.Args.config, "config", "c",
		defaults.AppDataConfigFilename, "configuration filename")

	// When cobra is ready to go call initConfig()
	cobra.OnInitialize(a.initConfig)
}
