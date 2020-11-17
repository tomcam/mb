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
