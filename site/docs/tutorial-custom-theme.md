# Tutorial: Creating a custom theme for Metabuzz

[Theme architecture](docs/theme-architecture.md)


TODO: Finish Intro stuff

Metabuzz themes can do a lot. That means they can get a little
complicated. But they don't have to be. You can start a theme
from scratch with almost no work, then add features to it as you
go. This tutorial does just that. It starts with little more
than a CSS file, then progresses to a complete Metabuzz theme
with support for a new custom pagetype, media queries, light and dark versions, header,
navbar, sidebar, and footer.

TODO: Would be nice to include
* new pagetype


TODO: 
* Explain how to find your theme directory; see theme-directory.html
* Mention -copytheme at some point
* Remove  below

```
mb new theme --from default --to empty
mb new theme --from empty --to _ 
```

## Creating a custom Metabuzz theme from scratch

The best way to create a custom Metabuzz theme is normally to modify an existing one. The `default` and `pillar` themes are good starting points but you can use any of them. For this tutorial we'll use the appropriately named `empty` theme, which keeps the boring parts of a theme, such as the stylesheet reset, and omits the fun parts, such as the design for article, header, footer, and sidebar. 

### Theme naming conventions

Metabuzz assumes your theme follows the same convention as a domain name. Simply put, that means it's case-insenstive, and allows only letters and the hypen character.

### The mb new theme command

* To create the theme named `mytheme`, use this command. Obviously replace `mytheme` with whatever your actual theme name should be.

```
mb new theme --from empty --to mytheme
```

A set of files gets created in the theme directory. Here's how to find it.

### Metabuzz theme file location: where to find your theme files

Metabuzz theme files are stored in the `.mb/themes` directory. By default that directory is found in a `metabuzz` directory where your operating system normally stores user application data. You can find out the directory location by running `mb info` on the command line: 

```
mb info
```

This displays a lot of information but the line of interest starts with `Theme directory`.

For a user named Taylor output could look something like this on MacOS:

```
/Users/taylor/Library/Application Support/metabuzz/.mb/themes
```

And this on Windows:

```
C:\Users\taylor\AppData\metabuzz\.mb\themes
```

### Changing the location of the theme file directory

You may wish to move the theme directory location. For example, it's much more convenient for all elements of a project to be in one directory if you use Git. It's easy to change. Just update the [Global configuration file](config-file.html) and move the theme directory to the new location.

## Metabuzz theme file structure

When you create a new theme, this is what gets generated at a minimum:

```
mytheme 
├── mytheme.css  
├── mytheme.toml  
├── fonts.css  
├── footer.md  
├── header.md  
├── layout.css
├── nav.md
├── reset.css
├── responsive.css
├── sidebar-left.css
├── sidebar-right.css
├── sidebar.md
├── sizes.css
├── theme-dark.css
└── theme-light.css
```







