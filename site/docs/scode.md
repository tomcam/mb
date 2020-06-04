===
templates="off"
===
# Scode

## TODO: See Zola's

https://www.getzola.org/documentation/content/shortcodes/

Inserts an the HTML file and possible parameters using a map, which must
have an attribute named `filename=` followed by the name of an HTML file.
The other attributes are used by the HTML file as appropriate. In this
example a YouTube instructional video on static site generators is inserted 
into the rendered markdown:


```
youtube = { filename="youtube.html", id = "dQw4w9WgXcQ" }
```

Given this file named `youtube.html` in the `.scodes` directory:

    <div>
    <iframe id="ytplayer" type="text/html" width="640" height="360"
      src="https://www.youtube.com/embed/{{- index "id" -}}?autoplay=0"
      frameborder="0" allow="autoplay; encrypted-media" allowfullscreen></iframe>
    </div>

And this in the front matter:

```
===
[List]
youtube = { filename="youtube.html", id = "dQw4w9WgXcQ", other = "other here" }
===
```


And this in the markdown:

```
{{"{{"}} scode .FrontMatter.List.youtube {{"}}"}}

```



