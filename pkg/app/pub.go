package app

import (
	"bytes"
	"fmt"
	"github.com/tomcam/mb/pkg/defaults"
	"github.com/tomcam/mb/pkg/errs"
	"github.com/tomcam/mb/pkg/mdext"
	"github.com/tomcam/mb/pkg/slices"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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
func (a *App) publishFile(filename string) error {
	var input []byte
	var err error
	// Get a fresh new Page object if doing more
	// than one file at a clip. Which is obviously
	// most of the time.
	var p Page
	a.Page = &p
	var f FrontMatter
	a.FrontMatter = &f
	a.Page.filePath = filename
	a.Verbose(filename)
	a.Page.filename = filepath.Base(filename)
	a.Page.dir = currDir()
	a.Verbose("%s", filename)
	// Read the whole Markdown file into memory as a byte slice.
	input, err = ioutil.ReadFile(filename)
	if err != nil {
		return errs.ErrCode("0102", filename)
	}
	// Extract front matter and parse.
	// Obviously that includes an optional theme or pagetype designation.
	// Starting at the Markdown, convert to HTML.
	// Interpret templates as well.
	a.Page.markdownStart, err = a.parseFrontMatter(filename, input)
	if err != nil {
		return errs.ErrCode("0103", filename)
	}

	// If no theme was specified in the front matter, but one was specified in the
	// site config, make the one specified in site.toml the theme.
	if a.Site.Theme != "" && a.FrontMatter.Theme == "" {
		a.FrontMatter.Theme = a.Site.Theme
	}
	// If no theme was specified at all, use the Metabuzz default.
	if a.FrontMatter.Theme == "" {
		a.FrontMatter.Theme = defaults.DefaultThemeName
	}
	a.loadTheme()
	// Parse front matter.
	// Convert article to HTML
	a.Article(filename, input)
	// xxx
	// Begin HTML document.
	// Open the <head> tag.
	a.startHTML()

	// If a title wasn't specified in the front matter,
	// put up a self-aggrandizing error message.
	a.titleTag()

	a.descriptionTag()

	a.headTags()

	// Output filename
	outfile := replaceExtension(filename, "html")
	relDir := relDirFile(a.Site.path, outfile)

	// START ASSEMBLING PAGE
	a.appendStr(a.Page.startHTML)
	a.appendStr(wrapTag("<title>", a.Page.titleTag, true))
	a.appendStr(metatag("description", a.Page.descriptionTag))
	a.appendStr(a.Page.headTags)

	// Hoover up any miscellanous files lying around,
	// like other HTML files, graphic assets, etc.
	a.localFiles(relDir)

	// Write the closing head tag and the opening
	// body tag
	a.closeHeadOpenBody()

	//a.appendStr(wrapTag("<article>", []byte(a.Page.Article), true))
	a.appendStr(a.pageRegionToHTML(&a.Page.Theme.PageType.Header, "<header>"))
	a.appendStr(a.pageRegionToHTML(&a.Page.Theme.PageType.Nav, "<nav>"))
	a.appendStr(a.pageRegionToHTML(&a.Page.Theme.PageType.Article, "<article>"))
	sidebar := strings.ToLower(a.FrontMatter.Sidebar)
	if sidebar == "left" || sidebar == "right" {
		a.appendStr(a.pageRegionToHTML(&a.Page.Theme.PageType.Sidebar, "<aside>"))
	}
	a.appendStr(a.pageRegionToHTML(&a.Page.Theme.PageType.Footer, "<footer>"))

	// Complete the HTML document with closing <body> and <html> tags
	a.appendStr(closingHTMLTags)

	// Strip out everything but the filename.
	base := filepath.Base(outfile)

	// Write everything to a temp file so in case there was an error, the
	// previous HTML file is preserved.
	tmpFile, err := ioutil.TempFile(a.Site.Publish, defaults.ProductName+"-tmp-")
	if err != nil {
		a.QuitError(errs.ErrCode("PREVIOUS", err.Error()))
	}

	if err = writeTextFile(tmpFile.Name(), string(a.Page.html)); err != nil {
		a.QuitError(errs.ErrCode("PREVIOUS", err.Error()))
	}
	// Ensure the file gets closed before exiting
	defer os.Remove(tmpFile.Name())
	// Get the relative directory.
	// relDir = relDirFile(a.Site.path, outfile)
	a.Page.Path = relDir
	// If there's a README.md and no index.md, rename
	// the output file to index.html
	if a.Page.filename == "README.md" && !a.Site.dirs[a.Page.dir].mdOptions.IsOptionSet(HasIndexMd) {
		base = "index.html"
	}

	// Generate the full pathname of the matching output file, as it will
	// appear in its published location.
	outfile = filepath.Join(a.Site.Publish, relDir, base)
	// If the write succeeded, rename it to the output file
	// This way if there was an existing HTML file but there was
	// an error in output this time, it doesn't get clobbered.
	if err = os.Rename(tmpFile.Name(), outfile); err != nil {
		a.QuitError(errs.ErrCode("0212", outfile))
	}

	if !fileExists(outfile) {
		a.QuitError(errs.ErrCode("0910", outfile))
	}
	a.Verbose("\tCreated file %s\n", outfile)
	a.fileCount++
	//
	// Success
	return nil
}

// Takes a buffer containing Markdown
// and converts to HTML.
// Doesn't know about front matter.
func (a *App) markdownBufferToBytes(input []byte) []byte {
	buf := new(bytes.Buffer)
	node := a.markdownAST(input)
	if err := a.newGoldmark().Renderer().Render(buf, input, node); err != nil {
		// TODO: Need something like displayErrCode("1010") or whatever
		a.QuitError(errs.ErrCode("0920", err.Error()))
		return nil
	}
	return buf.Bytes()
}

// markdownAST returns the goldmark AST for the input.
func (a *App) markdownAST(input []byte) ast.Node {
	ctx := parser.NewContext()
	p := a.newGoldmark().Parser()
	return p.Parse(text.NewReader(input), parser.WithContext(ctx))
}

// newGoldmark returns the a goldmark object with a parser and renderer.
func (a *App) newGoldmark() goldmark.Markdown {
	exts := []goldmark.Extender{
		extension.GFM,
		extension.DefinitionList,
		extension.Footnote,
		highlighting.NewHighlighting(
			highlighting.WithStyle(a.Site.MarkdownOptions.HighlightStyle),
			highlighting.WithFormatOptions()),
	}

	parserOpts := []parser.Option{parser.WithAttribute()}
	if a.Site.MarkdownOptions.HeadingIDs {
		parserOpts = append(parserOpts, parser.WithAutoHeadingID())
	}

	renderOpts := []renderer.Option{
		html.WithUnsafe(),
		html.WithXHTML(),
	}
	if a.Site.MarkdownOptions.hardWraps {
		renderOpts = append(renderOpts, html.WithHardWraps())
	}

	return goldmark.New(
		goldmark.WithExtensions(exts...),
		goldmark.WithParserOptions(parserOpts...),
		goldmark.WithRendererOptions(renderOpts...),
	)
}

// appendBytes() Appends a byte slice to the byte slice containing the rendered output
func (a *App) appendBytes(b []byte) {
	a.Page.html = append(a.Page.html, b...)
}

// appendStr() Appends a string to the byte slice containing the rendered output
func (a *App) appendStr(s string) {
	a.Page.html = append(a.Page.html, s...)
}

// MdFileToHTMLBuffer() takes a byte slice buffer containing
// a pure Markdown file as input, and returns
// a byte slice containing the file converted to HTML.
// It doesn't know about front matter.
// So it should be preceded by a call to App.parseFrontMatter()
// if there's any possibility that the file contains front matter.
// In the spirit of a browser, it simply returns an empty buffer on error.
func (a *App) MdFileToHTMLBuffer(filename string, input []byte) []byte {
	// Resolve any Go template variables before conversion to HTML.
	interp := a.interps(filename, string(input))
	// Convert markdown to HTML.
	return a.markdownBufferToBytes([]byte(interp))
}

func (a *App) addMdOption(dir string, mdOption MdOptions) {
	d := a.Site.dirs[dir]
	d.mdOptions |= mdOption
	a.Site.dirs[dir] = d
}

func (a *App) setMdOption(dir string, mdOption MdOptions) {
	d := a.Site.dirs[dir]
	d.mdOptions = mdOption
	a.Site.dirs[dir] = d
}

// publishLocalFiles() get called for every markdown file
// in the directory. It copies assets like image files & so forth
// from the source file's current directory to the publish location,
// creating a new subdirectory as needed.
// For example, if your article references ![cat](cat.png)
// then presumably cat.png is in the current directory.
// This copies all non-excluded files, such as cat.png and
// any other assets, from this directory
// into its matching publish directory,
// same as the source markdown file.
// Creates a subdirectory under Publish if in a subdirectory
// and one hasn't yet been created.
// Keeps track of which directories have had their assets copied to
// avoid redundant copies.
// Returns true if there are any markdown files in the current directory.
// Returns false if markdown files (or any files at all) are absent.
func (a *App) publishLocalFiles(dir string) bool {
	relDir := relativeDirectory(a.Site.path, dir)
	pubDir := filepath.Join(a.Site.Publish, relDir)

	// If this directory hasn't been created, create it.
	if !a.Site.dirs[pubDir].mdOptions.IsOptionSet(MarkdownDir) {
		if err := os.MkdirAll(pubDir, defaults.PublicFilePermissions); err != nil {
			a.QuitError(errs.ErrCode("0404", pubDir, err.Error()))
		}
		// Mark that the directory has been created so this
		// doesn't get repeated.
		//a.Site.dirs[dir].MdOptions |= MarkdownDir
		a.addMdOption(dir, MarkdownDir)
	}
	// Get the directory listing.
	candidates, err := ioutil.ReadDir(dir)
	if err != nil {
		a.QuitError(errs.ErrCode("1016", dir, err.Error()))
	}

	// Get list of files in the local directory to exclude from copy
	excludeFromDir := slices.NewSearchInfo(a.FrontMatter.ExcludeFilenames)

	// First check the directory to ensure there's at least 1 markdown file.
	hasMarkdown := false

	// Look for the specific file README.md, which competes with
	// index.md:
	// https://stackoverflow.com/questions/58826517/why-do-some-static-site-generators-use-readme-md-instead-of-index-md
	for _, file := range candidates {
		filename := file.Name()
		if hasExtensionFrom(filename, defaults.MarkdownExtensions) {
			hasMarkdown = true
		}

		if filename == "README.md" {
			//a.Site.dirs[dir].MdOptions |= HasReadmeMd
			a.addMdOption(dir, HasReadmeMd)

		}
		if strings.ToLower(filename) == "index.md" {
			//a.Site.dirs[dir].MdOptions |= HasIndexMd
			a.addMdOption(dir, HasIndexMd)
		}

	}

	if hasMarkdown {
		// Flag this as a directory that contains at least
		// 1 markdown file.
		a.addMdOption(dir, MarkdownDir)
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
				!hasExtensionFrom(filename, defaults.MarkdownExtensions) &&
				!excludeFromDir.Contains(filename) &&
				!strings.HasPrefix(filename, ".") {
				// It's a markdown file.
				// Got the file. Get its fully qualified name.
				copyFrom := filepath.Join(dir, filename)
				// Figure out the target directory.
				relDir := relDirFile(a.Site.path, copyFrom)
				// Get the target file's fully qualified filename.
				copyTo := filepath.Join(a.Site.Publish, relDir, filename)
				//a.Verbose("\tCopying '%s' to '%s'\n", copyFrom, copyTo)
				if err := Copy(copyFrom, copyTo); err != nil {
					//a.QuitError(err.Error())
					a.QuitError(errs.ErrCode("PREVIOUS", err.Error()))
				}
				// TODO: Get rid of fileExists() when possible
				if !fileExists(copyTo) {
					a.QuitError(errs.ErrCode("PREVIOUS", copyTo))
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
func (a *App) publishPageTypeAssets() {
	// Is the default aka parent theme?
	p := a.Page.Theme.PageType
	if p.name == "" || p.name == defaults.DefaultThemeName {
		// Default PageType
		a.publishAssets()
	} else {
		// It's a child theme aka pagetype.
		a.Page.Theme.PageType = p
		a.publishAssets()
		// TODO:This would happen more than once
		// with multiple pageTypes, so I should
		// eliminate that.
	}
}

// getMode() checks if the stylesheet is dark or light and adjusts as needed
func (a *App) getMode(stylesheet string) string {
	if stylesheet == "theme-light.css" && a.FrontMatter.Mode == "dark" {
		stylesheet = "theme-dark.css"
	}
	return stylesheet
}

// publishAssets() copies out the stylesheets, graphics, and other
// relevant files from the pageType (or default theme) directory
// to be published.
func (a *App) publishAssets() {
	p := a.Page.Theme.PageType
	a.publishPageAssets()
	a.publishThemeAssets()
	// Copy out different stylesheet depending on the
	// type of sidebar, if any.
	switch strings.ToLower(a.FrontMatter.Sidebar) {
	case "left":
		p.Stylesheets = append(p.Stylesheets, "sidebar-left.css")
	case "right":
		p.Stylesheets = append(p.Stylesheets, "sidebar-right.css")

	}
	a.copyStyleSheets(p)
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

	for _, file := range a.Page.Theme.PageType.otherAssets {
		from := filepath.Join(a.Page.Theme.PageType.PathName, file)
		// Create a matching directory for assets
		relDir := relDirFile(a.Site.path, a.Page.filePath)
		// Create a fully qualified filename for the published file
		// which means depositing it in the document directoyr, not
		// the assets directory.
		// TODO: What we really want is to put the assets in the assets directory.
		// After all, they're in the theme directory (example: social media icon files),
		// and CSS files specified in the TOML are correctly sent to the assets directory.
		// But to do that we'd need some concept of an asset directory in the theme, so instead
		// of something like ![facebook icon](facebook-24x24-red.svg) in, for example.
		// nav.md, we'd need to do something like specify what files get copied in the
		// theme's TOML, or have some kind of ![facebook icon]({{ThemeDir}}/facebook-24x24-red.svg)
		// prefix.
		// TODO:
		assetDir := filepath.Join(a.Site.Publish, relDir)
		to := filepath.Join(assetDir, file)
		if err := Copy(from, to); err != nil {
			a.QuitError(errs.ErrCode("0124", "from '"+from+"' to '"+to+"'"))
		}
	}
}

func (a *App) copyStyleSheets(p PageType) {
	// Copy shared stylesheets first
	for _, file := range a.Page.Theme.RootStylesheets {
		// If user has requested a dark theme, then don't copy skin.css
		// to the target. Copy theme-dark.css instead.
		// TODO: Decide whether this should be in root stylesheets and/or regular.
		file = a.getMode(file)
		// If it's a child theme, then copy its stylesheets from the parent
		// directory.
		if a.FrontMatter.isChild {
			file = filepath.Join("..", file)
		}
		a.copyStylesheet(file)
	}
	for _, file := range p.Stylesheets {
		file = a.getMode(file)
		a.copyStylesheet(file)
	}
}

// publishThemeAssets() obtains a list of non-stylesheet asset files in the current
// PageType directory that should be published, so, anything but Markdown, toml, HTML, and a
// few other excluded types. It writes these to App.Page.Theme.PageType.otherAssets
// It writes stylesheets to App.Page.Theme.currPageType.stylesheets.
// That because the otherAssets files can just get copied over, by the stylesheets
// file list needs to be turned into stylesheet links.
func (a *App) publishThemeAssets() {
	// First get the list of stylesheets specified for this PageType.
	// Get a directory listing of all the non-source files
	dir := a.Page.Theme.PageType.PathName
	// Get the full directory listing.
	candidates, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	for _, file := range candidates {
		filename := file.Name()
		// If it's a file...
		if !file.IsDir() {
			if !hasExtensionFrom(filename, defaults.MarkdownExtensions) &&
				!hasExtensionFrom(filename, defaults.ExcludedAssetExtensions) {
				// If it's a stylesheet, add to the private list
				if hasExtension(filename, ".css") {
				} else {
					a.Page.Theme.PageType.otherAssets = append(a.Page.Theme.PageType.otherAssets, filename)
				}
			}
		} else {
			// Special case for :
			if filename == (defaults.ThemeHelpSubdirname) {
				fmt.Println("Found special dir", filename)
			}
		}
	}
}

func (a *App) copyStylesheet(file string) {
	if strings.HasPrefix(strings.ToLower(file), "http") {
		a.appendStr(stylesheetTag(file))
		return
	}
	relDir := relDirFile(a.Site.path, a.Page.filePath)
	assetDir := filepath.Join(a.Site.AssetDir, relDir, defaults.ThemeDir, a.FrontMatter.Theme, a.FrontMatter.PageType, a.Site.AssetDir)
	from := filepath.Join(a.Page.Theme.PageType.PathName, file)
	to := filepath.Join(assetDir, file)
	var pathname string
	if strings.HasPrefix(strings.ToLower(file), "http") {
		pathname = file
		fmt.Println(pathname)
	} else {
		pathname = filepath.Join(defaults.ThemeDir, a.FrontMatter.Theme, a.FrontMatter.PageType, a.Site.AssetDir, file)
	}
	a.appendStr(stylesheetTag(pathname))
	to = filepath.Join(a.Site.Publish, relDir, defaults.ThemeDir, a.FrontMatter.Theme, a.FrontMatter.PageType, a.Site.AssetDir, file)
	if err := Copy(from, to); err != nil {
		a.QuitError(errs.ErrCode("0916", "from '"+from+"' to '"+to+"'"))
	}
}

// Look alongside the current file to assets to publish
// for example, it's a news article and it has an image.
// TODO: This willb repeated for each file in the directory,
// so I need a way to do it only once.
func (a *App) publishPageAssets() {
	candidates, err := ioutil.ReadDir(a.Page.dir)
	if err != nil {
		return
	}

	// Check if this has already been done
	if a.Page.assets == nil {
		for _, file := range candidates {
			filename := file.Name()
			if !file.IsDir() {
				if !hasExtensionFrom(filename, defaults.MarkdownExtensions) &&
					!hasExtensionFrom(filename, defaults.ExcludedAssetExtensions) &&
					!hasExtension(filename, ".css") {
					a.Page.assets = append(a.Page.assets, filename)
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
func (a *App) startHTML() {
	a.Page.startHTML = "<!DOCTYPE html>" + "\n" +
		"<html lang=" + a.Site.Language + ">" +
		`
	<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width,initial-scale=1">
	`
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

func (a *App) titleTag() {
	title := defaults.ProductName + ": Title needed here, squib"
	switch {
	case a.FrontMatter.Title != "":
		title = a.FrontMatter.Title
	default:
		node := a.markdownAST(a.Page.markdownStart)
		t := mdext.InferTitle(node, a.Page.markdownStart)
		if t != "" {
			title = t
		}
	}
	a.Page.titleTag = title
}

// Article() takes a document with optional front matter, parses
// out the front matter, and sends the Markdown portion to be converted.
// Write the HTML results to App.Page.Article
func (a *App) Article(filename string, input []byte) {
	// Extract front matter and parse.
	// Return the starting address of the Markdown.
	start, err := a.parseFrontMatter(filename, input)
	if err != nil {
		a.QuitError(errs.ErrCode("0103", filename))
	}
	// Resolve any Go template variables before conversion to HTML.
	interp := a.interps(filename, string(start))
	a.Page.Article = a.markdownBufferToBytes([]byte(interp))
	var w WebPage
	w.html = a.Page.Article
	a.Site.WebPages[a.Page.filePath] = w

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
func (a *App) headTags() {
	a.Page.headTags = a.headerFiles() +
		a.headTagGanalytics()
}

// headerFiles() finds all the files in the headers subdirectory
// and copies them into the HMTL headers of every file on the site.
func (a *App) headerFiles() string {
	var h string
	headers, err := ioutil.ReadDir(a.Site.headTagsPath)
	if err != nil {
		a.QuitError(errs.ErrCode("0706", a.Site.headTagsPath))
	}
	for _, file := range headers {
		h += fileToString(filepath.Join(a.Site.headTagsPath, file.Name()))
	}
	return h
}

// headTagGanalytics() generates a Google Analytics script, if a tracking
// ID is available. If not it returns an empty string so it's always
// safe to call.
func (a *App) headTagGanalytics() string {
	if a.Site.Ganalytics == "" {
		return ""
	}
	result := strings.Replace(ganalyticsTag, "XX-XXXXXXXXX-X", a.Site.Ganalytics, 1)
	// Only include if it worked
	if result == ganalyticsTag {
		return ""
	}
	return result + "\n"
}

// getDescription() does everything it can to generate a Description
// metatag for the file.
func (a *App) descriptionTag() {
	// Best case: user supplied the description in the front matter.
	if a.FrontMatter.Description != "" {
		a.Page.descriptionTag = a.FrontMatter.Description
	} else if a.Site.Branding != "" {
		a.Page.descriptionTag = a.Site.Branding
	} else if a.Site.Name != "" {
		a.Page.descriptionTag = a.Site.Name
	} else {
		a.Page.descriptionTag = "Powered by " + defaults.ProductName
	}

}

// localFiles() copies any files that happen to be lying around.
// It also generates stylesheet links
func (a *App) localFiles(relDir string) {
	// Copy any associated assets such as
	// images in the same directory.
	dirHasMarkdownFiles := a.publishLocalFiles(a.Page.dir)
	if dirHasMarkdownFiles {
		// Create its theme directory
		assetDir := filepath.Join(
			a.Site.Publish, relDir, defaults.ThemeDir, a.FrontMatter.Theme, a.FrontMatter.PageType, a.Site.AssetDir)
		if err := os.MkdirAll(assetDir, defaults.PublicFilePermissions); err != nil {
			a.QuitError(errs.ErrCode("0402", assetDir))
		}
		a.publishPageTypeAssets()
	}
}

// closeHeadOpenBody() writes the closing </head> tag
// and starts the <body> tag
func (a *App) closeHeadOpenBody() {
	var closer = `
</head>
<body>
`
	a.appendStr(closer)
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
func (a *App) pageRegionToHTML(pr *pageRegion, tag string) string {
	switch tag {
	case "<header>", "<nav>", "<article>", "<aside>", "<footer>":
		var path string
		path = filepath.Join(a.Page.Theme.PageType.PathName, pr.File)

		// A .sidebar file trumps all else.
		// See if there's a file with the same name as
		// the root source file but with a .sidebar extension.
		if tag == "<aside>" {
			// Base it on the root Markdown filename and the
			// extension .sidebar, so foo.md might also have
			// a foo.sidebar.
			// Construct a path to possible .sidebar file.
			sidebarfile := replaceExtension(a.Page.filePath, "sidebar")
			if fileExists(sidebarfile) {
				// If that .sidebar file exists, immediately
				// insert into the stream and leave,
				// because it's the highest priority.
				input := fileToBuf(sidebarfile)
				return wrapTag(tag, string(a.MdFileToHTMLBuffer(sidebarfile, input)), true)
			}
		}

		// Exception: a theme without an article pagetype specified is equivalent
		// to <article>{{ article }}</article>. So wrap the entire article in the
		// appropriate tag.
		if tag == "<article>" {
			if a.Page.Theme.PageType.Article.File == "" && a.Page.Theme.PageType.Article.HTML == "" {
				return wrapTag(tag, string(a.Page.Article), true)
			}
		}

		// Inline HTML is the highest priority
		if pr.HTML != "" {
			return pr.HTML
		}
		// Skip if there's no file specified
		if pr.File == "" {
			return ""
		}
		var input []byte
		// Error if the specified file can't be found.
		if !fileExists(path) {
			a.QuitError(errs.ErrCode("1015", path))
		}
		if isMarkdownFile(path) {
			input = fileToBuf(path)
			if tag == "<article>" {
				return string(a.MdFileToHTMLBuffer(path, input))
			} else {
				return wrapTag(tag, string(a.MdFileToHTMLBuffer(path, input)), true)
			}
		}
		return fileToString(path)
	default:
		a.QuitError(errs.ErrCode("1203", tag))
	}
	return ""
}

// metatag generates a meta tag. It's complicated.
func metatag(tag string, content string) string {
	const quote = `"`
	return ("\n<meta name=" + quote + tag + quote + " content=" + quote + content + quote + ">\n")
}
