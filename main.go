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
	fmt.Println("VERSION 4:30pm")
	fmt.Println("* Document the v validator. Write an article about it")
	fmt.Println("* Put something in the headtags directory")
}
