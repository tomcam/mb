

You can specify a header logo using a URL, of course. 

```
[company]
name = "METABUZZ"
#HeaderLogo = "https://icongr.am/clarity/helix.svg?size=48&color=448FA3"
HeaderLogo = "metabuzz-red.svg"
URL = "https://metabuzz.com"
```

Vertical social strip with colors generated from .Site.Social. Not dynamic.

```
[![Twitter logo](twitter-blue-30x30.svg)]({{ .Site.Social.Twitter }})

[![Facebook logo](facebook-blue-30x30.svg)]({{ .Site.Social.Facebook }})

[![LinkedIn logo](linkedin-blue-30x30.svg)]({{ .Site.Social.LinkedIn }})

[![YouTube logo](youtube-red-30x30.svg)]({{ .Site.Social.YouTube}})
```


Horizontal social strip with colors generated from .Site.Social. Not dynamic.

```
[![Twitter logo](twitter-blue-30x30.svg)]({{ .Site.Social.Twitter }})
[![Facebook logo](facebook-blue-30x30.svg)]({{ .Site.Social.Facebook }})
[![LinkedIn logo](linkedin-blue-30x30.svg)]({{ .Site.Social.LinkedIn }})
[![YouTube logo](youtube-red-30x30.svg)]({{ .Site.Social.YouTube}})
```



Horizontal social strip generated from entries in .Site.Social


{{ if .Site.Social.Twitter }}[![Twitter logo](twitter-blue-30x30.svg)]({{ .Site.Social.Twitter }}){{ end -}} {{ if .Site.Social.Facebook }}[![Facebook logo](facebook-blue-30x30.svg)]({{ .Site.Social.Facebook }}){{ end -}}{{ if .Site.Social.LinkedIn }}[![LinkedIn logo](linkedin-blue-30x30.svg)]({{ .Site.Social.LinkedIn }}){{ end -}} {{ if .Site.Social.YouTube }}[![YouTube logo](youtube-red-30x30.svg)]({{ .Site.Social.YouTube}}){{ end -}} 



{{- /*  Automatically name first item in header    
        based on company name, author name name
        if no company was specified, or just 
        the name of the theme if neither of those
        was specified.
        
*/ -}}
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



# Horizontal social strip generated from entries in .Site.Social

{{ if .Site.Social.Twitter }}[![Twitter logo](twitter-blue-30x30.svg)]({{ .Site.Social.Twitter }}){{ end -}} {{ if .Site.Social.Facebook }}[![Facebook logo](facebook-blue-30x30.svg)]({{ .Site.Social.Facebook }}){{ end -}}{{ if .Site.Social.LinkedIn }}[![LinkedIn logo](linkedin-blue-30x30.svg)]({{ .Site.Social.LinkedIn }}){{ end -}} {{ if .Site.Social.YouTube }}[![YouTube logo](youtube-red-30x30.svg)]({{ .Site.Social.YouTube}}){{ end -}} 


## How to create a Download section


## Clickable link in Markdown

```
[![Twitter logo](twitter-32x32-black.png)](https://twitter.com)
```


### Clickable Social Media icons in their original colors

This also removes hardcoding


## TODO: This code needs proper escaping

```
* [![Twitter logo](twitter-24x24-blue.svg)]({{ .Site.Social.Twitter }}) 
* [![Facebook logo](facebook-24x24-blue.svg)]({{ .Site.Social.Facebook }})
* [![LinkedIn logo](linkedin-24x24-blue.svg)]({{ .Site.Social.LinkedIn }})
* [![Pinterest logo](pinterest-24x24-red.svg)]({{ .Site.Social.Pinterest }}) 
* [![Instagram logo](instagram-24x24-magenta.svg)]({{ .Site.Social.Instagram }}) 
```

They can be found in the `/docs/img` directory

 
## CSS to style a link like a button 

From Personality.css


## Display icons without icon files

[Icongram](https://icongr.am/) is a wonderful, free service
that hosts SVG icons so you don't need to copy them into your
site as assets, for personal or commercial use. 
And you can request them in any size or color.

* Go to [Icongram](https://icongr.am) and choose a collection
you like, for exmaple, [Simple Icons](https://icongr.am/simple).

* Look for one you like, say, `arduino` there on the top row. 

Note that it has a caption underneath, `arduino` in this case.

* Create a URL to request the icon by building the URL like this,
replacing the items in bold with their actual values. 
Don't worry, an example follows.

```
https://icongr.am/**collection**/**label**.svg?size=**size**&color=**colorvalue**
```

So in this case you'd replace **collection** with `simple`, 
and **label** with  `arduino`, and **size** with `64`,
meaning the generate image will be 64x64 pixels, and
finally **colorvalue** is replaced by `0000FF`, which
means pure blue. Try clicking this or pasting it into 
your browser bar:

[https://icongr.am/simple/arduino.svg?size=64&color=0000FF ](https://icongr.am/simple/arduino.svg?size=64&color=0000FF)

You could write Markdown for it, of course.

```
![Arduino icon](https://icongr.am/simple/arduino.svg?size=64&color=0000FF)
```

And the result is:

![Arduino icon](https://icongr.am/simple/arduino.svg?size=64&color=0000FF)


<a id="gallery"></a>
## Social horizontal strip: row of social media icons

This displays a row of clickable social media icons. Here's how to use it:

* Copy the contents of this file into your Markdown. It can be an article, 
header, footer, whatever.

##### filename: product-snippets-social-stripe.md 

product-snippets-social-stripe.md

* Append the folllowing to the bottom of your [site.toml](site-toml.html) file
under the `[List]` header. Don't repeat the `[List]` part if one's already there. 

Obviously you can omit accounts you don't have, or add ones that aren't here.

The icons can be found at [https://icongr.am/simple](https://icongr.am/simple)
if you're missing any.

```
[List]
Iconsize = "24"
Iconcolor = "FFFFFF"
Iconurl = "https://icongr.am/simple/" 
  [[List.gallery]]
    alt = "Facebook logo"
    filename = "facebook.svg"
    url = "facebook.com/metabuzz"

  [[List.gallery]]
    alt = "LinkedIn logo"
    filename = "linkedin.svg"
    url = "linkedin.com/metabuzz"

  [[List.gallery]]
    alt = "Twitter logo"
    filename = "twitter.svg"
    url = "twitter.com/metabuzz"

  [[List.gallery]]
    alt = "YouTube logo"
    filename = "youtube.svg"
    url = "youtube.com/metabuzz"

  [[List.gallery]]
    alt = "Instagram logo"
    filename = "instagram.svg"
    url = "instagram.com/metabuzz"

  [[List.gallery]]
    alt = "Pinterest logo"
    filename = "pinterest.svg"
    url = "pinterest.com/metabuzz"

```

* Make sure everywher you see `url = ` you replace the URL with your own
social media web address.

