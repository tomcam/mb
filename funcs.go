package main

import (
	/*
		"html/template"
		"io/ioutil"
		"os"
		"path/filepath"
		"strings"
	*/
	"time"
)

// Return the current, local, formatted time.
// Can pass in a formatting string
// https://golang.org/pkg/time/#Time.Format
func (App *App) ftime(param ...string) string {
	var ref = "Mon Jan 2 15:04:05 -0700 MST 2006"
	var format string

	if len(param) < 1 {
		format = ref
	} else {
		format = param[0]
	}
	t := time.Now()
	return t.Format(format)
}
