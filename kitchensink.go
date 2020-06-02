package main

import (
	"fmt"
	"path/filepath"
)

type siteDescription struct {
	filename string
	dir      string
	mdtext   string
}

/* SVG file of an exciting 100x100px gray box */
var svgFile = `<?xml version="1.0" encoding="utf-8"?>
<svg  xmlns="http://www.w3.org/2000/svg">
  <rect x="0" y="0" width="100" height="100" style="fill: rgb(216, 216, 216);"/>
</svg>
`

var (
	siteTest = []siteDescription{
		{"index.md",
			"",
			`# Home
Go [one level deep](one/index.html), [two levels deep](two/three/index.html)

Host: {{ hostname }}

Time: {{ ftime }}

Location of this file: {{ path }}


**Box**

![100x100 SVG box](box-100x100.svg)

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

Location of this file: {{ path }}


Go [home 1](/index.html)

Go [home 2](\/index.html)

Go [home 3](/)

Go [home 4](/./index.html)

`},
	}

	// Create a test site to exercise important features
	// given a filename, the path to that filename,
	// and the Markdown text itself.
	// This probably won't end well but I can't think of a better
	// way to do this with limited time.
	testPages = []struct {
		filename string
		dir      string
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

// writeSiteFromArray() takes an array of
// structures containing a filename,
// a path to that filename, and the markdown
// text itself, and writes them out to
// a test site.
func writeSiteFromArray(sitename string, site []siteDescription) error {
	// First put an SVG graphic in the root
	//path := filepath.Join(site[0].dir, site[0].filename)
	path := filepath.Join(site[0].dir, "box-100x100.svg")
	err := writeTextFile(path, svgFile)
	if err != nil {
		return errCode("0211", err.Error(), "Sample SVG file")
		//return errCode("0211", "Sample SVG file")
	}
	for _, f := range site {
		path := filepath.Join(f.dir, f.filename)
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
		App.QuitError(err)
	}

	// Create directory structure for test site
	if err := createDirStructure(&testDirs); err != nil {
		return err
	}

	// Build the site from the array of data structures
	if err := writeSiteFromArray(sitename, siteTest); err != nil {
		return err
	}

	fmt.Println("Created site ", App.Site.Name)
	return nil

}
