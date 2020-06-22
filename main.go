package main

import (
"fmt"
)

func main() {
	// Get the execution environment ready.
	App := newDefaultApp()
	// Read configuration files, environment, and command line.
	if err := App.Cmd.Execute(); err != nil {
		App.QuitError(err)
	}
  fmt.Println("* Put something in the headtags directory") }
