package main

// newApp builds a new bare-bones app with just enough initialization for
// testing.
func newApp(mdSrc string) *App {
	return &App{
		Site: &Site{},
		Page: &Page{
			markdownStart: []byte(mdSrc),
		},
	}
}
