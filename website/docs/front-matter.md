===
templates="off"
===

# TODO:
* No empty front matter allowed


## excludedFiles

List of files in the current directory you don't want
copied to the Publish directory.
Must be literal filenames, not wildcards.

```
===
excludedfiles = [ "clientid.src", "productkey.txt" ]
===

# How to use your product key

Remember to keep your product keyk secret.

```

## PageType
See also [Theme](#theme)

## Theme
Allows you to set the visual appearance on a per-page basis
by naming a Metabuzz [theme](themes.html). If you don't
name a theme, Metabuzz looks for a default theme
set in the [site file](site-file.html#defaulttheme). If
it can't find one there, it uses the 
[default theme](themes.html#default-theme).


See also [PageType](#pagetype)

## List

TODO: Explain it needs to be last

## templates

For documentation purposes. If you're writing documentation that uses the template language, setting `templates="off"` prevents templates on that page from
being executed. Helps when you're documenting, well, templates.

```
===
templates="off"
===
```

For example, since there's no template function called `world` this
would normally produce an [0917]error if used in your Markdown, but if you 
set `templates="off"` you won't have that problem.
```
hello, {{ world. }}
```
