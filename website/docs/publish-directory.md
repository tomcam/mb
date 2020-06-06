# The Publish directory

The publish directory contains your website: the set of HTML files, theme files, CSS, image, sound, and other assets that Metabuzz generates from your Markdown files and other assets. It is in the `.pub` subdirectory immediately off the root of your site directory. Ultimately it will be copied to the `WWW` or whatever directory your web host uses to publish HTML files from.

You are not expected to make changes to the publish directory manually, because each time you run `mb build` the first thing that happens is that its contents are deleted. Otherwise you wouldn't have a reliable way of reliably creating a site if you need to roll back a version. 

Let's go through the creation of a minimal site to see how the publish directory is populated. 

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

* If you check out the `.pub` directory you'll see it's empty.

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




