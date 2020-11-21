package app

import (
	"fmt"
	"github.com/tomcam/mb/pkg/defaults"
	"github.com/tomcam/mb/pkg/errs"
	"io/ioutil"
	"os"
	"path/filepath"
)

// build() creates the site.
// It assumes it is in the source directory.
// Assets in the theme/pagetype directories are published, which
// includes anything other than HTML or Markdown files.
func (a *App) build() error {
	if !isProject(".") {
		return errs.ErrCode("1009", currDir())
	}

	var err error
	a.SiteDefaults()
	// Delete any existing publish dir
	if err := os.RemoveAll(a.Site.Publish); err != nil {
		return errs.ErrCode("0302", a.Site.Publish)
	}

	// Create an empty publish dir
	if err := os.MkdirAll(a.Site.Publish, defaults.PublicFilePermissions); err != nil {
		return errs.ErrCode("0403", a.Site.Publish)
	}

	if a.Site.path == "" {
		return errs.ErrCode("1018", "")
	}

	// Get a list of all files & directories in the site.
	if _, err = a.getProjectTree(a.Site.path); err != nil {
		return errs.ErrCode("0913", a.Site.path)
	}

	// Loop through the list of permitted directories for this site.
	for dir := range a.Site.dirs {
		// Change to each directory
		if err := os.Chdir(dir); err != nil {
			return errs.ErrCode("1101", dir)
		}
		// Get the files in just this directory
		files, err := ioutil.ReadDir(".")
		if err != nil {
			return errs.ErrCode("0703", dir)
		}

		// Go through all the Markdown files and convert.
		// Start search index JSON file with opening '['
		a.DelimitIndexJSON(a.Site.SearchJSONFilePath, true)
		commaNeeded := false
		for _, file := range files {
			if !file.IsDir() && isMarkdownFile(file.Name()) {
				// It's a Markdown file, not a dir or anything else.
				if commaNeeded {

					// TODO: Add error checking
					a.AddCommaToSearchIndex(a.Site.SearchJSONFilePath)
					commaNeeded = false
				}
				if err = a.publishFile(filepath.Join(dir, file.Name())); err != nil {
					return errs.ErrCode("PREVIOUS", err.Error())
				}
				commaNeeded = true
			}
		}
		// Close search index JSON file with ']'
		a.DelimitIndexJSON(a.Site.SearchJSONFilePath, false)

	}
	fmt.Printf("%v ", a.fileCount)
	if a.fileCount != 1 {
		fmt.Println("files")
	} else {
		fmt.Println("file")
	}

	return nil
}
