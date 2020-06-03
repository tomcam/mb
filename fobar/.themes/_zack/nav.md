{{$imgurl:=.Site.List.Iconurl}}
{{with .Site.List.gallery}}{{range $key, $value := .}}[![{{$value.alt}}]({{$imgurl}}{{$value.filename}}?&size=20&color=FFFFFF)]({{$value.url}}) {{end}}{{end}}
[HOME](/) [NEWS](/) [SOCIAL](/) [ABOUT](/) 

