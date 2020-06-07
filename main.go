package main

import (
"fmt"
)

func main() {
  fmt.Println("* Rootstylesheets are trying to copy remote stylesheet")
	// Get the execution evironment ready.
	App := newDefaultApp()
	// Read configuration files, environment, and command line.
	if err := App.Cmd.Execute(); err != nil {
		App.QuitError(err)
	}
}
