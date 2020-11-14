{{- /*  Automatically name first item in header    
        based on company name, author name name
        if no company was specified, or just 
        the name of the theme if neither of those
        was specified.
        
*/ -}}
![logo](square-image-32x32.png){{- if .Site.Company.Name -}}
{{- $name := .Site.Company.Name -}}
* [{{ $name -}}](/)
{{- else if .Site.Author.FullName -}}
{{- $name := .Site.Author.FullName -}}
* [{{ $name -}}](/)
{{- else }}
* [{{.FrontMatter.Theme}}](/)
{{- end }} 

