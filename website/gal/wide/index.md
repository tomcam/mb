===
theme="debut"
pagetype="gallery"
sidebar="none"

[List]
Title="METABUZZ THEME GALLERY"
DemoTheme="wide"
Next="future"
===

# **{{ .FrontMatter.List.DemoTheme }}** theme
* ![Screen shot of Wide theme](theme-1280x1024.png)
  ## {{ if .FrontMatter.List.DemoPageType }} PageType: **{{ .FrontMatter.List.DemoPageType }}**{{ else }}## {{ end }}
  An exceptionally lightweight, general-purpose theme with high information density and maximum flexibility.   
  ### Modes
  [Light theme](light.html) [Dark theme](dark.html)
  ### Sidebar support
  Light theme: [Left](light-sidebar-left.html) [Right](light-sidebar-right.html)  
  Dark theme: [Left](dark-sidebar-left.html) [Right](dark-sidebar-right.html) 
  #### CREATOR [Tom Campbell](https://metabuzz.com)
  #### LICENSE [MIT](https://metabuzz.com)
  ### Next: [{{ .FrontMatter.List.Next }}](../{{- .FrontMatter.List.Next -}}/index.html) 

# **{{ .FrontMatter.List.DemoTheme }}** theme
* ![Screen shot of Default dark theme with right sidebar](theme-dark-right-1280x1024.png)
  ## **Mode:** Dark
  ## **Sidebar:** Right

# **{{ .FrontMatter.List.DemoTheme }}** theme
* ![Screen shot of light theme with left sidebar](theme-light-left-1280x1024.png)
  ## **Mode:** Light 
  ## **Sidebar:** Left

# **{{ .FrontMatter.List.DemoTheme }}** theme
* ![Screen shot of dark theme with no sidebar](theme-dark-nosidebar-1280x1024.png)
  ## **Mode:** Dark
  ## **Sidebar:** None 

# **{{ .FrontMatter.List.DemoTheme }}** theme
* ![Screen shot of light theme with no sidebar](theme-light-nosidebar-1280x1024.png)
  ## **Mode:** Light
  ## **Sidebar:** None 

