===
theme="debut"
pagetype="gallery"
sidebar="none"

[List]
Title="METABUZZ THEME GALLERY"
DemoTheme="Wide"
Next="future"
===

# **{{ .FrontMatter.List.DemoTheme }}** theme ~~| Metabuzz~~
[![Screen shot of theme](theme-1280x1024.png)](dark.html) 
  ## {{ if .FrontMatter.List.DemoPageType }} PageType: **{{ .FrontMatter.List.DemoPageType }}**{{ else }}## {{ end }}

### About {{ .FrontMatter.List.DemoTheme }}
{{ inc "description.md" }}

{{ inc "variations.md" }}

### Creator 
[Tom Campbell](https://metabuzz.com)

### License 
[MIT](https://metabuzz.com)


