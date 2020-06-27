package app

// NewTestApp builds a new bare-bones app with just enough initialization for
// testing.
func NewTestApp(mdSrc string) *App {
	return &App{
		Site: &Site{},
		Page: &Page{
			markdownStart: []byte(mdSrc),
		},
	}
}
