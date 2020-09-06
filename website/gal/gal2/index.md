===
theme="debut"
pagetype="gallery"
sidebar="none"

[List]
Title="METABUZZ THEME GALLERY"
DemoTheme="gal2"
===
# **{{ .FrontMatter.List.DemoTheme }}** theme 

###
* ![Thumbnail screenshot of ](light-sidebar-left-256x256.png)
* ![Thumbnail screenshot of ](light-sidebar-right-256x256.png)
* ![Thumbnail screenshot of ](light-sidebar-none-256x256.png)
* ![Thumbnail screenshot of ](dark-sidebar-left-256x256.png)
* ![Thumbnail screenshot of ](dark-sidebar-right-256x256.png)
* ![Thumbnail screenshot of ](dark-sidebar-none-256x256.png)


[![Screen shot of theme](theme-1280x1280.png)](dark.html) 
  ## {{ if .FrontMatter.List.DemoPageType }} PageType: **{{ .FrontMatter.List.DemoPageType }}**{{ else }}## {{ end }}

### About {{ .FrontMatter.List.DemoTheme }}
{{ inc "description.md" }}

{{ inc "variations.md" }}

### Creator 
[Tom Campbell](https://metabuzz.com)

### License 
[MIT](https://metabuzz.com)


