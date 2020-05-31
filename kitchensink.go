package main

import (
	"fmt"
	"path/filepath"
)

type siteDescription struct {
		filename string
		dir string
		mdtext   string
}


var (


	siteTest = []siteDescription{
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

)
/*

   err := App.newSite(App.Site.Name)
   if err != nil {
     QuitError(err)
   } else {
     fmt.Println("Created site ", App.Site.Name)
   }
*/

// writeSiteFromArray() takes an array of 
// structures containing a filename, 
// a path to that filename, and the markdown
// text itself, and writes them out to
// a test site.
func writeSiteFromArray(sitename string, site []siteDescription) error {
  for _, f := range site {
    path := filepath.Join(f.dir,f.filename)
    err := writeTextFile(path, f.mdtext)
    if err != nil {
      return errCode("PREVIOUS", err.Error(), path)
    }
  }
  return nil
}


// kitchenSink() Generates a test site from an
// array of structures containing a filename, 
// a path to that filename, and the markdown
// text itself.
func (App *App) kitchenSink(sitename string) error {
	err := App.newSite(sitename)
	if err != nil {
		QuitError(err)
	}

	// Create directory structure for test site
	if err := createDirStructure(&testDirs); err != nil {
		return err
	}

  if err := writeSiteFromArray(sitename, siteTest); err != nil {
    return err
  }

  var path string
	for _, t := range testPages {
	  // filename,dir, mdtext
    //fmt.Printf("%s/%s\n", testPages[each].dir, testPages[each].filename)
    ///fmt.Printf("%s/%s\n", t.dir, t.filename)
    path = filepath.Join("/",t.dir,t.filename)
    fmt.Printf("%s\n%v------\n\n",path, t.mdtext)
  }

	fmt.Println("Created site ", App.Site.Name)
	return nil

}
