# Headers

## header.md

Most Metabuzz themes provide a `header.md` file in the [theme directory](theme-directory.html) that typically looks something like this:

```
* [Default](/)
* [Create](/)
* [Pricing](/)
* [Try it Free](/)
```

```markdown
{{"{{"}} if .Site.Company.Name {{"}}"}}
{{"{{-"}} $name := .Site.Company.Name {{"-}}"}}
* [{{"{{-"}} $name {{"-}}"}}](/)
{{"{{-"}} else if .Site.Author.FullName {{"-}}"}}
{{"{{-"}} $name := .Site.Author.FullName {{"-}}"}}
* [{{"{{"}} $name {{"-}}"}}](/)
{{"{{-"}} else {{"}}"}}
* [{{.FrontMatter.Theme}}](/)
{{"{{-"}} end {{"}}"}} 
* [Gallery](/)
* [Docs](/)
* [Download](/)
* [About](/)
```


* How to use the Debut theme
  + Typography
    - Start of page only: Follow first h1 with an h2
  + Images
    - Displaying image on left with text to the right 


