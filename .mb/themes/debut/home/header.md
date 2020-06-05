{{- $filename := index .Site.List.logomode .FrontMatter.Mode -}}
* [![Logo image]({{ $filename }}){{ .Site.Company.Name -}}]({{- .Site.Company.URL}})

