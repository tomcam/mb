===
theme="debut"
pagetype="gallery"
sidebar="none"

[List]
Title="METABUZZ THEME GALLERY"
DemoTheme="Personality"
Next="future"
===

# **{{ .FrontMatter.List.DemoTheme }}** theme ~~| Metabuzz~~
[![Screen shot of theme](theme-1280x1024.png)](dark.html) 
  ## {{ if .FrontMatter.List.DemoPageType }} PageType: **{{ .FrontMatter.List.DemoPageType }}**{{ else }}## {{ end }}

### DUDE

### About {{ .FrontMatter.List.DemoTheme }}
{{ inc "description.md" }}

### Live demos 

| No sidebar                | Sidebar                         |      
|:------------------------- |---------------------------------|
| [Light theme](light.html) | [Left](light-sidebar-left.html) [Right](light-sidebar-right.html)|

### Creator 
[Tom Campbell](https://metabuzz.com)

### License 
[MIT](https://metabuzz.com)


