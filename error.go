package main

import (
	"fmt"
	"os"
)

/*
	SECTIONS (this may not hold up)

	0100	- Error reading file
	0200	- Error creating file
	0300	- Error deleting file
	0400 	- Error creating directory
	0500	- Error determining directory name
	0600	- Error deleting directory
	0700	- Error reading directory
	0800	- Can't determine the name of something
	0900	- Problem generating something
	1000	- Something's missing that should be there
	1100	- Problem changing to a directory
  1200  - Syntax error!
*/

/* The reason many of these error codes have identical text is that
   the same error occurs but in different places. Since the
   Go lib returns identical error messages for each one, tracking
   down the error code shows us where the error occurred even if the
   executable is stripped of debug info

	Sample usage: return errCode("0401", err.Error())
	Sample usage: return errCode("0401", err.Error(), filename)
	err = copyDirAll(App.themesPath, App.Site.themesPath)
	if err != nil {
		QuitError(errCode("0911", "from '"+App.themesPath+"' to '"+App.Site.themesPath+"'"))
	}
  if err := copyDirOnly(from, to); err != nil {
	  msg := fmt.Sprintf("Unable to copy from pageType directory %s to new pageType directory %s", from, to)
	  return errCode("0906", msg)

		Sample usage:	 msg := fmt.Errorf("Error attempting to create project file %s: %v", projectFile, err.Error()).Error()
*/

var errMsgs = map[string]string{

	// Just print the last error
	"PREVIOUS": " ",

	// 0100	- Error reading file
	"0101": "Error reading front matter",                                            // filename
	"0102": "Unable to open file",                                                   // filename
	"0103": "Error reading front matter",                                            // filename
	"0104": "TOML error reading theme file",                                         // custom message + err.Error()
	"0105": "TOML error reading PageType file",                                      // custom message + err.Error()
	"0106": "Error copying file to publish",                                         // custom message
	"0107": "Error opening a file to publish",                                       // custom message + err.Error()
	"0108": "Error reading theme file",                                              // custom message + err.Error()
	"0109": "Unable to find theme TOML file",                                        // custom message
	"0110": "Error copying CSS file for new theme",                                  // custom message
	"0111": "Error copying file for new theme",                                      // custom message
	"0112": "File doesn't seem to exist",                                            // Filename
	"0113": "File isn't normal",                                                     // filename
	"0114": "Error opening file",                                                    // filename
	"0115": "Unable to find theme TOML file",                                        //  message
	"0116": "TOML error reading theme file",                                         // custom message + err.Error()
	"0117": "Unable to get shared stylesheets",                                      // Name of TOML file
	"0118": "Error reading configuration",                                           //
	"0119": "inc: unable to open location",                                          // location
	"0120": "inc: unable to open file",                                              // filename
	"0121": "inc: error reading file",                                               // filename
	"0122": "scode: unable to find file",                                            // filename
	"0123": "scode: error reading file",                                             // filename
	"0124": "Error copying a page asset",                                            // custom message
	"0125": "Error copying a style sheet",                                           // custom message
	"0126": "Error reading config file " + productName + "." + configFileDefaultExt, // Viper runtime error

	// 0200	- Error creating file
	"0201": "Error creating site configuration file",             // err.Error
	"0202": "Error closing site configuration file",              // err.Error
	"0203": "Error writing output file",                          // custom message + err.Error()
	"0204": "Error writing output file",                          // custom message
	"0205": "Error creating theme TOML file",                     // custom message
	"0206": "Error renaming CSS file for new theme",              // custom message
	"0207": "Error copying and renaming TOML file for new theme", // custom message
	"0208": "Error copying file to publish directory",            // Golang error, custom message
	"0209": "Error creating file",                                // filename
	"0210": "Error creating TOML file for new site",              // filename
	"0211": "Error creating sample file",                         // filename
  "0212": "Error renaming temporary output file",               // filename

	// 0250 - Error closing file
	"0251": "Error closing copy of file", // filename
	"0252": "Error closing TOML file",    // filename

	// 0300	- Error deleting file
	"0301": "Unable to delete theme file",          // custom message
	"0302": "Error deleting old publish directory", // directory name

	// 0400	- Error creating directory
	"0401": "Error creating site directory",    // Name
	"0402": "Error publishing asset directory", // Name
	"0403": "Error recreating publish directory",
	"0404": "Error creating publish directory",   // dir name
	"0405": "Error recreating publish directory", // dir name
	"0406": "Error copying directory",            // custom msg with both dirs, Golang error,
	"0407": "Error creating site directory",      // dir name

	// 0500	- Error determining directory name
	"0501": "",

	// 0600 - Error deleting directory
	"0601": "",

	// 0700	- Error reading directory
	"0701": "Can't copy theme",                      // Custom msg
	"0702": "Error copying directory",               // Go error + dir name
	"0703": "Error reading directory",               // directory name
	"0704": "Missing name of source directory",      // directory name to copy from
	"0705": "Missing name of target directory",      // directory name to copy to
	"0706": "Unable to read from headers directory", // Name of headers directory
	"0707": "Directories are identical:",            // custom message

	// 0800	- Can't determine the name of something
	"0801": "",

	// 0900	- Problem generating something
	"0901": "Problem creating TOML object",                         // err.Error
	"0902": "Error creating new site.toml file",                    // Full custom error message
	"0903": "Error writing to file",                                // Full custom error message
	"0904": "Theme name taken",                                     // custom message
	"0905": "Couldn't create directory for new theme",              // custom message
	"0906": "Problem creating new theme files",                     // custom message
	"0907": "Pagetype name taken",                                  // custom message
	"0908": "Problem creating TOML object",                         // runtime error
	"0910": "Problem creating output file",                         // filename
	"0911": "Unable to copy themes directory to site directory",    // custom message
	"0912": "Problem converting markdown file",                     //
	"0913": "Unable to read project directory",                     //
	"0914": "Error creating a temporary file",                      // filename
	"0915": "Unable to copy scodes directory to site directory",    // custom message
	"0916": "Unable to copy a style sheet",                         // custom message
	"0917": "Error creating ",                                      // filename, Golang message
	"0918": "Can't copy a theme onto itself. That would be silly.", //
	"0919": "Pagetype name is already taken",                       // custom message

	// 0950 - Something's already there
	"0951": "Site already exists:", // sitename

	// 1000	- Something's missing that should be there
	"1001": "Missing front matter and markdown", // filename
	"1002": "This isn't a project directory",    // full custom message
	"1003": "Unable to find theme",              // full custom message
	"1004": "The",                               // full custom message
	"1005": "PageType not found",                // full custom message
	"1006": "PageType not found",                // full custom message
	"1007": "Error reading theme",               // full custom message
	"1008": "No theme file by the name",         // custom message
	"1009": "Not a project directory",           //  Dir name
	"1010": "Not a theme file",                  // full custom message
	"1011": "No Publish directory specified",    //
	"1012": "Please specify a site name",
	"1013": "Please specify a site name",
	"1014": "Unable to determine application configuration data directory",
	"1015": "Theme TOML specifies a file that can't be found", // filename
	"1016": "Unable to read directory",                        // filename
	"1017": "Missing name of the theme to create.",
  "1018": "No path specified for the project",                // 

	// 1100 - Problem changing to a directory
	"1101": "Can't change to source directory", // directory name
	"1102": "Can't change to source directory", // directory name
	"1103": "Can't change to site directory",   // directory name
	"1104": "Can't change to site directory",   // directory name
	"1105": "Can't change to site directory",   // directory name
	"1106": "Can't change to site directory",   // directory name

	// 1200 - Syntax error!
	"1201": "inc: Couldn't execute template in",   // filename
	"1202": "scode: Couldn't execute template in", // filename
	"1203": "Unknown tag type",                    // Tag name
	"1204": "Error executing",                     // filename, Go error message
}

type errMsg struct {
	// Key to a map of error messages
	key string

	// If key is the word "PREVIOUS", this will contain an error
	// message from an earlier action, typically a return from the
	// Go runtime.
	previous string
	extra    []string
}

// Error() looks up e.key, which is an error code expressed as
// a string (for example, "1001") and returns its associated map entry.
// But there's likely much more happening:
// -  If e.key is "PREVIOUS" it suggests that an error message
//    that didn't get displayed probably
//    should, and its contents in e.previous are returned.
// -  If e.extra has something, say, a filename, it should be
//    used to customize the error message.
// Why the fake number? Because it gets appended to "mbz" in an error message,
// and I plan for Metabuzz to be so popular that people would be
// looking up error codes search engines, e.g. mbz1001. And it's a
// ghetto way of keeping error codes unique while being kind of sorted
// in the source code.
func (e *errMsg) Error() string {
	var msg error
	// Make sure the error code has documentation
	if errMsgs[e.key] == "" {
		msg = fmt.Errorf("ERROR CODE %s UNTRACKED: please contact the developer\nMore info: %s\n",
			errorCodePrefix+e.key, e.previous)
		return msg.Error()
	}

	// Error message from an earlier error return needs to be seen.
	if e.key == "PREVIOUS" {
		return fmt.Errorf("%s\n", e.previous).Error()
	}

	if e.previous != "" {
		msg = fmt.Errorf("%s %s (error code %s%s)\n",
			// The most common case: an error code with customization
			errMsgs[e.key], e.previous, errorCodePrefix, e.key)
	} else {
		msg = fmt.Errorf("%s (error code %s%s)\n",
			// The slighly unusual case of an error code with no customization
			errMsgs[e.key], errorCodePrefix, e.key)
	}
	return msg.Error()
}

// New() allocates a map entry for errMsgs.
func New(key string, previous string, extra ...string) error {
	return &errMsg{key, previous, extra}
}

// errCode() takes an error code, say "0110", and
// one or two optional strings. It adds the error code
// to the error messages map so that message can be looked
// up. The additional parameters give information such
// as whether a previous error message should be displayed,
// or something to customize the message that's already in
// the error messages map, like a filename.
// When calling a Go runtime routine that could return
// an error message, make err.Error() the second
// parameter so its contents are included, like this:
/*
	f, err := os.Create(filename)
	if err != nil {
		return errCode("0210", err.Error(), filename)
	}
*/

func errCode(key string, previous string, extra ...string) error {
	var e error
	if len(extra) > 0 {
		e = New(key, previous, extra[0])
	} else {
		e = New(key, previous)
	}
	return e
}

// App.Verbose() displays a message followed
// by a newline to stdout
// if the verbose flag was used. Formats it like Fprintf.
func (App *App) Verbose(format string, a ...interface{}) {
	if App.Flags.Verbose {
		fmt.Println(App.fmtMsg(format, a...))
	}
}

// App.Warning() displays a message followed by a newline
// to stdout, preceded by the text "Warning: "
// Overrides the verbose flag. Formats it like Fprintf.
func (App *App) Warning(format string, a ...interface{}) {
	fmt.Println("Warning: " + App.fmtMsg(format, a...))
}

// fmtMsg() formats string like Fprintf and writes to a string
func (App *App) fmtMsg(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

func displayErrCode(errCode string) {
}

// displayError() shows the specified error message
// without exiting to the OS.
func displayError(e error) {
	fmt.Println(e.Error())
}

// QuitError() displays the error passed to it and exits
// to the operating system, returning a 1 (any nonzero
// return means an error ocurred).
// Normally functions that can generate a runtime error
// do so by returning an error. But sometimes there's a
// constraint, for example, fulfilling an interface method
// that doesn't support this practice.
func (App *App) QuitError(e error) {
	if App.Page.filePath != "" {
		fmt.Printf("%s ", App.Page.filePath)
	}
	displayError(e)
	if e == nil {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
