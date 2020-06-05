package main
import(
  "fmt"
)
func main() {
  fmt.Println("* Work on new pagetype next")

	// Get the execution evironment ready.
	App := newDefaultApp()
	// Read configuration files, environment, and command line.
	if err := App.Cmd.Execute(); err != nil {
		App.QuitError(err)
	}
}
