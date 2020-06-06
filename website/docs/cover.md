# Things to cover in the documentation

## Extensions

Explain goldmark extensions

 with bookmarks/anchors/ID attributes

## Themes

New architecture

```
Stylesheets = ["reset.css", "fonts.css", "sizes.css", "layout.css", "theme-light.css", "home.css", "responsive.css"]
```

Plus of course: sidebar-left.css, sidebar-right.css, theme-dark.css


* Ideally each theme would be documented with:
  - A visual guide showing what the header, navbar, etc. are
  - A visual guide showing each part of the client area
  - Expected use of Markdown for that particular theme, e.g. Debut Home likes an h1 followed by an h2
  - Complete documentation of special theme features
  - A slide show/gallery showing 
    + Default theme
    + Default theme dark
    + Default theme with sidebar it implements


## Tutorials

* Checklists
* Footnotes https://michelf.ca/projects/php-markdown/extra/#footnotes 

```
Visit [^quoted]
...

[^quoted]:

Thanks for making it this far
```

## Markdown

There's no underline convention in Commonmark!

When discussing free use of HTML, bring up symbols like &middot; , the cophright one and simliar
## Walkthroughs

* Walkthroughs should cover things you only do once (or occasionally) vs things you'll
do every day
* Have a complete walkthrough getting it on the web with GitHub and also with Netlify




