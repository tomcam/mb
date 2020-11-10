## Creating a custom Metabuzz theme from scratch

The easiest way to create a custom Metabuzz theme is normally to modify an existing one. For this tutorial we'll use the appropriately named `empty` theme, which starts from scratch.

### Theme naming conventions

Metabuzz assumes your theme follows the same convention as a domain name. Simply put, that means it's case-insenstive, and allows only letters and the hypen character.

### The mb new theme command

* To create the theme named `mytheme`, use this command. Obviously replace `mytheme` with whatever your actual theme name should be.

```
:: Create a new theme called mytheme.
:: Base it on the existing theme named empty.
mb new theme mytheme from empty
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

All of these elements are combined with an .md source file
to create an HTML output file in the [publish directory](publish-directory.html). The `.toml` file directs that process.





