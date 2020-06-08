===
execute="off"
===

# Adding themes to the gallery directory

The gallery directory is at `mb/website/gal` Each theme has its own directory, so, for example, the Wide theme is at `mb/website/gal/wide`.

## Editing index.md

The theme directory in the gallery's `index.md` file is shown below. You need to change:

* The `theme=` line
* The `pagetype` line, if any.
* The `DemoTheme` branding name

You need to create
* A 1280x1024 image of the dark version of the theme, if any, named `theme-{theme name}-dark-1280x1024.png`


## How the demo directory works

The `demo` directory starts with `index.md`, which 

## Things I did
- deleted "theme=" from front matter of:
  - dark-sidebar-left.md
  - dark-sidebar-right.md

## Editing files in the demo directory

* edit index.md to use the correct theme
* Replace description.md to suggest uses for the theme and why it was created
* Put sample images in there for the sidebar and intro.md
* Replace intro.md for inserted text
* Replace left-sidebar-example.md
* Replace right-sidebar-example.md

### Files to leave untouched in the demo directory
* dark-sidebar-left.md
* light-sidebar-left.md
* dark-sidebar-left.sidebar
* dark-sidebar-right.sidebar
* dark.md


