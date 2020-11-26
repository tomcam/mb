package app

import (
  //"fmt"
	"bytes"
	"encoding/json"
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
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
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

  /* True means load parent theme */
  a.loadTheme(true)
  /* False means load pagetype, if any (aka child theme */
  a.loadTheme(false)
	// Convert article to HTML
	a.Article(filename, input)
	// Begin HTML document.
	// Open the <head> tag.
	a.startHTML()

	// Output filename
	outfile := replaceExtension(filename, "html")
	relDir := relDirFile(a.Site.path, outfile)
	// Strip out everything but the filename.
	base := filepath.Base(outfile)

	title := a.titleTag()

	// Extract the title and body from this page.
	// Convert to JSON, and add as a record to the
	// search index file.
	node := a.markdownAST(a.Page.markdownStart)
	docPath := "/" + path.Join(relDir, strings.TrimSuffix(base, ".html"))
	if strings.HasSuffix(docPath, "/index") {
		docPath = strings.TrimSuffix(docPath, "index")
	}
	doc := mdext.Doc{
		Path:  docPath,
		Title: title,
		Body:  mdext.BuildDocBody(node, a.Page.markdownStart),
	}
	// xxx
	if err := appendIndexDoc(a.Site.SearchJSONFilePath, doc); err != nil {
		a.QuitError(errs.ErrCode("PREVIOUS", err.Error()))
	}

	a.descriptionTag()

	a.headTags()

	// START ASSEMBLING PAGE
	a.appendStr(a.Page.startHTML)
	a.appendStr(wrapTag("<title>", title, true))
	a.appendStr(metatag("description", a.Page.descriptionTag))
	a.appendStr(a.Page.headTags)

	// Hoover up any miscellaneous files lying around,
	// like other HTML files, graphic assets, etc.
	a.localFiles(relDir)

	// Write the closing head tag and the opening
	// body tag
	a.closeHeadOpenBody()

	a.appendStr(a.layoutElementToHTML(&a.Page.Theme.PageType.Header, "<header>"))
	a.appendStr(a.layoutElementToHTML(&a.Page.Theme.PageType.Nav, "<nav>"))
	a.appendStr(a.layoutElementToHTML(&a.Page.Theme.PageType.Article, "<article>"))
	sidebar := strings.ToLower(a.FrontMatter.Sidebar)
	if sidebar == "left" || sidebar == "right" {
		a.appendStr(a.layoutElementToHTML(&a.Page.Theme.PageType.Sidebar, "<aside>"))
	} else {
		// TODO: If you have sidebar="ight" for example and/or the word "sidebar"
		// appears in the main Markdown, you get a bonus error message
		///a.QuitError(errs.ErrCode("1019", filename))
	}
	a.appendStr(a.layoutElementToHTML(&a.Page.Theme.PageType.Footer, "<footer>"))

	// Complete the HTML document with closing <body> and <html> tags
	a.appendStr(closingHTMLTags)

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

// appendIndexDoc appends doc as newline-delimited JSON to file.
// Called for each source .md file
func appendIndexDoc(file string, doc mdext.Doc) error {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	b, err := json.Marshal(doc)
	b = append(b, '\n') // use newline delimited JSON; 1 line per file
	if err != nil {
		return err
	}
	n, err := f.Write(b)
	if err != nil {
		return err
	}
	if n < len(b) {
		return io.ErrShortWrite
	}
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

	for _, file := range candidates {
		filename := file.Name()
		if hasExtensionFrom(filename, defaults.MarkdownExtensions) {
			hasMarkdown = true
		}

    // Look for the specific file README.md, which competes with
    // index.md:
    // https://stackoverflow.com/questions/58826517/why-do-some-static-site-generators-use-readme-md-instead-of-index-md
		if filename == "README.md" {
			a.addMdOption(dir, HasReadmeMd)

		}
		if strings.ToLower(filename) == "index.md" {
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

func (a *App) publishAssets() {
	p := a.Page.Theme.PageType
	a.findPageAssets()
	a.findThemeAssets()
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
    // xxxx
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

/*
// CSSDir() computes the fully qualified directory
// name for .CSS files, based on App.assetPath()
// (it is assumed to be a subdirectory of assetPath())
func (a *App) CSSDir() string {
	return filepath.Join(a.assetPath(), a.Site.CSSDir)
}
*/

// ImageDir() computes the fully qualified directory
// name for image files, based on App.assetPath()
// (it is assumed to be a subdirectory of assetPath())
func (a *App) ImageDir() string {
	return filepath.Join(a.assetPath(), a.Site.ImageDir)
}

// assetPath() computes the fully qualified directory
// name for assets, based on Site.AssetDir, etc.
func (a *App) assetPath() string {
	return filepath.Join(a.Site.path, a.Site.Publish, a.Site.AssetDir, "themes", a.FrontMatter.Theme)
}

// relTargetThemeDir() computes the relative destination directory
// name for theme assets
func (a *App) relTargetThemeDir() string {
  // xxxx Look up the html <base> tag
	return filepath.Join("/", a.Site.AssetDir, "themes", a.FrontMatter.Theme, a.FrontMatter.PageType)
}

// fullTargetThemeDir() computes the fully qualified destination directory
// name for theme assets
func (a *App) fullTargetThemeDir() string {
	return filepath.Join(a.Site.Publish, a.relTargetThemeDir())
}

// copyStyleSheet() takes the name of a stylesheet specified
// in the theme and copies it to the destination (publish)
// directory.
func (a *App) copyStyleSheet(file string) {
	// Pass through if not a local file
	if strings.HasPrefix(strings.ToLower(file), "http") {
		a.appendStr(stylesheetTag(file))
		return
	}

  var from string
	// Get fully qualified source filename to copy.
	/// xxxfrom := filepath.Join(a.parentThemeFullDirectory(), file)
  if a.FrontMatter.isChild {
	  from = filepath.Join(a.childThemeFullDirectory(), file)
  } else {
	  from = filepath.Join(a.parentThemeFullDirectory(), file)
  }

	// Relative path to the publish directory for themes
	pathname := filepath.Join(a.relTargetThemeDir(), file)
	// Write out the link
	a.appendStr(stylesheetTag(pathname))

	to := filepath.Join(a.fullTargetThemeDir(), file)
	if from == to {
		a.QuitError(errs.ErrCode("0922", "from '"+from+"' to '"+to+"'"))
	}

	// Actually copy the style sheet to its destination
	if err := Copy(from, to); err != nil {
		a.QuitError(errs.ErrCode("0916", "from '"+from+"' to '"+to+"'"))
	}
}

func (a *App) copyRootStylesheets() {
	for _, file := range a.Page.Theme.PageType.RootStylesheets {
		// If user has requested a dark theme, then don't copy skin.css
		// to the target. Copy theme-dark.css instead.
		// TODO: Decide whether this should be in root stylesheets and/or regular.
		file = a.getMode(file)
		// If it's a child theme, then copy its stylesheets from the parent
		// directory.
		if a.FrontMatter.isChild {
			file = filepath.Join("..", file)
		}
		a.copyStyleSheet(file)
	}
}

// copyStyleSheets() takes the list of style sheets
// given in the theme TOML file and copies them to the
// publishing (asset) directory
func (a *App) copyStyleSheets(p PageType) {
	dir := a.fullTargetThemeDir()
  /*
	if dirExists(dir) {
		fmt.Println("Directory " + dir + " already exists. Sayanara. xxx")
		return
	}
  */
	if err := os.MkdirAll(dir, defaults.PublicFilePermissions); err != nil {
		a.QuitError(errs.ErrCode("0402", dir))
	}

	a.copyRootStylesheets()

	for _, file := range p.Stylesheets {
		file = a.getMode(file)
		a.copyStyleSheet(file)
	}
	// responsive.css is always last
	a.copyStyleSheet("responsive.css")
}



// findThemeAsets() obtains a list of all non-source files.
// XXX I think I can merge this with findPageAssets()
func (a *App) findThemeAssets() {
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
				// ...and not a stylesheet, add to the list.
				if !hasExtension(filename, ".css") {
					a.Page.Theme.PageType.otherAssets = append(a.Page.Theme.PageType.otherAssets, filename)
				}
			}
		} else {
			// Special case for :
			if filename == (defaults.ThemeHelpSubdirname) {
				// Ignore it by design. This is help for the
				// user at design time. Don't want it cluttering
				// up the directory at publish time.
			}
		}
	}
}


// Look alongside the current file to assets to publish
// for example, it's a news article and it has an image.
func (a *App) findPageAssets() {
	candidates, err := ioutil.ReadDir(a.Page.dir)
	// TODO: Better error handling
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

func (a *App) titleTag() string {
	if a.FrontMatter.Title != "" {
		return a.FrontMatter.Title
	}

	node := a.markdownAST(a.Page.markdownStart)
	if t := mdext.InferTitle(node, a.Page.markdownStart); t != "" {
		return t
	}

	return defaults.ProductName + ": Title needed here, squib"
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

// layoutElementOverride() now has direction from the theme TOML to generate HTML
// for layout elements such as header, aside, etc.
// replaceWith is the extension to search for, for example,
// "sidebar" or "header".
// Look for a file in the current directory
// sharing the name of the Markdown page, so a page
// named contact.md might have an optional contact.sidebar page, for
// example. If that file is found, it takes precedence over everything
// else in the theme TOML file. For example, suppose the theme.toml
// file has this section:
//
//  [Sidebar]
//    HTML = ""
//    File = "sidebar.md"
//
// If replaceWith is "sidebar", then look (in this example) for a
// file named contact.sidebar. If that file exists, it takes
// precedents over all contents of the Sidebar part of the
// theme TOML file.
//
// Returns the contents of that file as pure HTML ready to be
// inserted into the output, or "".
func (a *App) layoutElementOverride(pr *layoutElement, tag string, replaceWith string) string {

	// Look for the matching file, so if the source Markdown
	// file is contact.md, look for a contact.sidebar file.
	elFile := replaceExtension(a.Page.filePath, replaceWith)
	if fileExists(elFile) {
		// If that .sidebar file exists, immediately
		// insert into the stream and leave,
		// because it's the highest priority.
		input := fileToBuf(elFile)
		// Generate HTML from the Markdown contents of this file.
		return wrapTag(tag, string(a.MdFileToHTMLBuffer(elFile, input)), true)
	}
	return ""
}

// layoutElementToHTML() takes an page region (header, nav, article, sidebar, or footer)
// and converts it to HTML. All we know is that it's been specified
// but we don't know whether's a Markdown file, inline HTML, whatever.
func (a *App) layoutElementToHTML(pr *layoutElement, tag string) string {
	var html string
	pathname := filepath.Join(a.Page.Theme.PageType.PathName, pr.File)
	switch tag {
	case "<header>":
		html = a.layoutElementOverride(pr, tag, "header")
	case "<footer>":
		html = a.layoutElementOverride(pr, tag, "footer")
	case "<aside>":
		html = a.layoutElementOverride(pr, tag, "sidebar")
	case "<nav>":
		html = a.layoutElementOverride(pr, tag, "navbar")
	case "<article>":
		// Exception: the theme TOML file doesnt have any entries under
		// "[article]" but there is markdown on the page. This might be
		// the most common case. Doing this allows the user to create a
		// website simply by creating a file named
		// index.md or README.md in Markdown.

		// This differs from all the other
		// layout elements in the theme TOML file, which require you to
		// specify something in File or HTML in order for them
		// to appear in the HTML output.
		// A theme TOML without an "[article]" layout element specified
		// is equivalent to <article>{{ article }}</article>.
		// So wrap the entire article in the
		// appropriate tag.
		if a.Page.Theme.PageType.Article.File == "" && a.Page.Theme.PageType.Article.HTML == "" {
			// NO article.md specified
			return wrapTag(tag, string(a.Page.Article), true)
		} else {
			// article.md specified
			html = a.layoutElementOverride(pr, tag, "article")
		}
	default:
		a.QuitError(errs.ErrCode("1203", tag))
	}

	// If nonempty, the header or whatever was generated
	// from a file in case something like this:
	//   [header]
	//   File = "header.md"
	if html != "" {
		return html
	}

	// If nonempty, the header or other page element
	// was generated from a file in case something like this:
	//   [header]
	//   HTML = "<header>Super simple header</header>"
	// Note that in this case Metabuzz doesn't supply
	// the tage, so you have to add the <header> or <nav>
	// or whatever.
	if pr.HTML != "" {
		return pr.HTML
	}

	// Skip if there's no file specified
	if pr.File == "" {
		return ""
	}

	// Pretty common case. Convert a page layout element from the
	// theme TOML file (e.g. header.md or aside.md)
	// into HTML.
	var input []byte
	// Error if the specified file can't be found.
	if !fileExists(pathname) {
		a.QuitError(errs.ErrCode("1015", pathname))
	}
	if isMarkdownFile(pathname) {
		input = fileToBuf(pathname)
		if tag == "<article>" {
			return wrapTag(tag, string(a.MdFileToHTMLBuffer(pathname, input)), true)
		} else {
			return wrapTag(tag, string(a.MdFileToHTMLBuffer(pathname, input)), true)
		}
	}
	return fileToString(pathname)
}

// metatag generates a meta tag. It's complicated.
func metatag(tag string, content string) string {
	const quote = `"`
	return ("\n<meta name=" + quote + tag + quote + " content=" + quote + content + quote + ">\n")
}

// TODO: Document this function
// TODO: Add proper error checking
// TODO: Just hold a file descripter in *App?
func (a *App) AddCommaToSearchIndex(file string) error {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	if _, err = f.Write([]byte(",\n")); err != nil {
		f.Close()
		return err
	}
	return nil
}

// xxx
// TODO: Just hold a file descripter in *App?
// TODO: Document this function
func (a *App) DelimitIndexJSON(file string, opening bool) error {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	if opening == true {
		if _, err = f.Write([]byte("[\n")); err != nil {
			f.Close()
			return err
		}
		return nil
	}
	if _, err = f.Write([]byte("]\n")); err != nil {
		f.Close()
		return err
	}
	return nil
}
