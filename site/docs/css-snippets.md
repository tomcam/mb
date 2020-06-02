# CSS Snippets

article > h3:before { content: 'foo'; }


## Search

Add to nav.md:

* [News](/)
* [Specials](/)
* [Privacy](/)
* [Contact](/)
* [![Twitter logo](twitter-32x32-white.svg)]({{ .Site.Social.Twitter }}) [![Facebook logo](facebook-32x32-white.svg)]({{ .Site.Social.Facebook }}) [![LinkedIn logo](linkedin-32x32-white.svg)]({{ .Site.Social.LinkedIn }}) [![YouTube logo](youtube-32x32-white.svg)]({{ .Site.Social.YouTube}})
* <form style="display:inline;float:right;margin-right:1em;"action="https://www.google.com/search" method="GET"><span>&#9906;</span><input type="text" name="q" placeholder="Search"></form>


Currently working on a version that adds site:{{hostname}} automatically

https://www.google.com/search?site%3Atom.im+&q=apps
https://www.google.com/search?&q=apps+site%3Atom.im

Faling:

* <form style="display:inline;float:right;margin-right:1em;"action="https://www.google.com/search" method="GET"><span>&#9906;</span><input type="text" name="{{- site -}}{{-  %3A  -}}{{hostname}}{{- + -}}q" placeholder="Search"></form>




* <form style="display:inline;float:right;margin-right:1em;"action="https://www.google.com/search" method="GET"><span>&#9906;</span><input type="text" name="q" placeholder="Search"></form>


Adapted from: https://stackoverflow.com/questions/12036038/is-there-unicode-glyph-symbol-to-represent-search


## Beautify gradient

Adapted Fromt https://stackoverflow.com/questions/45763742/css3-swirl-gradient

My versionjlkji

First add this to the top of your light theme-light.css or theme-dark.css:
```
/* Default (paragraph) text is used for main branding in header */
:root
{
  --gradient-dark:#201D33;
  --gradient-light:#FF7700;
}

```
Then:

```
header {
  background: linear-gradient(to right, var(--gradient-dark) , var(--gradient-light));
  overflow: hidden;
}

header:before, header:after {
  content: "";
  display: block;
  position: relative;
  iborder-radius: 50% 50% 50% 50% / 60% 60% 40% 40%;
  border-radius: l0% 20% 30% 40% / 60% 50% 40% 10%;

}

header:before {
  box-shadow: -3rem 0 4rem -1rem var(--gradient-light);
  transform: skewX(-10deg);
}

header:after {
  transform: skewX(-10deg);
  box-shadow: inset -4rem 0 6rem -3rem var(--gradient-light);
}

```



### Original:

```
header {
  background: linear-gradient(to right, #50bcf3 , #60ec94);
  overflow: hidden;
}

header:before, header:after {
  content: "";
  display: block;
  position: relative;
  border-radius: 50% 50% 50% 50% / 60% 60% 40% 40%;

}

header:before {
  box-shadow: -30px 0 40px -10px #60ec94;
  transform: skewX(-10deg);
}

header:after {
  transform: skewX(-10deg);
  box-shadow: inset -40px 0 70px -30px #60ec94;
}
```


## The Debut theme header has a split line unordered list

It's pretty good, and includes an SVG graphic in place of a logo


## Tilted logo thingie

###

```
/*
 * --------------------------------------------------
 * Special feature: Unicode character rotated at 45
 * degrees in its own color (if desired). You can
 * eliminate this by deleting the content item
 * between single quotes.
 * Thank you, methodofaction! 
 * https://stackoverflow.com/questions/9779919/css-rotate-a-pseudo-after-or-before-content 
 * --------------------------------------------------
 */

/* ⣠ ▦ ␥ ▛ ◍ ☶ ⚙ */
header > ul:before  {
 	content:'⿹';
  color:red;
  transform:rotate(30deg);
	text-decoration:none;
	font-size:var(--branding-height);
	padding-bottom:6rem;
	margin-right:1em;
	display:inline-block;
}

```



### Other try:

```
/* ⣠ ▦▦  */
header > ul:before  {
	content:'⣠';
	text-decoration:none;
	font-size:8rem;
	padding-bottom:4rem;
	margin-right:0em;
	display:inline-block;
	-moz-transform: rotate(215deg);
	-o-transform: rotate(215deg);
	-ms-transform: rotate(215deg);
	transform:rotate(135deg);
	display:inline-block;
}

header > ul:before  {
	content:'▦';
	text-decoration:none;
	transform:rotate(45deg);
	margin-right:.25rem;
	font-size:3rem;
	display:inline-block;
	-moz-transform: rotate(45deg);
	-o-transform: rotate(45deg);
	-ms-transform: rotate(45deg);
	 transform: rotate(45deg);
}

```

## Drop shadow on img for screenshots

Give screenshots a drop shadow

```
article > img, article > p > img {
  width:80%;
  box-shadow: rgb(128,128,128) 1px 1px 3px 0px;
  display:block;
}
```
## Getting close to the right custom checkmarks

```
article > ul > li {list-style:none;}
article > ul > li::before {content: '';color:var(--article-bg);}
article > ul > li > input[type=checkbox]{position: relative;cursor: pointer;}
article > ul > li > input[type="checkbox"] ~ LABEL {color:var(--article-bg);}
/* Checkbox itself. This appears when either checked or unchecked. */
/* Paints over the existing text with background color. Ugh. */
article > ul > li > input[type=checkbox]:before {
         content: '';
         display: block;
         position: absolute;
         width: 1.4em;
         height: 1.5em;
         top: -.4em;
         left: 0;
         background-color:var(--article-bg); /* xxx */
}

/* Style a checked item. */
article > ul > li > input[type=checkbox]:checked:after {
         content: '\2714';
         display: block;
	 margin-left:2rem;
         width: 1em;
         height: 2em;
         position: absolute;
         top: -.5em;
         left: -3.5em;
}

/* Style an unchecked item. */
article > ul > li > input[type=checkbox]:not(:checked):after {
         content: '\25EF';
         display: block;
	 margin-left:2rem;
         width: 1em;
         height: 2em;
         position: absolute;
         top: -.5em;
         left: -3.5em;
}
```



## Gallery idea

Combine an h3 followed by a bullet list to get a gallery 3 elements wide.
This is from the Journey Home pagetype. I may have done a more advanced version so
check it.

### Three items
* Title ![stuff](image.png)
* Title ![stuff](image.png)
* Title ![stuff](image.png)


```
article > h3 + ul {
	clean:left;
	padding:0;
	margin:0;
	list-style-type:none;
}
article > h3 + ul li {
	/* There are 3 "cells" across the width of the container */
	max-width:33.3%;
	width:33.3%;
	padding:0;
	margin:0;
	/* Put the cells next to each other */
	float:left;
}

/* 3rd cell in the row starts a new line */
article > h3 + ul li:nth-child(3) {
	display:block;
}
article > h3 + ul li > img {
	/* Leave some white space to the right of the image */
	width:90%;
	/* Prevent text from crawling into the white space */
	margin-right:10%;
}


article > h3 + ul li > img + strong { color:blue;}


```

## Commenting

Pillar.css has the best comments

Also see boilerplate.txt

## Recent bug fixes


Default experiment with new theme foo
themename.css:

* Delete this line:

article {width:var(--max-article-width);}  

* Add:

article {max-width:var(--full-page);}

* Add:

article > p,                                                                                        
  article > ul,                                                                                     
  article > ol,                                                                                     
  article > li,                                                                                     
  article > dd,                                                                                     
  article > dl,                                                                                     
  article > hr {                                                                                    
    font-size:1.5rem;                                                                               
    line-height:2.25rem;                                                                            
    margin-bottom:2em;                                                                              
    max-width:var(--max-article-width);                                                             
  }     

article > p > code, article > code, article > pre, article > pre > code {                           
    font-family:var(--code);                                                                        
    font-size:1em;                                                                                  
    overflow:auto;                                                                                  
    margin-bottom:1em;                                                                              
    line-height:1.5em;                                                                              
    max-width:var(--max-article-width);                                                             
    }    

* The following worked for Pillar-style themes:

skin.css
    --max-article-width: 90%;                                                                       
    --max-aside-width: 20%;                                                                         
    --full-page:60%;                                                                                
    --left-gap:15%;                                                                                 
    --right-gap:15%;                                                                                
    --left-start:0%;                                                                                
                                                   
sidebar-left.css:
article, aside {float:right}
article {margin-left:0;}      

sidebar-right.css:
  article, aside {float:left}
  article {clear:left;}
	aside {margin-left:0;padding-left:0;} 


## Style link like a button. Tken from Journey home pagetype

```
/* Style link like button. Expands when mouse hovers over it. */
article > ul > li > p > a,
article > ul > li > p > a:link,
article > ul > li > p > a:visited
	{font-size:.65em;color:white;background-color:black;padding:.25rem .5rem .25rem .5rem;border-radius:.25em;text-decoration:none;}
article > ul > li > p > a:hover,
article > ul > li > p > a:active
	{padding:.5rem .75rem .5rem .75rem;}
```

## Responsive

This often helps:

```
article{width:100%;max-width:100%;}
```

### Tables

I think many themes could use this:

```
article >  table {width:var(--max-article-width);
```

## Experimental

Suppose you want the subhead over the H1 effect.
If you didn't want to use up a header style, you could do this:

``` css
header + nav + article > p:first-child {color:purple;font-size:4rem;}
```


## Getting the CSS right for tom.im theme

THE KEY WAS REMOVING var(--max-article-width); from the article > p and > code
declarations in themename.css

article > pre, article > p, article > ul, article > ol, article > li {line-height:1.5em; max-width:var(--max-article-width)}
article > pre > code {font-family:var(--code);max-width:var(--max-article-width);overflow:auto;}



article > pre, article > p, article > ul, article > ol, article > li, article > dd, article > dl {
	font-family:var(--formal);
	font-size:1.5rem;
	line-height:2.25rem;
	margin-bottom:2em;
	/* REMOVE!!!! */max-width:var(--max-article-width)

C}

### skin.css

```
	--max-article-width: 60%; 
	--max-aside-width: 20%;   
	--full-page:60%; 
	--left-gap:10%;
	--right-gap:0%;
	--left-start:0%;
```

### sidebar-left.css
```
	--max-article-width:70%;
	header, nav, footer {padding-left:5%;}
	article {padding-left:0;margin-left:0;margin-right:0;}
	aside{margin-right:0;}
```

### sidebar-right.css

```
	aside {margin-left:0;padding-left:0;padding-right:0%;}
```

## Simple footer for eprsona blog

```
![Landsteading logo](landsteading-fake-80x40.png)
Copyright 2020 Your name &middot; 12345 Main St. &middot; Springfield &middot; USA &middot; 12345 &middot; +1 (555) 1212 
```

## Full width tables for articles

```
article >  table {width:var(--max-article-width);margin-top:2rem;margin-bottom:3rem;margin-right:var(--right-gap);}
article > table > thead > tr > th {color:var(--header-bg);background-color:var(--header-fg);}
article > table > td,th {padding:1rem;}
article > table > tbody > tr > td {padding:1rem;border-bottom:.1px solid gray;}
```

## Using unordered lists on a row of links

You would precede this with the container tag, for example:

```
nav > ul > li > a:hover {etc}
```

This is a common convention that restyles the links so they look consistent 
except when hovered over or active, since usability guidelines require visual
feedback in those cases.

```
/* REPLACE nav WITH CONTAINING ELEMENT  and --nav-fg with appropriate color */
/* NAV */
nav > img {display:inline;}
nav > ul > li {display:inline;list-style-type:none;} 
nav > ul > li > a,
nav > ul > li > a:link,
nav > ul > li > a:visited {display:inline;text-decoration:none;color:var(--nav-fg);} 
/* For accessibility purposes show visual distinction when hovering or active */
nav > ul > li > a:hover,
nav > ul > li > a:active {text-decoration:underline;} 

```



This gives you a couple of options.

## Simplest: A single big branded item. Just put it in the header.md file like this:

```
My Awesome Product
```

This gives it a big branded presence.

``` css

/* HEADER */
/* Simplest case: Default (paragraph) text is used for main branding in header */
header > p {color:var(--header-fg);background-color:var(--header-bg);font-size:var(--branding-height);font-weight:700;display:inline;} 
```

## More flexible: Several items, the first one being branded


If you have more than one item, use an unordered list. These styles show the first item with branding applied,
then most of the remaining links showing slightly larger than other text, and the last link
styled like a button.

Your Markdown would look like this:

```
* [Pillar](/)  
* [Home](/)  
* [About](/)  
* [FREE TRIAL](/)  
```

And the CSS for it is:

```
header > ul {line-height:var(--header-height);}
header > ul {margin-left:0;padding-left:0;padding-right:1rem;list-style-type:none;display:inline;}
header > ul > li{display:inline;} 
header > ul > li > a {text-decoration:none;color:var(--header-fg);margin-right:2rem;} 
header > ul > li > a:hover {text-decoration:underline;} 
header > ul > li:first-child > a {color:var(--header-fg);font-size:var(--branding-height);font-weight:700;padding-right:2em;}
```

## Using a single header style to convey multiple kinds of information

## TODO: Obviously needs images

(Taken from Marlow)

```
#### **NEWS** | *Tom Campbell*  ~~{{ ftime "January 02, 2006" }}~~
```

``` cs
/* h4 Normal text is actually meant only for separator characters */
article > h4 {margin-top:1rem;color:var(--trim-fg)}
/* h4 **bold** formatting is good for the section name */ 
article > h4 > strong {color:black;}
/* h4 *em* formatting is bold but stil lighter; good for author credit */
article > h4 > em {color:gray;font-style:normal;font-weight:bold;font-size:1.2em;}
/* h4 ~~strikeout~~ formatting is abused for a dateline, and pulled to the right */
article > h4 > del {float:right;margin-right:4em;text-decoration:none;color:var(--trim-fg);font-size:1.2em;}
```


## Initial in a big box


This abuses the `~~` markup for strikethrought to put a big fat box around the text:

```
header p > del {font-size:var(--branding-height);font-weight:900;text-decoration:none;color:var(--header-fg);margin:1em;padding:.25em;border:1rem solid var(--header-fg);} 
```


### Center the footer, header, etc.

This example also centers any table inside the footer

``` css
footer, footer > table {padding-left:0;padding-right:0;width:100%;margin-left:auto;margin-right:auto;}
```

Or in some cases:

``` css
footer {padding-left:0;padding-right:0;text-align:center;width:100%;margin-left:auto;margin-right:auto;}

```

### Clear and simple code display


``` css
article > p > code {letter-spacing: -.05em;padding: .1em .3em;line-height: 1.4;text-shadow: 0 1px 0 #FFF;background-color: white;border: 1px solid darkgray;border-radius: 3px;box-shadow: 0 1px 0 rgba(20, 20, 20), 0 0 0 .1em #FFF inset;white-space: nowrap;}
article > pre > code, p > code {font-family:'Lucida Console','Consolas','Menlo','Monaco',monospace,monospace;font-size:1em;}
article > pre  {border: 1px solid lightgray;padding:1rem;background-color:ivory;}

```
<code style= "font-family:'Lucida Console','Consolas','Menlo','Monaco',monospace,monospace;font-size:1em;letter-spacing: -.05em;padding: .1em .3em;line-height: 1.4;text-shadow: 0 1px 0 #FFF;background-color: white;border: 1px solid darkgray;border-radius: 3px;box-shadow: 0 1px 0 rgba(20, 20, 20), 0 0 0 .1em #FFF inset;white-space: nowrap;">hello, world.</code>

