package main

import (
	"github.com/tomcam/mb/pkg/app"
)

func main() {
	// Get the execution environment ready.
	App := app.NewDefaultApp()
	// Read configuration files, environment, and command line.
	if err := App.Cmd.Execute(); err != nil {
		App.QuitError(err)
	}
}
