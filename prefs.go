package main

// Prefs contains information about the user's system, such as the user's
// name and home directory. Also global preferences, such
// as which theme to use as the default or where the
// shortcode directory should be.
type Prefs struct {
	// Where to find global values such as themes.
	// Determined at startup.
	configDir string
}
