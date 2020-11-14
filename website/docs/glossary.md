# Glossary

## Commonmark

See [Markdown](#markdown)

## Global configuration file

The [global configuration file](config-file.html) is a file named `metabuzz.toml` normally stored in a subdirectory named `.mb` that contains information that applies to all projects you create with Metabuzz---for example, where the theme files are stored.

## Layout element

## Markdown

Markdown is a sensible way to represent text files so that they read easily as plain text if printed out, but also carry enough semantic meaning that they can be converted into HTML. 
The closest thing to an industry standard for Markdown is [CommonMark](https://commonmark.org). Metabuzz converts all CommonMark text according to specification, and includes extensions for things like tables, strikethrough, and autolinks. See the source to [Goldmark](https://github.com/yuin/goldmark) for more information on extensions.

Take this example of Markdown you might use in a document:

```
# Introduction

*hello*, world
```

The above would be converted in HTML that looks like this.

```
<h1>Introduction</h1>
<p><em>hello</em>, world.</p>
```

That means the `# Introduction` actually represents the HTML heading type `h1`, which is the hightest level of organization. `## Introduction` would generated an `h2` header, and so on. 

The asterisk characters are replaced by the `<em>` tag pair, which means they have the semantic power of emphasis. This is represented by HTML as italics, although you could override it in CSS.


## Pagetype

See also [theme](#theme)

## Project

See [site](#site)


## Publish directory

The [publish directory](publish-directory.html) contains your website: the set of HTML files, theme files, CSS, image, sound, and other assets that Metabuzz generates from your Markdown files and other assets. It is in the `.pub` subdirectory immediately off the root of your [site directory](site.html). Ultimately it will be copied to the WWW or whatever directory your web host uses to publish HTML files from.


## README.md

`README.md` has a special property. If there's a Markdown file named `README.md` in a directory, it gets renamed to `index.html` and becomes the "home page" for that directory, even if there's already an `index.md` file. Wait, what? It's because [GitHub](https://guides.github.com/features/wikis/) uses that convention for `README.md` and GitHub is the big dog. Hey, we didn't make the rules.

<a id="site"></a>
## Site, aka site directory

*Project* and *site* normallly have the same informal meaning: a [directory](site-directory.html) containing all themes, the sites Markdown documents, graphic assets, stylesheets, and related files required to create a website. It's created automatically when you use the [build command](tutorial01.html#building-your-site).

## Site configuration file 

The [site configuration file](site-file.html), also called simply the *site file*, holds information about the project you're working on. It's a file named `.site.toml` stored in a subdirectory of your project named `.site`. Example of site-specific data includes the company name, the URL of the site, the author of the site, and so on.

You can have as many site files as you want. They are completely independent, so you can create all the websites you want as long as the Markdown and other files go in different directories.


## Template language

The [template language](template-language.html) doesn't refer to themes, which in some content management systems are called templates. Instead, the template language is a text-replacment system that adds features to your website that can't be added using pure HTML. Metabuzz uses the [Go template package](https://golang.org/pkg/text/template/) unchanged, so if you have any questions that aren't handled by the Metabuzz documentation you'll find it either there or in the [Go template package source code](https://golang.org/src/text/template/template.go).


## Theme

Every Metabuzz site has a [theme](themes.md), which is a collection of stylesheets and graphic images structured in a particular way. A theme has its own folder, which is used as the name of the theme, and a confguration file listing what files comprise the theme. If you haven't specified a theme in your [site file](#site-configuration-file) or page [front matter](#front-matter) then the theme named Default is used.

A theme is technically a [pagetype](#pagetype). The only difference between the two is that the theme may contain other pagetypes. For example, a blog-oriented theme might have a home pagetype and a blog pagetype.

See also [pagetype](#pagetype)

## Theme TOML file

Each Markdown file

[layout elements](#layout-element)

```
[Header]
  HTML = ""
  File = "header.md"

[Nav] 
  HTML = ""
  File = "nav.md"

[Article]
  HTML = ""
  File = ""

[Sidebar]
  HTML = ""
  File = "sidebar.md"

[Footer]
  HTML = ""
  File = "footer.md"
```




