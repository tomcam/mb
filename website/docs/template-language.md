===
templates="off"
===
# TODO: Check these code samples in the Markdown source


## Show clickable icon only if you have the appropriate entry in the Site file

```
{{ if .Site.Social.Twitter }}[![Twitter logo](twitter-gray-30x30.svg)]({{ .Site.Social.Twitter }}){{ end -}}
```
More complicated:

```
[Learn More](/) [Download the App](/) [Contact](/)  [Media](/) {{ if .Site.Social.Twitter }}[![Twitter logo](twitter-blue-30x30.svg)]({{ .Site.Social.Twitter }}){{ end -}} {{ if .Site.Social.Facebook }}[![Facebook logo](facebook-blue-30x30.svg)]({{ .Site.Social.Facebook }}){{ end -}}{{ if .Site.Social.LinkedIn }}[![LinkedIn logo](linkedin-blue-30x30.svg)]({{ .Site.Social.LinkedIn }}){{ end -}} {{ if .Site.Social.YouTube }}[![YouTube logo](youtube-red-30x30.svg)]({{ .Site.Social.YouTube}}){{ end -}} 
```

Inline mention: `{{"{{"}} ftime "3:04pm" {{"}}"}}` 

Code fence:


```
{{"{{"}} ftime "3:04pm" {{"}}"}}
``` 

# Template language

While Markdown and HTML can take you a long way, there are still some things they can't do. Metabuzz extends them with a "template" language, created as part of the [Go language](https://golang.org/pkg/text/template/) that Metabuzz is written in.

TODO: a bit more intro matter

## Hands-off themes 

Metabuzz ships with some themes that you can personalize by putting text in one place and never having to touch the templates at all. Here's how that works.

Suppose your navbar, sidebar, and footer all have image links to your Twitter account and your Twitter handle is `trymetabuzz`. One way to update them would be to go into `navbar.md`, `sidebar.md`, and `footer.md` and change them to reflect it like this: 

```
[![Twitter logo](twitter-blue-30x30.svg)](https://twitter.com/trymetabuzz)
```

But there's an easier way. The hands-off templates have image links that look like this:

```
[![Twitter logo](twitter-blue-30x30.svg)]({{ \.Site.Social.Twitter }})
```



```
{{"{{"}} $author := index .Site.Authors 0 {{"-}}"}}
{{"{{-"}} if .Site.Company.Name  {{"}}"}}                                                 * [{{"{{-"}} .Site.Company.Name {{"}}"}}](/)                                              {{"{{-"}} else if $author.FullName {{"}}"}}
* [{{"{{-"}} $author.FullName {{"}}"}}](/)
{{"{{-"}} else {{"}}"}}
* [{{"{{-"}} .FrontMatter.Theme {{"}}"}}](/)
{{"{{-"}} end {{"}}"}}
* [Events](/)
* [Podcast](/)
* [Subscribe](/)
```





