package main

import (
	"bytes"
	"fmt"
	//"github.com/gohugoio/hugo/markup/tableofcontents"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	//"text/template"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var (
	// Credit to anonymous user at:
	// https://play.golang.org/p/OfQ91QadBCH
	// Match an h1 in Markdown
	h1, _ = regexp.Compile("(?m)^\\s*#{1}\\s*([^#\\n]+)$")
	// Match headers 2-6 in Markdown
	anyHeader, _ = regexp.Compile("(?m)^\\s*#{2,6}\\s*([^#\\n]+)$")
	// Match everything after the pound sign on a line starting with the pound sign
	notPound, _ = regexp.Compile("(?m)[^#|\\s].*$")

	closingHTMLTags = `
</body>
</html>
`
)

// publishFile() is the heart of this program. It converts
// a Markdown document (with optional TOML at the beginning)
// to HTML.
func (App *App) publishFile(filename string) error {
	var input []byte
	var err error
	// Get a fresh new Page object if doing more
	// than one file at a clip. Which is obviously
	// most of the time.
	var p Page
	App.Page = &p
	var f FrontMatter
	App.FrontMatter = &f
	// Probably belongs in build.go if anything
	//App.siteDefaults()
	App.Page.filePath = filename
	App.Page.filename = filepath.Base(filename)
	App.Page.dir = currDir()
	App.Verbose("%s", filename)
	// Read the whole Markdown file into memory as a byte slice.
	input, err = ioutil.ReadFile(filename)
	if err != nil {
		return errCode("0102", filename)
	}
	// TODO: Wrong. Setting default them properly means checking site.
	App.FrontMatter.PageType = ""

	// Extract front matter and parse.
	// Obviously that includes an optional theme or pagetype designation.
	// Starting at the Markdown, convert to HTML.
	// Interpret templates as well.
	App.Page.markdownStart, err = App.parseFrontMatter(filename, input)
	if err != nil {
		return errCode("0103", filename)
	}

	// If no theme was specified in the front matter, but one was specified in the
	// site config, make the one specified in site.toml the theme.
	if App.Site.Theme != "" && App.FrontMatter.Theme == "" {
		App.FrontMatter.Theme = App.Site.Theme
	}
	App.loadTheme()
	// Parse front matter.
	// Convert article to HTML
	App.Article(filename, input)
	// xxx
	// Begin HTML document.
	// Open the <head> tag.
	App.startHTML()

	// If a title wasn't specified in the front matter,
	// Generate title tag contents from headers. Find the first
	// h1. If that fails, find the first h2-h6. If that
	// fails, put up a self-aggrandizing error message.
	App.titleTag()

	App.descriptionTag()

	App.headTags()

	// Output filename
	outfile := replaceExtension(filename, "html")
	relDir := relDirFile(App.Site.path, outfile)

	// START ASSEMBLING PAGE
	App.appendStr(App.Page.startHTML)
	App.appendStr(wrapTag("<title>", App.Page.titleTag, true))
	App.appendStr(metatag("description", App.Page.descriptionTag))
	App.appendStr(App.Page.headTags)

	// Hoover up any miscellanous files lying around,
	// like other HTML files, graphic assets, etc.
	App.localFiles(relDir)

	// Write the closing head tag and the opening
	// body tag
	App.closeHeadOpenBody()

	//App.appendStr(wrapTag("<article>", []byte(App.Page.Article), true))
	App.appendStr(App.pageRegionToHTML(&App.Page.Theme.PageType.Header, "<header>"))
	App.appendStr(App.pageRegionToHTML(&App.Page.Theme.PageType.Nav, "<nav>"))
	App.appendStr(App.pageRegionToHTML(&App.Page.Theme.PageType.Article, "<article>"))
	sidebar := strings.ToLower(App.FrontMatter.Sidebar)
	if sidebar == "left" || sidebar == "right" {
		App.appendStr(App.pageRegionToHTML(&App.Page.Theme.PageType.Sidebar, "<aside>"))
	}
	App.appendStr(App.pageRegionToHTML(&App.Page.Theme.PageType.Footer, "<footer>"))

	// Complete the HTML document with closing <body> and <html> tags
	App.appendStr(closingHTMLTags)

	// Strip out everything but the filename.
	base := filepath.Base(outfile)

	// Write everything to a temp file so in case there was an error, the
	// previous HTML file is preserved.
	tmpFile, err := ioutil.TempFile(App.Site.Publish, productName+"-tmp-")
	if err != nil {
		App.QuitError(errCode("PREVIOUS", err.Error()))
	}

	if err = writeTextFile(tmpFile.Name(), string(App.Page.html)); err != nil {
		App.QuitError(errCode("PREVIOUS", err.Error()))
	}
	// Ensure the file gets closed before exiting
	defer os.Remove(tmpFile.Name())
	// Get the relative directory.
	//relDir = relDirFile(App.Site.path, outfile)
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
		App.QuitError(errCode("0212", outfile))
	}

	if !fileExists(outfile) {
		App.QuitError(errCode("0910", outfile))
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
	App.Page.html = append(App.Page.html, b...)
}

// appendStr() Appends a string to the byte slice containing the rendered output
func (App *App) appendStr(s string) {
	App.Page.html = append(App.Page.html, s...)
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

	// If this directory hasn't been created, create it.
	if !optionSet(App.Site.dirs[pubDir], markdownDir) {
		if err := os.MkdirAll(pubDir, PUBLIC_FILE_PERMISSIONS); err != nil {
			App.QuitError(errCode("0404", pubDir, err.Error()))
		}
		// Mark that the directory has been created so this
		// doesn't get repeated.
		App.Site.dirs[dir] |= markdownDir
	}
	// Get the directory listing.
	candidates, err := ioutil.ReadDir(dir)
	if err != nil {
		App.QuitError(errCode("1016", dir, err.Error()))
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
				//App.Verbose("\tCopying '%s' to '%s'\n", copyFrom, copyTo)
				if err := Copy(copyFrom, copyTo); err != nil {
					//App.QuitError(err.Error())
					App.QuitError(errCode("PREVIOUS", err.Error()))
				}
				// TODO: Get rid of fileExists() when possible
				if !fileExists(copyTo) {
					App.QuitError(errCode("PREVIOUS", copyTo))
				}
			}
		}
	}
	return true
}

// publishPageTypeAssets() figures out what assets are used for this pageType, namely
// stylesheets, graphics, and anything that's not HTML or markdown.
// Ideally these assets are sparingly used, for example, a logo for a header.
// In the spirit of HTML, missing assets are ignored.
// There's a simple form of inheritance. If you publish a PageType that's
// not the anonymous one, it means it's a child, so copy the parent assets
// as well. (You can exclude files using Exclude in the theme's toml file.)
func (App *App) publishPageTypeAssets() {
	// Is the default aka parent theme?
	p := App.Page.Theme.PageType
	if p.name == "" || p.name == defaultThemeName {
		// Default PageType
		App.publishAssets()
	} else {
		// It's a child theme aka pagetype.
		App.Page.Theme.PageType = p
		App.publishAssets()
		// TODO:This would happen more than once
		// with multiple pageTypes, so I should
		// eliminate that.
	}
}

// getMode() checks if the stylesheet is dark or light and adjusts as needed
func (App *App) getMode(stylesheet string) string {
	if stylesheet == "theme-light.css" && App.FrontMatter.Mode == "dark" {
		stylesheet = "theme-dark.css"
	}
	return stylesheet
}

// publishAssets() copies out the stylesheets, graphics, and other
// relevant files from the pageType (or default theme) directory
// to be published.
func (App *App) publishAssets() {
	p := App.Page.Theme.PageType
	App.publishPageAssets()
	App.publishThemeAssets()
	// Copy out different stylesheet depending on the
	// type of sidebar, if any.
	switch strings.ToLower(App.FrontMatter.Sidebar) {
	case "left":
		p.Stylesheets = append(p.Stylesheets, "sidebar-left.css")
	case "right":
		p.Stylesheets = append(p.Stylesheets, "sidebar-right.css")

	}

	// Copy shared stylesheets first
	for _, file := range App.Page.Theme.RootStylesheets {
		// If user has requested a dark theme, then don't copy skin.css
		// to the target. Copy theme-dark.css instead.
		// TODO: Decide whether this should be in root stylesheets and/or regular.
		file = App.getMode(file)
		/*
			if file == "theme-light.css" && App.FrontMatter.Mode == "dark" {
				file = "theme-dark.css"
			}
		*/
		// If it's a child theme, then get its stylesheets from the parent
		// directory.
		if App.FrontMatter.isChild {
			file = filepath.Join("..", file)
		}
		App.copyStylesheet(file)
	}
	for _, file := range p.Stylesheets {
		// Add the stylesheet tag
		// And copy the stylesheet itself
		file = App.getMode(file)
		/*
			if file == "theme-light.css" && App.FrontMatter.Mode == "dark" {
				file = "theme-dark.css"
			}
		*/
		// Create a matching directory for assets
		relDir := relDirFile(App.Site.path, App.Page.filePath)

		// Get full path of stylesheet to copy from theme directory.
		from := filepath.Join(App.Page.Theme.PageType.PathName, file)
		// Get directory to which this file will be copied for publishing
		assetDir := filepath.Join(App.Site.Publish, relDir, themeDir, App.FrontMatter.Theme, App.FrontMatter.PageType, App.Site.AssetDir)
		// Create a fully qualified filename for the published file
		to := filepath.Join(assetDir, file)
		var pathname string
		if !strings.HasPrefix(file, "http") {
			// Create the full pathname for the link tag, say, "themes/reference/assets/reset.css"
			pathname = filepath.Join(themeDir, App.FrontMatter.Theme, App.FrontMatter.PageType, App.Site.AssetDir, file)
		} else {
			// If it's hosted elsewhere the filename is sufficient.
			pathname = file
		}
		// Turn it into a "link" tag.
		App.appendStr(stylesheetTag(pathname))
		if !strings.HasPrefix(file, "http") {
			if err := Copy(from, to); err != nil {
				App.QuitError(errCode("0125", "from '"+from+"' to '"+to+"'"))
			}
		}
	}
	// Copy other files in the theme directory to the target publish directory.
	// This is whatever happens to be
	// in the theme directory with sizes.css, fonts.css, etc. Since those files
	// are stylesheets specified in the .TOML (or determined dynamically, like
	// sidebar-left.css and sidebar-right.css) it's easy. You generate a stylesheet
	// tag for them and then copy them right to the published theme directory.
	// The other files are dealt with here. Probably they would typically
	// be graphics files. They will be copied not to the
	// asset directory, but right into the document directory.
	// Which feels counterinutive
	// and kind of wrong, because they are likely to be something like social media
	// icons. More on this situation below, but of course they are actually
	// part of the page itself.

	// xxx

	for _, file := range App.Page.Theme.PageType.otherAssets {
		from := filepath.Join(App.Page.Theme.PageType.PathName, file)
		// Create a matching directory for assets
		relDir := relDirFile(App.Site.path, App.Page.filePath)
		// Create a fully qualified filename for the published file
		// which means depositing it in the document directoyr, not
		// the assets directory.
		// TODO: What we really want is to put the assets in the assets directory.
		// After all, they're in the theme directory (example: social media icon files),
		// and CSS files specified in the TOML are correctly sent to the assets directory.
		// But to do that we'd need some concept of an asset directory in the theme, so instead
		// of something like ![facebook icon](facebook-24x24-red.svg) in, for example.
		// nav.md, we'd need to do something like specify what files get copied in the
		// theme's TOML, or have some kind of ![facebook icon]({{themeDir}}/facebook-24x24-red.svg)
		// prefix.
		// TODO:
		assetDir := filepath.Join(App.Site.Publish, relDir)
		to := filepath.Join(assetDir, file)
		// xxx
		if err := Copy(from, to); err != nil {
			App.QuitError(errCode("0124", "from '"+from+"' to '"+to+"'"))
		}
	}
}

// publishThemeAssets() obtains a list of non-stylesheet asset files in the current
// PageType directory that should be published, so, anything but Markdown, toml, HTML, and a
// few other excluded types. It writes these to App.Page.Theme.PageType.otherAssets
// It writes stylesheets to App.Page.Theme.currPageType.stylesheets.
// That because the otherAssets files can just get copied over, by the stylesheets
// file list needs to be turned into stylesheet links.
func (App *App) publishThemeAssets() {
	// First get the list of stylesheets specified for this PageType.
	// Get a directory listing of all the non-source files
	dir := App.Page.Theme.PageType.PathName
	// Get the full directory listing.
	candidates, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	for _, file := range candidates {
		filename := file.Name()
		// If it's a file...
		if !file.IsDir() {
			if !hasExtensionFrom(filename, markdownExtensions) &&
				!hasExtensionFrom(filename, excludedAssetExtensions) {
				// If it's a stylesheet, add to the private list
				if hasExtension(filename, ".css") {
				} else {
					App.Page.Theme.PageType.otherAssets = append(App.Page.Theme.PageType.otherAssets, filename)
				}
			}
		} else {
			// Special case for :
			if filename == (THEME_HELP_SUBDIRNAME) {
				fmt.Println("Found special dir", filename)
			}
			// xxx
		}
	}
}

func (App *App) copyStylesheet(file string) {
	relDir := relDirFile(App.Site.path, App.Page.filePath)
	//assetDir := filepath.Join(App.Site.AssetDir, relDir, themeDir, App.FrontMatter.Theme, App.FrontMatter.PageType, App.Site.AssetDir)
	assetDir := filepath.Join(App.Site.AssetDir, relDir, themeDir, App.FrontMatter.Theme, App.FrontMatter.PageType, App.Site.AssetDir)
	from := filepath.Join(App.Page.Theme.PageType.PathName, file)
	to := filepath.Join(assetDir, file)
	pathname := filepath.Join(themeDir, App.FrontMatter.Theme, App.FrontMatter.PageType, App.Site.AssetDir, file)
	App.appendStr(stylesheetTag(pathname))
	// assetDir only exists if there was at least 1
	// markdown file in that directory. If it doesn't exist,
	// there's no reason to copy this file
	//to = filepath.Join(App.Site.Publish, relDir, themeDir, App.FrontMatter.Theme, App.FrontMatter.PageType, App.Site.AssetDir, file)
	to = filepath.Join(App.Site.Publish, relDir, themeDir, App.FrontMatter.Theme, App.FrontMatter.PageType, App.Site.AssetDir, file)
	if err := Copy(from, to); err != nil {
		App.QuitError(errCode("0916", "from '"+from+"' to '"+to+"'"))
	}
}

// Look alongside the current file to assets to publish
// for example, it's a news article and it has an image.
// TODO: This willb repeated for each file in the directory,
// so I need a way to do it only once.
func (App *App) publishPageAssets() {
	candidates, err := ioutil.ReadDir(App.Page.dir)
	if err != nil {
		return
	}

	// Check if this has already been done
	if App.Page.assets == nil {
		for _, file := range candidates {
			filename := file.Name()
			if !file.IsDir() {
				if !hasExtensionFrom(filename, markdownExtensions) &&
					!hasExtensionFrom(filename, excludedAssetExtensions) &&
					!hasExtension(filename, ".css") {
					App.Page.assets = append(App.Page.assets, filename)
				}
			}
		}
	}
}

// stylesheetTag() produces just that.
// Given the name of a stylesheet, like say "markdown.css",
// return it in a link tag.
func stylesheetTag(stylesheet string) string {
	// If no stylesheet specified just return empty string
	if stylesheet == "" {
		return ""
	}
	return `<link rel="stylesheet" href="` + stylesheet + `">` + "\n"
}

// startHTML() begins the HTML document and opens the head tag
func (App *App) startHTML() {
	App.Page.startHTML = ("<!DOCTYPE html>" + "\n" +
		"<html lang=" + App.Site.Language + ">" +
		`
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
`)
}

// firstHeader() returns the first header it founds in the markdown.
// It looks through the whole text for an h1. If not found,
// it looks for the first h2 to h6 it can find.
// Otherwise it returns ""
func firstHeader(markdown string) string {
	result := header1(markdown)
	if result != "" {
		return result
	}
	return header2To6(markdown)
}

// header1() extracts the first h1 it finds in the markdown
func header1(s string) string {
	any := h1.FindString(strings.Trim(s, "\n\t\r"))
	if any != "" {
		return (notPound.FindString(any))
	} else {
		return ""
	}
}

// header2To6() extracts the first h2-h2 it finds.
func header2To6(s string) string {
	any := anyHeader.FindString(strings.Trim(s, "\n\t\r"))
	if any != "" {
		return notPound.FindString(any)
	} else {
		return ""
	}
}

func (App *App) titleTag() {
	//App.appendStr("\n<title>" + title + "</title>\n")
	var title string
	if App.FrontMatter.Title != "" {
		title = App.FrontMatter.Title
	} else {
		title = firstHeader(string(App.Page.markdownStart))
	}
	if title == "" {
		title = productName + ": Title needed here, squib"
	}
	App.Page.titleTag = title

}

// Article() takes a document with optional front matter, parses
// out the front matter, and sends the Markdown portion to be converted.
// Write the HTML results to App.Page.Article
func (App *App) Article(filename string, input []byte) {
	// Extract front matter and parse.
	// Return the starting address of the Markdown.
	start, err := App.parseFrontMatter(filename, input)
	if err != nil {
		App.QuitError(errCode("0103", filename))
	}
	// Resolve any Go template variables before conversion to HTML.
	interp := App.interps(filename, string(start))
	App.Page.Article = App.markdownBufferToBytes([]byte(interp))
	var w WebPage
	w.html = App.Page.Article
	App.Site.WebPages[App.Page.filePath] = w

}

// stripHeading() returns the string following a Markdown heading.
// It is guaranteed a string of the form "### foo",
// where there can be 1-6 # characters. As with findFirstHeading()
// it's not the most general purpose routine ever and may
// need revisiting, because it assumes a lot about the
// Markdown format, which is more flexible than this.
// TODO: Check to see if it's even used
func stripHeading(heading string) string {
	match := strings.Index(heading, " ")
	l := len(heading)
	if l < 0 {
		return ""
	}
	return (heading[match+1 : l])
}

// headTags() inserts miscellaneous items such as Google Analytics tags
// into the header before it's close.
func (App *App) headTags() {
	App.Page.headTags = App.headerFiles() +
		App.headerTagGanalytics()
}

// headerFiles() finds all the files in the headers subdirectory
// and copies them into the HMTL headers of every file on the site.
func (App *App) headerFiles() string {
	var h string
	headers, err := ioutil.ReadDir(App.Site.headersPath)
	if err != nil {
		App.QuitError(errCode("0706", App.Site.headersPath))
	}
	for _, file := range headers {
		h += fileToString(filepath.Join(App.Site.headersPath, file.Name()))
	}
	return h
}

// headerTagGanalytics() generates a Google Analytics script, if a tracking
// ID is available. If not it returns an empty string so it's always
// safe to call.
func (App *App) headerTagGanalytics() string {
	if App.Site.Ganalytics == "" {
		return ""
	}
	result := strings.Replace(ganalyticsTag, "XX-XXXXXXXXX-X", App.Site.Ganalytics, 1)
	// Only include if it worked
	if result == ganalyticsTag {
		return ""
	}
	return result + "\n"
}

// getDescription() does everything it can to generate a Description
// metatag for the file.
func (App *App) descriptionTag() {
	// Best case: user supplied the description in the front matter.
	if App.FrontMatter.Description != "" {
		App.Page.descriptionTag = App.FrontMatter.Description
	} else if App.Site.Branding != "" {
		App.Page.descriptionTag = App.Site.Branding
	} else if App.Site.Name != "" {
		App.Page.descriptionTag = App.Site.Name
	} else {
		App.Page.descriptionTag = "Powered by " + productName
	}

}

// localFiles() copies any files that happen to be lying around.
// It also generates stylesheet links
func (App *App) localFiles(relDir string) {
	// Copy any associated assets such as
	// images in the same directory.
	dirHasMarkdownFiles := App.publishLocalFiles(App.Page.dir)
	if dirHasMarkdownFiles {
		// Create its theme directory
		assetDir := filepath.Join(App.Site.Publish, relDir, themeDir, App.FrontMatter.Theme, App.FrontMatter.PageType, App.Site.AssetDir)
		if err := os.MkdirAll(assetDir, PUBLIC_FILE_PERMISSIONS); err != nil {
			App.QuitError(errCode("0402", assetDir))
		}
		App.publishPageTypeAssets()
	}
}

// closeHeadOpenBody() writes the closing </head> tag
// and starts the <body> tag
func (App *App) closeHeadOpenBody() {
	var closer = `
</head>
<body>
`
	App.appendStr(closer)
}

func wrapTag(tag string, contents string, block bool) string {
	var newline string
	if block {
		newline = "\n"
	}
	if len(tag) > 3 {
		output := newline + tag + contents + tag[:1] + "/" + tag[1:] + newline
		return output
	}
	return ""
}

// Wraps the contents within a block/style tag,
// so it turns <p>hello, world.<p> into
// <article><p>hello, world.<p></article>
// If block is true, adds newlines strictly for
// clarity in the output HTML.
func wrapTagBytes(tag string, html []byte, block bool) string {
	var newline string
	if block {
		newline = "\n"
	}
	if len(tag) > 3 {
		output := newline + tag + string(html) + tag[:1] + "/" + tag[1:] + newline
		return output
	}
	return ""
}

// pageRegionToHTML() takes an page region (header, nav, article, sidebar, or footer)
// and converts it to HTML. All we know is that it's been specified
// but we don't know whether's a Markdown file, inline HTML, whatever.
func (App *App) pageRegionToHTML(a *pageRegion, tag string) string {
	switch tag {
	case "<header>", "<nav>", "<article>", "<aside>", "<footer>":
		var path string
		path = filepath.Join(App.Page.Theme.PageType.PathName, a.File)

		// A .sidebar file trumps all else.
		// See if there's a file with the same name as
		// the root source file but with a .sidebar extension.
		if tag == "<aside>" {
			// Base it on the root Markdown filename and the
			// extension .sidebar, so foo.md might also have
			// a foo.sidebar.
			// Construct a path to possible .sidebar file.
			sidebarfile := replaceExtension(App.Page.filePath, "sidebar")
			if fileExists(sidebarfile) {
				// If that .sidebar file exists, immediately
				// insert into the stream and leave,
				// because it's the highest priority.
				input := fileToBuf(sidebarfile)
				return wrapTag(tag, string(App.MdFileToHTMLBuffer(sidebarfile, input)), true)
			}
		}
		// Inline HTML is the highest priority
		if a.HTML != "" {
			return a.HTML
		}
		// Skip if there's no file specified
		if a.File == "" {
			return ""
		}
		var input []byte
		// Error if the specified file can't be found.
		if !fileExists(path) {
			App.QuitError(errCode("1015", path))
		}
		if isMarkdownFile(path) {
			input = fileToBuf(path)
			return wrapTag(tag, string(App.MdFileToHTMLBuffer(path, input)), true)
		}
		return fileToString(path)
	default:
		App.QuitError(errCode("1203", tag))
	}
	return ""
}

// Generates a meta tag
func metatag(tag string, content string) string {
	const quote = `"`
	return ("\n<meta name=" + quote + tag + quote + " content=" + quote + content + quote + ">\n")
}
