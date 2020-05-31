package main

import (
	"fmt"
	"path/filepath"
)

var (
  // Create a test site to exercise important features
  // given a filename, the path to that filename,
  // and the Markdown text itself. 
  // This probably won't end well but I can't think of a better
  // way to do this with limited time.
	testPages = []struct {
		filename string
		dir string
		mdtext   string
	}{
		{"index.md",
			"",
			`# Home
Go [one level deep](one/index.html), [two levels deep](two/three/index.html)
`},
		{"index.md",
			"one",
			`# Page 1
This page is 1 level deep.

The time is {{ ftime }}
`},
		{"index.md",
			"two/three",
			`# Page 2
This page is 2 levels deep.

Go [home](/index.html)
`},
	}

	/* Directory structure for the test site */
	testDirs = [][]string{
		{"one"},
		{"two", "three"},
	}

	rootMd = `# Home page at root

Link 1 directory deep: [one](one/index.html)

Link 2 directories deep: [three](three/index.html)
`
)

/*

   err := App.newSite(App.Site.Name)
   if err != nil {
     QuitError(err)
   } else {
     fmt.Println("Created site ", App.Site.Name)
   }
*/

// Generates a test app
func (App *App) kitchenSink(sitename string) error {
	err := App.newSite(sitename)
	if err != nil {
		QuitError(err)
	}

	// Create directory structure for test site
	if err := createDirStructure(&testDirs); err != nil {
		return err
	}

  var path string
	for _, t := range testPages {
	  // filename,filepath mdtext
    //fmt.Printf("%s/%s\n", testPages[each].dir, testPages[each].filename)
    ///fmt.Printf("%s/%s\n", t.dir, t.filename)
    path = filepath.Join("/",t.dir,t.filename)
    fmt.Printf("%s\n%v------\n\n",path, t.mdtext)
  }

	fmt.Println("Created site ", App.Site.Name)
	return nil

}
