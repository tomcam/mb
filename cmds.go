package main

import (
	"flag"
	"github.com/spf13/cobra"
	//"github.com/spf13/viper"
)

var (

	// Declare command-line subcommand to display config info
	cmdInfo = flag.NewFlagSet("info", flag.ExitOnError)

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

	// From part of Name for new theme from=default
	NewThemeFrom string
	// to part of name for new theme to=mytheme
	NewThemeTo string
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
			Short: "info TODO: Document this",
			Long:  `info: TODO: Long version`,
			Run: func(cmd *cobra.Command, args []string) {
				App.info()
			},
		}

		/*****************************************************
		  TOP LEVEL COMMAND: build
		 *****************************************************/
		cmdBuild = &cobra.Command{
			Use:   "build",
			Short: "build: Generates the site HTML and copies to publish directory",
			Long: `
      build: Generates the site HTML and copies to publish directory 

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
        if len(args) > 0 {
          for i, _ := range args {
            err := App.publishFile(args[i])
            if err != nil {
              QuitError(err)
            }
          }
        } else {
          App.Warning("%s", "Fake build xxx")
        }
			},
		}

		/*****************************************************
		  TOP LEVEL COMMAND: new
		 *****************************************************/
		cmdNew = &cobra.Command{
			Use:   "new",
			Short: "new commands: cli new site|theme",
			Long:  `new commands: Long version`,
			/*
							Run: func(cmd *cobra.Command, args []string) {
								//App.cmdNew(args)
				        // xxx
				      },
			*/
		}

		/*****************************************************
		    Subcommand: new site
		*****************************************************/
		//err      error

		cmdNewSite = &cobra.Command{}

		/*****************************************************
		     Subcommand: new theme
		*****************************************************/
		cmdNewTheme = &cobra.Command{}
	)

	// Example command line:
	// new theme --from=pillar
	cmdNewTheme.Flags().StringVarP(&App.Args.NewThemeFrom, "from", "f", DEFAULT_THEME_NAME, "name of theme to copy from")
	// Example command line:
	// new theme --to=mytheme
	cmdNewTheme.Flags().StringVarP(&App.Args.NewThemeTo, "to", "t", "", "name of theme to create (required)")
	cmdNewTheme.MarkFlagRequired("to")

	// newTheme
	// xxx
	// Example command line:
	// new theme
	cmdNew.AddCommand(cmdNewTheme)

	// Example command line:
	// new site foo --name=mysite
	cmdNewSite.Flags().StringVarP(&App.Args.NewSiteName, "name", "n", "", "name of new site (follow file naming conventions)")

	// Example command line:
	// new site
	cmdNew.AddCommand(cmdNewSite)

	// Example command line:
	// new
	App.Cmd.AddCommand(cmdNew)
	App.Cmd.AddCommand(cmdBuild)
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
