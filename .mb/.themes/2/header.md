{{- /*  The first list item has 2 elements: the name
        of an optional logo file and then the first
        item on the unordered list, which gets
        special branding. If the user has specified 
        the name of of a logo fie in .Site.Company.HeaderLogo
        than that file gets displayed. Otherwise it's skipped.
      
        That's followed by name, which is chosen from the
        following sources in order:

        * Under "[List]" at the bottom of your page's front matter,
          an entry in the form title="your title here":

          [List]
          Title="replace this"

        * Under "[Company] in your site file, an entry
          in the form name="company name":

          [Company]
          Name="replace this"

        * Under "[Author] in your site file, an entry
          in the form fullname="author name":

          [Author]
          FullName="replace this"

           Automatically name first item in header    
        
        * If nothing else is specified the name of the theme
          appears in its place

*/ -}}
{{- if .Site.Company.HeaderLogo -}}
![logo]({{- .Site.Company.HeaderLogo -}})
{{- end -}}
{{- if .FrontMatter.List.Title -}}
{{ $name := .FrontMatter.List.Title }}
* [{{ $name }}](/)
{{- else if .Site.Company.Name }}
{{ $name := .Site.Company.Name }}
* [{{ $name }}](/)
{{ else if .Site.Author.FullName }}
{{ $name := .Site.Author.FullName }}
* [{{ $name }}](/)
{{- else }}
* [{{.FrontMatter.Theme}}](/)
{{- end }} 



{{- if .Site.Company.Name -}}
{{- $name := .Site.Company.Name -}}
* [{{ $name -}}](/)
{{- else if .Site.Author.FullName -}}
{{- $name := .Site.Author.FullName -}}
* [{{ $name -}}](/)
{{- else }}
* [{{.FrontMatter.Theme}}](/)
{{- end }} 
* [Create](/)
* [Pricing](/)
* [Try it Free](/)




