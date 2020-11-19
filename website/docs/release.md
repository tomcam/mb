# How to do a release

* Fill out as many parts of the Metabuzz site file as possible.
* Make sure website/docs/site-file-example.md is up to date 
Put it in website/docs/site-file-example.md
* Obviously run unit tests...
* Check file sizes of directories in gallery
* Increment semver, maybe by inserting it to the top of a version.txt file
* The makefile should ensure a release isn't being duplicated
* Search for xxx in source
* Create a release with goreleaser
* Embed the latest version of the theme directory
* Manual steps
  - [Validate CSS](https://validator.w3.org/nu/#textarea) for all themes (batch version [here](https://validator.github.io/validator/#usage)
  - Get rid of VNU_CHECK.HTML files or ideaaly start using temp files https://stackoverflow.com/questions/10982911/creating-temporary-files-in-bash
  - Check gallery in old phone, looking for
    + Responsive screen version
    + Dark version with both sidebars, and none
    + Light version with both sidebars, and none
    + Look at all pagetypes
  - Ideally, remove all unused files from each theme directory
