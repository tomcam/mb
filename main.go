package main

import (
//"fmt"
)

func main() {
  fmt.Println("* Put something in the headtags directory")
	// Get the execution environment ready.
	App := newDefaultApp()
	// Read configuration files, environment, and command line.
	if err := App.Cmd.Execute(); err != nil {
		App.QuitError(err)
	}
}
