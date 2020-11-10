# Setting up your theme

# TODO: Not even close to finished. Barely started, in fact

Metabuzz themes are highly customizable if you know CSS but they're not too hard to tweak even if you don't know it too well. Here are some easy things you can do to change the appearance of your theme dramatically.


# Setting up your theme

# TODO: Not even close to finished. Barely started, in fact

Metabuzz themes are highly customizable if you know CSS but they're not too hard to tweak even if you don't know it too well. Here are some easy things you can do to change the appearance of your theme dramatically.

## Changing sizes of the header and navbar

Your [theme directory](theme-directory.html)# Theme structure customization

## Good intro to templates

Generating social media icons only if they have entries in .Site.Social. See
[template-examples.txt](template-examples.txt) for examples. Also shows how to use
-} to prevent a newline. Could put all this automatic stuff into a Masala theme.

## Colors and borders should be localized in theme-light.css and theme-dark.css

Borders because theuy have color informaiton embedded in them, though sadly you
can't search for the keyword color becaus eit's not part of the broder rule.

which unfortunately results in parallel construction, e.g.

pillar.css:

```
header > ul > li > a,
header > ul > li > a:link,
header > ul > li > a:visited
	{text-decoration:none;}
header > ul > li > a:hover,
header > ul > li > a:active
	{text-decoration:underline;}
```
theme-light.css:

```
header > ul > li > a,
header > ul > li > a:link,
header > ul > li > a:visited
	{color:var(--header-fg);}
header > ul > li > a:hover,
header > ul > li > a:active
	{color:var(--header-fg);}
```


## Theme architecture

Anything in `RootStylesheets` is copied to all pagetypes.

```
Stylesheets = ["sizes.css", "journey.css", "responsive.css"]
RootStylesheets = ["reset.css", "fonts.css", "layout.css", "theme-light.css" ]
```

  Plus of course: sidebar-left.css, sidebar-right.css, theme-dark.css
* Revise comments, especially in sizes.css

## CSS suggestions


* `p` should use margin-top more than margin-bottom because of `ol` and `ul`. Most of the time they follow a `p`., 
* `img` should have 0 or small margins for image credits/accessibility reasons. It should normally be followed by something like a caption, a photo credit, or both.
* `h1` should use margin bottom more than margin top, because the style these days
is to show a mini header over it. But `h2` and below should use margin-top more, just like `p`
* :


## Cover

Should probably have a .help directory in the theme explaining things like:
* stylessheets below
* Within that a .sample directory used for thumbnail and maybe docs

```
Stylesheets = ["sizes.css", "journey.css", "responsive.css"]
RootStylesheets = ["reset.css", "fonts.css", "layout.css", "theme-light.css" ]
```
* theme-dark and theme-light.css
* Right, left, and no sidebars
* Correct order of style sheets, e.g.
  + responsive.css comes after themename.css
  + sidebar*.css come after themename.css

## Floating page-in-a-page thing like Chill theme

$ d2 copytheme

Name of theme to copy? [default] pillar
Name of new theme? chill

edit .d2/themes/chill/skin.css:

```
--bg: #F6F6F8;
--header-fg: purple;
--header-bg: var(--bg);
--nav-fg: var(--header-bg);
--nav-bg: var(--fg);
--trim-fg: #E56236
```

edit .d2/themes/chill/chill.css
```
article {background-color:white;box-shadow: rgb(128,128,128) 1px 1px 3px 0px;} 
```

## Easiest way to handle :last-child bullets you want to suppress

Just leave the bullet blank, as in this example from pillar:

# TODO: Following must be processed because it doesn't work in markup

```
{{- /*  Automatically name first item in header    
        based on company branding, then name, then
        author name, then theme name, in descending
        order of importance.
*/ -}}
{{ if .Site.Company.Name }}
{{- $name := .Site.Company.Name -}}
{{ else if .Site.Author.FullName }}
{{- $name := .Site.Author.FullName -}} 
{{ else if .FrontMatter.Theme }}
{{- $name := .FrontMatter.Theme -}}
* [{{- $name -}}](/)
{{ end }} * [Product](/)
* [Create](/)
* [Pricing](/)
* [Resources](/)
* 
```

## Changing sizes of the header and navbar

Your [theme directory](theme-directory.html)
