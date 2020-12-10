# TODO
* Lots of illos when theme appearance settles down
* [Customizing a Metabuzz theme](customizing-theme.html)
* [Creating a Metabuzz theme from scratch](theme-from-scratch.html)
* [Theme architecture](theme-architecture.html)
* [Theme TOML file](theme-toml-file.html)
* [Theme directory](theme-directory.html)
* [Setting up theme](setting-up-theme.html)
* [Adding a theme to the Metabuzz gallery](add-theme.html)
* [Gallery directory](gallery-dir.html)

All Metabuzz pages have a consistent visual appearance called a *theme*.
It's defined by a set of CSS files in the [theme directory](theme-directory.html) and by behaviors imposed by the 
[Metabuzz theme framework](metabuzz-theme-framework.html).

A theme can have a "child" theme that inherits the visual traits
of the parents. It uses exactly the same set traits, which can then
be overridden. These child themes are called *pagetypes*. Internally
a theme is also a pagetype. PageTypes let you create a blog theme,
for example, but a section of the blog might be a pagetype called 
"homepage".


## Default theme

If you don't
name a theme in your page's [front matter](front-matter.html#theme), Metabuzz looks for a default theme set in the [site file](site-file.html#defaulttheme). If it can't find one there, it uses the 
[Wide](../gal/wide/index.html) theme.

## How Metabuzz generates a web page




## Specifying a theme in the front matter {#frontmatter}  

To choose a different theme, insert these lines to the very beginning of your file. In this example you'll choose the Debut theme, but it could just as well be any other theme like Pillar or Journey.

```
===
theme = "debut"
===
```

So a complete example might look like this:

```
===
theme = "debut"
===

# Debut theme test

How does the debut theme look?
```

<!-- DO NOT CHANGE THIS HEADER NAME because it is referenced elsewhere -->
## Light and dark themes


