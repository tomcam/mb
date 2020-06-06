# The Site Directory

When you create a site with Metabuzz using the `mb new site [sitename]` command, Metabuzz creates a subdirectory is created with whatever name you used in place of `[sitename]`. It populates that directory with:

* A stub `index.md` Markdown home page
* Copies of all themes Metabuzz ships with in a directory named `.themes`
* An empty `.common` directory you can use for boilerplate material 
* A `.headers` directory with HTML you might want to inject into the `<head>` section of web pages your site creates
* An empty `.pub` directory that gets emptied and repopulated with rendered HTML each time you [build your site](tutorial01.html#building-your-site)
* An `.scodes` directory containg HTML and template fragments used for shortcode purposes, for example, embedding YouTube videos
* A `.site` directory with configuration data used by your site as a whole, for example, 


## A quick site tour


### Create a site

* First create the site, here called `test`.

```
mb new site test
```

* Go to the new site's directory.

```
cd test
```

* Look at the stub `index.md` file:

##### file: index.md

```
# test

Welcome to test
```

### Examine the site's contents

* Now look at the site directory. It will have contents something like this.

```
index.md
.common
.headers
.site
.pub
.scodes
.themes
```

* If you check ouit the `.pub` directory you'll see it's empty.

### Build the site and look at the output

* Now create the site:

```
mb build
```

* And take a look at the `.pub` directory, also known as the Publish directory:

```
index.html
facebook-32x32-white.svg
linkedin-32x32-white.svg
twitter-32x32-white.svg
youtube-32x32-white.svg
themes/

```

* Here's what portions of `index.html` look like, with the rest replaced by `...` for brevity:

##### file: index.html

```
<DOCTYPE html>
<html lang=en>
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">

<title>metabuzz: Title needed here, squib</title>

<meta name="description" content="test">
<link rel="stylesheet" href="themes/wide/reset.css">
<link rel="stylesheet" href="themes/wide/fonts.css">
<link rel="stylesheet" href="themes/wide/layout.css">
<link rel="stylesheet" href="themes/wide/sizes.css">
<link rel="stylesheet" href="themes/wide/theme-light.css">
<link rel="stylesheet" href="themes/wide/wide.css">
<link rel="stylesheet" href="themes/wide/responsive.css">

...

<header><ul>
<li><a href="/">Default</a></li>
<li><a href="/">Create</a></li>
<li><a href="/">Pricing</a></li>
<li><a href="/">Try it Free</a></li>
</ul>
</header>

<nav><ul>
<li><a href="/">News</a></li>
<li><a href="/">Specials</a></li>
<li><a href="/">Privacy</a></li>
<li><a href="/">Contact</a></li>
<li><a href="">
  <img src="twitter-32x32-white.svg" alt="Twitter logo" /></a> 
  <a href=""><img src="facebook-32x32-white.svg" alt="Facebook logo" />
...
</nav>

<article><h1 id="test">test</h1>
<p>Welcome to test</p>
</article>

<footer><h2 id="wide"><a href="/">wide</a></h2>
<table>
<thead>
<tr>
<th>SECTIONS</th>
...
</tr>
</tbody>
</table>
</footer>

</body>
</html>

```

* Look at the `.pub/themes` subdirectory.

```
wide/
```

Wide is the theme Metabuzz uses if you don't specify a different one in the [front matter](front-matter.html) or the [site file](site-file.html).

If you had other pages using different themes they'd be copied under the `.pub/themes` directory as well.

* Finally, inspect the contents of the `.pub/themes/wide` directory:

```
wide.css
fonts.css
layout.css
reset.css
responsive.css
sizes.css
theme-light.css
```


## Contents of the site directory


test
├── .common
├── .headers
├── .pub
│   ├── facebook-32x32-white.svg
│   ├── index.html
│   ├── linkedin-32x32-white.svg
│   ├── themes
│   │   └── wide
│   │       ├── fonts.css
│   │       ├── layout.css
│   │       ├── reset.css
│   │       ├── responsive.css
│   │       ├── sizes.css
│   │       ├── theme-light.css
│   │       └── wide.css
│   ├── twitter-32x32-white.svg
│   └── youtube-32x32-white.svg
├── .scodes
├── .site
│   └── site.toml
├── .themes
│   ├── aventurine
│   │   ├── aventurine.css
│   │   ├── aventurine.toml
│   │   ├── fonts.css
│   │   ├── footer.md
│   │   ├── header.md
│   │   ├── layout.css
│   │   ├── nav.md
│   │   ├── reset.css
│   │   ├── responsive.css
│   │   ├── sidebar-left.css
│   │   ├── sidebar-right.css
│   │   ├── sidebar.md
│   │   ├── sizes.css
│   │   └── theme-light.css
│   ├── cv
│   │   ├── aside.html
│   │   ├── aside.md
│   │   ├── cv.toml
│   │   ├── footer.html
│   │   ├── footer.md
│   │   ├── header.html
│   │   ├── header.md
│   │   ├── layout.css
│   │   ├── nav.html
│   │   ├── nav.md
│   │   └── responsive.css
│   └── wide 
│      ├── .docs
│      │   └── howto.md
│      ├── facebook-32x32-white.svg
│      ├── fonts.css
│      ├── footer.md
│      ├── header.md
│      ├── home
│      │   ├── facebook-32x32-white.svg
│      │   ├── fonts.css
│      │   ├── footer.md
│      │   ├── header.md
│      │   ├── home.css
│      │   ├── home.toml
│      │   ├── layout.css
│      │   ├── linkedin-32x32-white.svg
│      │   ├── nav.md
│      │   ├── reset.css
│      │   ├── responsive.css
│      │   ├── sidebar-left.css
│      │   ├── sidebar-right.css
│      │   ├── sidebar.md
│      │   ├── sizes.css
│      │   ├── theme-dark.css
│      │   ├── theme-light.css
│      │   ├── twitter-32x32-white.svg
│      │   └── youtube-32x32-white.svg
│      ├── layout.css
│      ├── linkedin-32x32-white.svg
│      ├── nav.md
│      ├── reset.css
│      ├── responsive.css
│      ├── sidebar-left.css
│      ├── sidebar-right.css
│      ├── sidebar.md
│      ├── sizes.css
│      ├── theme-dark.css
│      ├── theme-light.css
│      ├── twitter-32x32-white.svg
│      ├── wide.css
│      ├── wide.toml
│      └── youtube-32x32-white.svg
│      ├── blog
│      │   ├── blog.toml
│      │   ├── facebook-32x32-white.svg
│      │   ├── fonts.css
│      │   ├── footer.md
│      │   ├── header.md
│      │   ├── layout.css
│      │   ├── linkedin-32x32-white.svg
│      │   ├── nav.md
│      │   ├── pillar.toml
│      │   ├── reset.css
│      │   ├── sidebar-left.css
│      │   ├── sidebar-right.css
│      │   ├── sidebar.md
│      │   ├── skin.css
│      │   ├── twitter-32x32-white.svg
│      │   ├── youtube-32x32-white.svg
│      │   └── zack.css
│      ├── facebook-32x32-white.svg
│      ├── fonts.css
│      ├── footer.md
│      ├── header.md
│      ├── layout.css
│      ├── linkedin-32x32-white.svg
│      ├── nav.md
│      ├── pillar.toml
│      ├── reset.css
│      ├── responsive.css
│      ├── sidebar-left.css
│      ├── sidebar-right.css
│      ├── sidebar.md
│      ├── theme-light.css
│      ├── twitter-32x32-white.svg
│      ├── youtube-32x32-white.svg
│      ├── zack.css
│      └── zack.toml
└── index.md

