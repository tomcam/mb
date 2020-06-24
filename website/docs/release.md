# How to do a release

* Obviously run unit tests...
* Check file sizes of directories in gallery
* Increment semver, maybe by inserting it to the top of a version.txt file
* The makefile should ensure a release isn't being duplicated
* Create a release with goreleaser
* Embed the latest version of the theme directory
* Manual steps
  - [Validate CSS](https://validator.w3.org/nu/#textarea) for all themes (batch version [here](https://validator.github.io/validator/#usage)
  - Check gallery in old phone, looking for
    + Responsive screen version
    + Dark version with both sidebars, and none
    + Light version with both sidebars, and none
    + Look at all pagetypes
