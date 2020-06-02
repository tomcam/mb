package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// build() creates the site.
// It assumes it is in the source directory.
// Assets in the theme/pagetype directories are published, which
// includes anything other than HTML or Markdown files.
func (App *App) build() error {
	if !isProject(".") {
		return errCode("1009", currDir())
	}
  App.Site.path = currDir()
	App.siteDefaults()

	var err error
	// Delete any existing publish dir
	if err := os.RemoveAll(App.Site.Publish); err != nil {
		return errCode("0302", App.Site.Publish)
	}
	// Now create an empty publish dir
	err = os.MkdirAll(App.Site.Publish, PUBLIC_FILE_PERMISSIONS)
	if err != nil {
		return errCode("0403", App.Site.Publish)
	}

	// Get a list of all files & directories in the site.
	if _, err = App.getProjectTree(App.Site.path); err != nil {
		return errCode("0913", "")
	}

	// Loop through the list of permitted directories for this site.
	for dir, _ := range App.Site.dirs {
		// Change to each directory
		if err := os.Chdir(dir); err != nil {
			return errCode("1101", dir)
		}
		// Get the files in just this directory
		files, err := ioutil.ReadDir(".")
		if err != nil {
			return errCode("0703", dir)
		}

		// Go through all the Markdown files and convert.
		for _, file := range files {
			if !file.IsDir() && isMarkdownFile(file.Name()) {
				if err := App.publishFile(filepath.Join(dir, file.Name())); err != nil {
					return errCode("PREVIOUS", err.Error())
				}
			}
		}
	}
	fmt.Printf("%v ", App.fileCount)
	if App.fileCount != 1 {
		fmt.Println("files")
	} else {
		fmt.Println("file")
	}

	return nil
}
