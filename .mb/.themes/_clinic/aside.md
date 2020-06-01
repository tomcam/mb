{{$imgurl:=.Site.List.Iconurl}}{{$color:=.Site.List.black}}{{$size:=.Site.List.Iconsize}}
{{with .Site.List.gallery}}
{{range $key, $value := .}}
[![{{$value.alt}}]({{$imgurl}}{{$value.filename}}?&size={{$size}}&color={{$color}})]({{$value.url}})
{{end}}{{end}}

