package app

import (
	"bufio"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/plus3it/gorecurcopy"
	"github.com/spf13/viper"
	"github.com/tomcam/mb/pkg/defaults"
	"github.com/tomcam/mb/pkg/errs"
	"github.com/tomcam/mb/pkg/slices"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// cfgString() obtains a value set from a config file, environment
// variable, whatever. Simple abstraction over viper
func cfgString(option string) string {
	return viper.GetString(option)
}

// cfgBool() obtains a value set from a config file, environment
// variable, whatever. Simple abstraction over viper
func cfgBool(option string) bool {
	return viper.GetBool(option)
}

// configDir() returns the user's application configuration directory,
// or just "." for the current directory if it can't be
// determined through system calls.
// https://golang.org/pkg/os/#UserConfigDir
func configDir() string {
	if cfgDir, err := os.UserConfigDir(); err != nil {
		// Can't determine from the OS, so just use the current directory.
		return filepath.Join(".", defaults.GlobalConfigurationDirName)
	} else {
		// Got an actual valid global application data directory
		return filepath.Join(cfgDir, defaults.ProductName, defaults.GlobalConfigurationDirName)
	}
}

// copyDirOnly() copies a directory nonrecursively.
// Doesn't other directories
// Thanks to https://github.com/plus3it/gorecurcopy/blob/master/gorecurcopy.go
func copyDirOnly(source, dest string) error {
  err := os.MkdirAll(dest, defaults.PublicFilePermissions)
	if err != nil {
		return errs.ErrCode("0410",err.Error(),dest)
	}
	entries, err := ioutil.ReadDir(source)
	if err != nil {
		// TODO: More specific errors
		return errs.ErrCode("0708",err.Error(),source)
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(source, entry.Name())
		destPath := filepath.Join(dest, entry.Name())
		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
		  return errs.ErrCode("0129",err.Error(),sourcePath)
		}
		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			// Do nothing. What's the syntax for that?
		case os.ModeSymlink:
			if err := CopySymLink(sourcePath, destPath); err != nil {
		    return errs.ErrCode("0130",err.Error(),sourcePath)
			}
		default:
			if err := Copy(sourcePath, destPath); err != nil {
		    return errs.ErrCode("0214",err.Error(),sourcePath)
			}
		}
	}
	return nil
}

// copyDirAll() does a recursive copy of the directory and its subdirectories
func copyDirAll(source, dest string) error {
	if source == "" {
		return errs.ErrCode("0704", source)
	}
	if dest == "" {
		return errs.ErrCode("0705", dest)
	}

	if dest == source {
		return errs.ErrCode("0707", "from '"+source+"' to '"+dest+"'")
	}

	err := gorecurcopy.CopyDirectory(source, dest)
	//return ErrCode("0406","from '"+source+"' to '"+dest+"'",err.Error())
	if err != nil {
		return errs.ErrCode("0406", "from '"+source+"' to '"+dest+"'", "")
	}
	return nil
}

// copyNewName() copies file source to file target, then deletes source
func copyNewName(source, target string) error {
	if !fileExists(source) {
		return errs.ErrCode("PREVIOUS", "")
	}
	if err := Copy(source, target); err != nil {
		return errs.ErrCode("0111", "trying to copy "+source+" to "+target)
	}
	if !fileExists(target) {
		return errs.ErrCode("0110", "tried to create "+target)
	}
	if deleteFile(source) != nil {
		return errs.ErrCode("0301", source)
	}
	return nil
}

func CopySymLink(source, dest string) error {
	link, err := os.Readlink(source)
	if err != nil {
		return err
	}
	return os.Symlink(link, dest)
}

// Copy() does just that. It copies a single file named source to
// the file named in dest.
func Copy(src, dest string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return errs.ErrCode("0112", src)
	}

	if !sourceFileStat.Mode().IsRegular() {
		return errs.ErrCode("0113", src)
	}
	source, err := os.Open(src)
	if err != nil {
		return errs.ErrCode("0114", src)
	}
	destination, err := os.Create(dest)
	if err != nil {
		return errs.ErrCode("0209", dest)
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	if err != nil {
		return errs.ErrCode("0251", dest)
	}
	// Success
	return nil
}

// createDirStructuure() Creates the specified site structure
// in the current directory.
func createDirStructure(dirs *[][]string) (err error) {
	// Obtain current directory in a portable way.
	basedir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Build up a directory tree for each row
	// in dirs
	for _, row := range *dirs {
		path := basedir
		for _, subdir := range row {
			// Append the next subdirectory in the path
			// in a portable way
			path = filepath.Join(path, subdir)
		}
		err := os.MkdirAll(path, defaults.PublicFilePermissions)
		if err != nil {
			return err
		}
	}
	return nil
}

// curDir() returns the current directory name.
func currDir() string {
	if path, err := os.Getwd(); err != nil {
		return "unknown directory"
	} else {
		return path
	}
}

// deleteFile() returns true if it can remove the named file
func deleteFile(filename string) error {
	return os.Remove(filename)
}

// deleteFileMust() tries to remove the named file but doesn't
// return an error if there's a problem.
func deleteFileMust(filename string) {
	if !fileExists(filename) {
		return
	}
	_ = os.Remove(filename)
}

// dirExists() returns true if the name passed to it is a directory.
func dirExists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	} else {
		return false
	}
}

// fileExists() returns true, well, if the named file exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// fileToBuf() reads the named file into a byte slice and returns
// that byte slice. In the spirit of HTML it simply returns an empty
// slice on failure.
func fileToBuf(filename string) []byte {
	var input []byte
	var err error
	// Read the whole file into memory as a byte slice.
	input, err = ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}
	}
	return input
}

// fileToString() sucks up a file and returns its contents as a string.
// Fails quietly  if unable to open the file, since
// we're just generating HTML.
func fileToString(infile string) string {
	input, err := ioutil.ReadFile(infile)
	if err != nil {
		return ""
	}
	return string(input)
}

// hasExtension() returns true if the string ends in the specified extension
// (case insensitive). Need to supply the period too:
// if hasExtension(filename, ".aside") {
func hasExtension(filename, extension string) bool {
	return filepath.Ext(filename) == extension
}

// hasExtension() Returns true if the fully qualified filename
// ends in any of the extensions listed in extensions.
func hasExtensionFrom(path string, extensions *slices.SearchInfo) bool {
	return extensions.Contains(filepath.Ext(path))
}

/*
// hasExcludedExtension() returns true if the supplied string
// is on the list of excluded filetypes.
func hasExcludedExtension(filename string) bool {
	return hasExtensionFrom(filename, ExcludedAssetExtensions)
}

// hasMarkdownExtension() returns true if the supplied string
// (presumably a filename) ends in ".md", ".markdown", etc.
func hasMarkdownExtension(filename string) bool {
	return hasExtensionFrom(filename, MarkdownExtensions)
}
*/

// homeDir() returns the user's home directory, or just "." for
// the current directory if it can't be determined through system
// calls.
func homeDir() string {
	var home string
	var err error
	if home, err = os.UserHomeDir(); err != nil {
		return "."
	}
	return home
}

// inputString() gets a string from the keyboard and returns it
// See also promptString()
func inputString() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

// isDirectory() returns true if the specified
// path name is a directory.
func isDirectory(pathName string) bool {
	f, err := os.Stat(pathName)
	if err != nil {
		return false
	}
	return f.Mode().IsDir()
}

// isFile() returns true if the specified
// path name is an existing file.
func isFile(pathName string) bool {
	f, err := os.Stat(pathName)
	if err != nil {
		return false
	}
	return f.Mode().IsRegular()
}

// isMarkdownFile() returns true of the specified filename has one of
// the extensions used for Markdown files.
func isMarkdownFile(filename string) bool {
	return hasExtensionFrom(filename, defaults.MarkdownExtensions)
}

// isProject() looks at the structure of the specified directory
// and tries to determine if there's already a project here.
// It does so by looking for site config subdirectory.
func isProject(path string) bool {
	// If the directory doesn't exist, that's easy.
	if !dirExists(path) {
		return false
	}

	// The directory exists. Does it contain a site directory?
	return isSiteDir(path)

}

// isSiteDir() looks for the special name used for the subdirectory
// used to hold site config file & info
func isSiteDir(path string) bool {
	return dirExists(siteDir(path))
}

// promptString() displays a prompt, then awaits for keyboard
// input and returns it on completion.
// See also inputString(), promptYes()
func promptString(prompt string) string {
	fmt.Print(prompt + " ")
	return inputString()
}

// promptStringDefault() displays a prompt, then awaits for keyboard
// input and returns it on completion. It precedes the end of the
// prompt with a default value in brackets.
// See also inputString(), promptYes()
func promptStringDefault(prompt string, defaultValue string) string {
	fmt.Print(prompt + " [" + defaultValue + "] ")
	answer := inputString()
	if answer == "" {
		return defaultValue
	} else {
		return answer
	}
}

// promptYes() displays a prompt, then awaits
// keyboard input. If the answer starts with Y,
// returns true. Otherwise, returns false.
// See also inputString(), promptString()
func promptYes(prompt string) bool {
	// See also inputString(), promptYes()
	answer := promptString(prompt)
	return strings.HasPrefix(strings.ToLower(answer), "y")
}

// relativeDirectory() takes a base directory,
// for example, /users/tom/mysite, and a subdirectory, for
// example, /users/tom/mysite/articles/
// and returns the relative directory, which would be
// the directory named /articles in this case.
func relativeDirectory(baseDir, subDir string) string {
	if baseDir == subDir {
		return ""
	}
	// Begin at the end of the base directory
	start := len(baseDir)
	// Extract the target directory from the
	// input filename
	stop := len(subDir)
	// End at the beginning of the filename
	//stop := l - start
	// TODO: Playing with fire?
	if stop <= 0 {
		stop = start
		return string(subDir[start:stop])
	}
	return string(subDir[start+1 : stop])
}

// relDirFile() takes a base directory,
// for example, /users/tom/mysite, and a filename, for
// example, /users/tom/mysite/articles/announce.md,
// and returns the relative directory, which would be
// the directory named /articles in this case.
func relDirFile(baseDir, filename string) string {
	// Begin at the end of the base directory
	// xxx
	start := len(baseDir)
	// Extract the target directory from the
	// input filename
	l := len(filepath.Dir(filename))
	// End at the beginning of the filename
	stop := l - start
	// TODO: Playing with fire?
	if stop < 0 {
		stop = start
	}
	return string(filename[start : start+stop])
}

// replaceExtension() is passed a filename and returns a filename
// with the specified extension.
func replaceExtension(filename string, newExtension string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename)) + "." + newExtension

}

// siteDir() returns the expected name of a site subdirectory  in
// the given path.
func siteDir(path string) string {
	return filepath.Join(path, defaults.GlobalConfigurationDirName, defaults.SiteConfigDir)
}

// userName() returns the
// current user name according to the OS,
// or "" if none is found.
func userName() string {
	u, err := user.Current()
	if err != nil {
		return ""
	}
	return u.Username
}

// pressEnter() just displays the optional prompt watis for hte user to press enter.. It doesn't return anthing.
// input and returns it on completion.
func pressEnter(prompt string) {
	fmt.Print(prompt + " ")
	_ = inputString()
}

// WriteTextFile creates a file called filename without checking to see if it
// exists, then writes contents to it.
func writeTextFile(filename, contents string) error {
	var out *os.File
	var err error
	if out, err = os.Create(filename); err != nil {
		return errs.ErrCode("0204", "Problem creating file %v: %v\n", filename, err.Error())
	}
	if _, err = out.WriteString(contents); err != nil {
		return errs.ErrCode("0903", "Problem writing to file %v: %v\n", filename, err.Error())
	}
	return nil
}

// TODO: This fails in its only uses, the readPrefs() call. Figure it out.
func readTomlFile(filename string, target interface{}) (err error) {
	var input []byte
	if input, err = ioutil.ReadFile(filename); err != nil {
		return err
	}

	if _, err = toml.Decode(string(input), target); err != nil {
		return err
	}

	return nil
}

// writeTomlFile() creates a TOML file based on the filename and
// data structure passed in.
func writeTomlFile(filename string, target interface{}) error {
	f, err := os.Create(filename)
	if err != nil {
		return errs.ErrCode("0210", err.Error(), filename)
	}
	if err = toml.NewEncoder(f).Encode(target); err != nil {
		return errs.ErrCode("0908", err.Error())
	}
	if err := f.Close(); err != nil {
		return errs.ErrCode("0252", filename)
	}
	return nil
}
