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

### Try all versions of this theme live 

| No sidebar                | Sidebar                         |                                  |
|:------------------------- |---------------------------------|----------------------------------|
| [Light theme](light.html) | [Left](light-sidebar-left.html) | [Right](light-sidebar-right.html)|
| [Dark theme](dark.html)   | [Left](dark-sidebar-left.html)  | [Right](dark-sidebar-right.html) |



### About {{ .FrontMatter.List.DemoTheme }}
{{ inc "description.md" }}

#### CREATOR [Tom Campbell](https://metabuzz.com)
#### LICENSE [MIT](https://metabuzz.com)


