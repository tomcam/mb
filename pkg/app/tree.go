package app

import (
	"github.com/tomcam/mb/pkg/defaults"
	"github.com/tomcam/mb/pkg/errs"
	"github.com/tomcam/mb/pkg/slices"
	"os"
	"path/filepath"
	"strings"
)

// Called by  getProjectTree()
// Builds a list all files and all directories in the project.
// Excludes the assets directory and the publish directory.
func (a *App) visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		// Find out what directories to exclude
		exclude := slices.NewSearchInfo(a.excludeDirs())
		if err != nil {
			// Quietly fail if unable to access path.
			return err
		}
		isDir := info.IsDir()

		// Skip any directory to be excluded, such as
		// the pub directory itself and anything
		// the user specified in the siteConfig's "Exclude"
		name := info.Name()

		// Exclude this directory if it starts with "."
		if strings.HasPrefix(name, ".") && isDir {
			return filepath.SkipDir
		}

		// Exclude this directory if found on the, ah, exclusion list.
		if exclude.Contains(name) && isDir {
			a.Verbose("Excluding directory: %s", name)
			return filepath.SkipDir

		}

		if exclude.Contains(name) {
			a.Verbose("Excluding: %s", name)
			return nil
		}

		if isDir {
			a.setMdOption(path, NormalDir)
			// a.Site.dirs[path].MdOptions = NormalDir
		}

		*files = append(*files, path)
		return nil
	}

}

// Obtain a list of all files in the specified project tree starting
// at the root.
// Ignore directories starting with a .
// Ignore the assets directory
func (a *App) getProjectTree(path string) (tree []string, err error) {
	// This should only be called once so I imagine the
	// following is unnecessary
	var files []string
	err = filepath.Walk(path, a.visit(&files))
	if err != nil {
		return []string{}, errs.ErrCode("0702", err.Error(), path)
	}
	// fmt.Fprintf(os.Stdout, "Directory tree for %+v\n", files)
	return files, nil
}

// Returns a list of files and directories to be excluded from the source directory when the
// project is built. It's based on internal configuration (for example, it excludes the
// publish directory) and any existing excludes (for example, Exclude=["pub", "node_modules"])
// in the site config file.
func (a *App) excludeDirs() []string {
	// fmt.Println("Excluded in site.toml:", a.Site.ExcludeDirs)
	// Add the publish directory if it isn't already there.
	return append(a.Site.ExcludeDirs,
		defaults.CommonDir,
		defaults.HeadTagsDir,
		defaults.PublishDir,
		defaults.SCodeDir,
		defaults.SiteConfigDir,
		defaults.ThemeDir)
}
