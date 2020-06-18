package main

import "fmt"

// generateTOC() reads the Markdown source and returns a string
// array containing each header and its level in the document.
func (App *App) generateTOC(level int) {
	// Actual markdown source is a byte slice at App.Page.markdownStart
	//fmt.Printf("Markdown source: %v\n", string(App.Page.markdownStart))

	// Generate fake TOC
	var entry TOCEntry
	for i := 0; i < level; i++ {
		entry.header = fmt.Sprintf("Fake header level %v", i+1)
		entry.level = i
		App.Page.TOC = append(App.Page.TOC, entry)
	}
}
