package main

import (
	"bytes"
	//"fmt"
	//"github.com/gohugoio/hugo/markup/tableofcontents"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	//"regexp"
	//"text/template"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// publishFile() is the heart of this program. It converts
// a Markdown document (with optional TOML at the beginning)
// to HTML.
func (App *App) publishFile(filename string) error {
	var input []byte
	var err error
	App.Page.filePath = filename
	App.Page.filename = filepath.Base(filename)
	App.Page.dir = currDir()
	App.Verbose("%s", filename)
	// Read the whole Markdown file into memory as a byte slice.
	input, err = ioutil.ReadFile(filename)
	if err != nil {
		return errCode("0102", filename)
	}

	// Extract front matter and parse.
	// Starting at the Markdown, convert to HTML.
	// Interpret templates as well.
	App.Convert(filename, input)

	// Output filename
	outfile := replaceExtension(filename, "html")
	// Strip out everything but the filename.
	//base := filepath.Base(htmlkjlfile)

	// Write everything to a temp file so in case there was an error, the
	// previous HTML file is preserved.
	tmpFileBaseName := PRODUCT_NAME + "-tmp-"
	tmpFile, err := ioutil.TempFile(App.Site.Publish, tmpFileBaseName)
	// Ensure the file gets closed before exiting
	defer os.Remove(tmpFile.Name())
	if err != nil {
		return errCode("0914", err.Error())
	}

	// This is why we're here! Translate from Markdown to HTML
	err = writeTextFile(tmpFile.Name(), string(App.Page.Article))
	if err != nil {
		return errCode("PREVIOUS", "")
	}

	// Copy any associated assets such as
	// images in the same directory.
	//dirHasMarkdownFiles := App.publishLocalFiles()
	_ = App.publishLocalFiles(App.Page.dir)

	// Create the output filename by replacing the name of the input file with an html extension
	//outfile := replaceExtension(filename, "html")
	// Strip out everything but the filename.
	base := filepath.Base(outfile)

	// Write everything to a temp file so in case there was an error, the
	// previous HTML file is preserved.
	tmpFile, err = ioutil.TempFile(App.Site.Publish, PRODUCT_NAME+"-tmp-")
	writeTextFile(tmpFile.Name(), string(App.Page.Article))
	// Ensure the file gets closed before exiting
	defer os.Remove(tmpFile.Name())
	// Get the relative directory.
	relDir := relDirFile(App.Site.path, outfile)
	App.Page.Path = relDir
	// If there's a README.md and no index.md, rename
	// the output file to index.html
	if App.Page.filename == "README.md" && !optionSet(App.Site.dirs[App.Page.dir], hasIndexMd) {
		base = "index.html"
	}

	// Generate the full pathname of the matching output file, as it will
	// appear in its published location.
	outfile = filepath.Join(App.Site.Publish, relDir, base)
	// If the write succeeded, rename it to the output file
	// This way if there was an existing HTML file but there was
	// an error in output this time, it doesn't get clobbered.
	if err = os.Rename(tmpFile.Name(), outfile); err != nil {
		return err
	}

	if !fileExists(outfile) {
		QuitError(errCode("0910", outfile))
	}
	App.Verbose("\tCreated file %s", outfile)
	App.fileCount++
	//
	// Success
	return nil
}

// Takes a buffer containing Markdown
// and converts to HTML.
// Doesn't know about front matter.
func (App *App) markdownBufferToBytes(input []byte) []byte {
	autoHeadings := parser.WithAttribute()
	if App.Site.MarkdownOptions.headingIDs == true {
		autoHeadings = parser.WithAutoHeadingID()
	}

	if App.Site.MarkdownOptions.hardWraps == true {
		// TODO: Figure out how to get this damn thing in as an option
		//hardWraps := html.WithHardWraps()
	}

	// TODO: Make the following option: Footnote,
	markdown := goldmark.New(
		goldmark.WithExtensions(extension.GFM, extension.DefinitionList, extension.Footnote,
			highlighting.NewHighlighting(
				highlighting.WithStyle(App.Site.MarkdownOptions.HighlightStyle),
				highlighting.WithFormatOptions(
				//highlighting.WithLineNumbers(),
				),
			)),
		goldmark.WithParserOptions(
			autoHeadings, // parser.WithAutoHeadingID(),
			parser.WithAttribute(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
			/* html.WithHardWraps(), */
			html.WithXHTML(),
		),
	)

	var buf bytes.Buffer
	if err := markdown.Convert(input, &buf); err != nil {
		// TODO: Need something like displayErrCode("1010") or whatever
		App.Warning("Error converting Xxx")
		return []byte{}
	}
	return buf.Bytes()
}

// appendBytes() Appends a byte slice to the byte slice containing the rendered output
func (App *App) appendBytes(b []byte) {
	App.html = append(App.html, b...)
}

// appendStr() Appends a string to the byte slice containing the rendered output
func (App *App) appendStr(s string) {
	App.html = append(App.html, s...)
}

// MdFileToHTMLBuffer() takes a byte slice buffer containing
// a pure Markdown file as input, and returns
// a byte slice containing the file converted to HTML.
// It doesn't know about front matter.
// So it should be preceded by a call to App.parseFrontMatter()
// if there's any possibility that the file contains front matter.
// In the spirit of a browser, it simply returns an empty buffer on error.
func (App *App) MdFileToHTMLBuffer(filename string, input []byte) []byte {
	// Resolve any Go template variables before conversion to HTML.
	interp := App.interps(filename, string(input))
	// Convert markdown to HTML.
	return App.markdownBufferToBytes([]byte(interp))
}

// Convert() takes a document with optional front matter, parses
// out the front matter, and sends the Markdown portion to be converted.
// Write the HTML results to App.Page.Article
func (App *App) Convert(filename string, input []byte) (start []byte, err error) {
	// Extract front matter and parse.
	// Return the starting address of the Markdown.
	start, err = App.parseFrontMatter(filename, input)
	if err != nil {
		return []byte{}, errCode("0103", filename)
	}
	// Resolve any Go template variables before conversion to HTML.
	interp := App.interps(filename, string(start))
	App.Page.Article = App.markdownBufferToBytes([]byte(interp))
	return App.Page.Article, nil
}

// publishLocalFiles() get called for every markdown file
// in the directory. It copies assets like image files & so forth
// from the source file's current directory to the publish location,
// creating a new subdirectory as needed.
// For example, if your article references ![cat](cat.png)
// then presumbably cat.png is in the current directory.
// This copies all nonexcluded files, such as cat.png and
// any other assets, from this directory
// into its matching publish directory,
// same as the source markdown file.
// Creates a subdirectory under Publish if in a subdirectory
// and one hasn't yet been created.
// Keeps track of which directories have had their assets copied to
// avoid redundant copies.
// Returns true if there are any markdown files in the current directory.
// Returns false if markdown files (or any files at all) are abset.
func (App *App) publishLocalFiles(dir string) bool {
	relDir := relativeDirectory(App.Site.path, dir)
	pubDir := filepath.Join(App.Site.Publish, relDir)

	if !optionSet(App.Site.dirs[pubDir], markdownDir) {
		if err := os.MkdirAll(pubDir, PUBLIC_FILE_PERMISSIONS); err != nil {
			// TODO: Have this function return error?
			App.infoLog.Printf(errCode("0404", pubDir).Error())
		}
		App.Site.dirs[dir] |= markdownDir
	}
	// Get the directory listing.
	candidates, err := ioutil.ReadDir(dir)
	if err != nil {
		// TODO: Return error?
		App.infoLog.Println("publishLocalFiles(): NO files found")
		return false
	}

	// Get list of files in the local directory to exclude from copy
	var excludeFromDir = searchInfo{
		list:   App.FrontMatter.ExcludeFilenames,
		sorted: false}

	// First check the directory to ensure there's at least 1 markdown file.
	hasMarkdown := false

	// Look for the specific file README.md, which competes with
	// index.md:
	// https://stackoverflow.com/questions/58826517/why-do-some-static-site-generators-use-readme-md-instead-of-index-md
	for _, file := range candidates {
		filename := file.Name()
		if hasExtensionFrom(filename, markdownExtensions) {
			hasMarkdown = true
		}

		if filename == "README.md" {
			App.Site.dirs[dir] |= hasReadmeMd
		}
		if strings.ToLower(filename) == "index.md" {
			App.Site.dirs[dir] |= hasIndexMd
		}

	}

	if hasMarkdown {
		// Flag this as a directory that contains at least
		// 1 markdown file.
		App.Site.dirs[dir] |= markdownDir
	} else {
		// No markdown files found, so return
		return false
	}
	for _, file := range candidates {
		filename := file.Name()
		// Don't copy if it's a directory.
		if !file.IsDir() {
			// Don't copy if its extension is on one of the excluded lists.
			if !hasExtension(filename, ".css") &&
				!hasExtensionFrom(filename, markdownExtensions) &&
				!excludeFromDir.Found(filename) &&
				!strings.HasPrefix(filename, ".") {
				// It's a markdown file.
				// Got the file. Get its fully qualified name.
				copyFrom := filepath.Join(dir, filename)
				// Figure out the target directory.
				relDir := relDirFile(App.Site.path, copyFrom)
				// Get the target file's fully qualified filename.
				copyTo := filepath.Join(App.Site.Publish, relDir, filename)
				if err := Copy(copyFrom, copyTo); err != nil {
					//QuitError(err.Error())
					QuitError(errCode("PREVIOUS", err.Error()))
				}
			}
		}
	}
	return true
}
