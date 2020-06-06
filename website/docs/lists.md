===
templates="off"
===
# Lists

The `[List]` section in the front matter of a MarkDown file is a powerful
way to do a kind of search and replace. It's sort of a tiny database,
and lets you give a name to things like
a name and address, job title, or even multiple kinds kind of text in one place.
If they have to be changed later, just change them in the `[List]` section.

Technically, the `[List]` section is a collection of key/value pairs
expressed in TOML format. The [TOML specification](https://github.com/toml-lang/toml)
and can be found along with its source code on [GitHub](https://github.com/toml-lang/toml). 

Metabuzz normally uses a drastically simplified subset, however, and that's what this documentation covers.

The `[List]` section is optional. 

## IMPORTANT! Always ensure the [List] section ends the file

You'll be warned about this later, but it's important that the `[List]`
appears last in your file, after things like `title`, `excludedfiles`, and
so on.

## Example file

Here is an example Markdown file containing

##### file listex01.md

```
===
\[List\]
title="Metabuzz: The Greatest Hits"
books=["Cujo","Perdido"]
mixed=["Double quoted",'single-quoted','"mixed" quotes']
biglist= [
  ["beat","me"],
  ["with","a","stick"],
]
===

Welcome to {{ .FrontMatter.title }}
```



