# How to create a gallery entry

A gallery entry is complex but it's actually easy to create. The main things you need to do are:

1. [Search and replace](#replace) the source files with your new theme name
2. [Create the description file](#description)
2. [Create the optional intro file](#introfile)

## How the gallery files work together

The gallery shows each theme's main features by populating a typical page, but since there are at least 16 combinations possible with just the dark, light, left sidebar, right sidebar, and no sidebar options they make extensive use  of including common material to avoid duplication.

### Light and dark copies of the home page

There are several copies of the "home" (`index.md`) page because it shows the light version and the dark version (named `index.md` and `dark.md`). To keep things consistenet these pages have no content of their own. Instead, they use the [inc](functions/inc.html) function to include the sample content (most of which is found in `description.md`.

`index.md` and `dark.md` differ only in their front matter. `index.md` specifies in the front matter that it uses the light theme. It's technically unnecessary because light theme is the default:

```
===
mode="light"
===
```

And of course `dark.md` specifies the opposite:
```
===
mode="dark"
===
```
### Right and left sidebar copies of the home page

There are other copies of the home page. The names of these example files are `light-sidebar-left.md` and `light-sidebar-right.md`. These are the home page in light mode with a left sidebar and a right sidebar. Since you can't have both on one page no fourth option is needed (the third option is no sidebar, which is the default).

Matching those to are `dark-sidebar-right.md` and `dark-sidebar-left.md`, the equivalent demo pages for the dark theme.


* Just once: Ensure the following files are in the project's .common directory:
  + [allfiles.md](https://github.com/tomcam/d2/blob/master/gal/.common/allfiles.md)
  + [kitchen.md](https://github.com/tomcam/d2/blob/master/gal/.common/kitchen.md)
  + [showtheme.md](https://github.com/tomcam/d2/blob/master/gal/.common/showtheme.md)
  + [variations.md](https://github.com/tomcam/d2/blob/master/gal/.common/variations.md)

## How to customize the source files

Quick take: the files in the .common directory don't have to be touched. Most of the other .md files only need the name of the theme changed.

### Replace theme name in source files with the new theme name. {#replace}

Many files (such as `index.md`) include front matter naming the theme (and often that's all that needs to be changed in those files). 

Those files are:

* `index.md`
* `dark.md`
* `dark-sidebar-left.md`
* `dark-sidebar-right.md`
* `light-sidebar-left.md`
* `light-sidebar-right.md`


The section to replace looks like this:

```
===
theme="CHANGEME"
===
```

* Replace "CHANGEME" with the name of your theme in each of those files. Do not forget the quotes. 

| File				| Purpose			                      | Links to, includes, or expects        |
| ----------------------------- | --------------------------------------------------- | --------------------------------------|
| `index.md` 			| Light theme home page with no sidebars              | allfiles.md			      |
| `dark.md` 			| Dark theme home page with no sidebars               | allfiles.md			      |
| `light-sidebar-right.md`  	| Light theme home page with right sidebar            | light-sidebar-right.sidebar           |
| `light-sidebar-right.sidebar` | Light theme home page right sidebar                 | right-sidebar-example.md              |
| `right-sidebar-example.md`	| Theme-independent right sidebar example markup      |					      |
| `left-sidebar-example.md`	| Theme-independent leftsidebar example markup        |					      |
| `dark-sidebar-left.md`  	| Dark theme home page with left sidebar              | dark-sidebar-left.sidebar	      |	
| `dark-sidebar-right.md`  	| Dark theme home page with right sidebar             | dark-sidebar-right.sidebar            |
| `dark-sidebar-right.sidebar`  | Dark theme home page right sidebar                  | right-sidebar-example.md              |
| `dark-sidebar-left.md`  	| Dark theme home page with left sidebar              | dark-sidebar-left.sidebar	      |	
| `description.md`              | Explains and illustrates main theme features        |                                       |

### Create the description file {#description}

The file `description.md` shows your theme's best features and sells your audience on why they should use this theme. It should have a screen shot at or near the top. Ideally it will have many other images showing features visually and how to use Markdown (or, if necessary, CSS) to use your theme's unique features to the fullest. 

### Create the (optional) intro file {#introfile}

The optional intro file contains whatever markup might normally start a page using this theme.

Here was one I used with the default theme:

```
![Picture of vintage-style portable phonograph](img-sample-portable-phono-640x417.jpg)

#### Photo by [Suzy Hazelwood](www.pexels.com/photo/black-speaker-beside-vinyl-sleeve-3077740/) from Pexels  

#### By **Tom Campbell** | Staff Writer

#### [![Twitter logo](twitter-white-24x24.svg)]({{"{{"}} .Site.Social.Twitter {{"}}"}}) [![Facebook logo](facebook-white-24x24.svg)]({{"{{"}} .Site.Social.Facebook {{"}}"}}) [![LinkedInlogo](linkedin-white-24x24.svg)]({{"{{"}} .Site.Social.LinkedIn {{"}}"}}) [![YouTube logo](youtube-white-24x24.svg)]({{"{{"}} .Site.Social.YouTube {{"}}"}})
```


## TODO: READ RANDOM Notes
`.sample` directory contains skeletons of everything needed for a gallery directory. Note that index.new is so named because thre may be an existing index file in the directory.
theme="CHANGEME"

# To replace the theme from CHANGEME to zz do this:
```
find . -type f -name "*.md" -exec sed -i -e 's/CHANGEME/zz/g' {} +
```

