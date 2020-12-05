package app

import (
	"fmt"
	"os"
)

// App.Verbose() displays a message followed
// by a newline to stdout
// if the verbose flag was used. Formats it like Fprintf.
func (a *App) Verbose(format string, ss ...interface{}) {
	if a.Flags.Verbose {
		fmt.Println(a.fmtMsg(format, ss...))
	}
}

// App.Warning() displays a message followed by a newline
// to stdout, preceded by the text "Warning: "
// Overrides the verbose flag. Formats it like Fprintf.
func (a *App) Warning(format string, ss ...interface{}) {
	fmt.Println("Warning: " + a.fmtMsg(format, ss...))
}

// fmtMsg() formats string like Fprintf and writes to a string
func (a *App) fmtMsg(format string, ss ...interface{}) string {
	return fmt.Sprintf(format, ss...)
}

// displayError() shows the specified error message
// without exiting to the OS.
func displayError(e error) {
	fmt.Println(e.Error())
}

// QuitError() displays the error passed to it and exits
// to the operating system, returning a 1 (any nonzero
// return means an error occurred).
// Normally functions that can generate a runtime error
// do so by returning an error. But sometimes there's a
// constraint, for example, fulfilling an interface method
// that doesn't support this practice.
func (a *App) QuitError(e error) {
	if a.Page.filePath != "" {
		fmt.Printf("%s ", a.Page.filePath)
	}
	displayError(e)
	if e == nil {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

// Quit() is a quick and dirty function that displays the
// message and exits to the operating system.
// TODO: Don't allow these in production!
func (a *App) Quit(format string, ss ...interface{}) {
	fmt.Println(a.fmtMsg(format, ss...))
	os.Exit(1)
}
