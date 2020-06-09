# Internal: Adding a theme to the live gallery

The live gallery demos show off all the main features of a Metabuzz theme, and they're generated automatically from text and images you supply. This document demonstrates that process.

## TODO: 


* I removed demo directory
* The old demo/s index.md is now called light.md
* Remove index.sidebar? I don't think it's being used.
* Avoid copying graphics needed only by the wide theme. Or delete them, whatever. Could be as easy as naming them wide_foo.png or whatever, then deleting from tht file pattern.

The live gallery generates a live gallery entry for each theme showing:

* A screenshot of the theme
* A live demo link to the light mode theme
* A live demo link to the dark mode theme
* And for each mode (dark and light), there are versions of live demo links to:
  - The theme without a sidebar
  - The theme with a left sidebar
  - The theme with a right sidebar

It does this using consistent names for files and front matter settings for those files internally. The `addtheme` utility copies the directory for the `wide` theme and tweaks the front matter of the files to point to the new theme.

The checklist below shows what you need to create for each new theme you add to the gallery.

This document assumes you're in the `website/gal` directory, which is the root directory for the live gallery. When you see a reference to the `demo` directory it's that subdirectory for the theme. For example, with the `wide` theme you'd find it in the `wide/demo` directory. 

## How the gallery works

The gallery entry shows is built from these parts:
1. A typical but abbreviated version of the theme page with some actual content, typically an image link and maybe a photo credit. That comes from `demo/intro.md`. 
1. Following contents the contents `demo/intro.md` is a description of the theme in a paragraph or less, representing the contents of the `demo/description.md` file.
1. The description is followed by a demo of all supported Markdown features automatically generated from a file in the `.mb/common` directory.

## How to create a live gallery demo for a new theme

To add a live demo to the gallery you run the local addtheme utility, then replace files and graphics as discussed below.

### Run `addtheme`

The first thing to do is run `addtheme` followed by the name of the theme.


* From the `gal` directory, run `addtheme` with the name of the theme following. In this case, it's the `genuine` theme:

```
./addtheme genuine
```

## Rebuild the website

The entire source for the Metabuzz is included. You can see your work by rebuilding it:

```
cd ../website
mb ../build
```

Then view it:

```
:: Windows users can skip the open 
open index.html
```

However that version recycles content unique to the `wide` theme. The next order of business is to replace all that.

## Checklist: What's needed to add a theme to the gallery


### Write a description and put in demo/description.md

* **Replace demo/description.md:** Put a compelling one-paragraph description in Markdown.

### Add images to the theme's root directory
* **Create a 1280x1024 screenshot of the theme.** Either light or dark mode is fine. Use the one you think makes your theme look best.


### Edit index.md in the theme's root directory

* **Edit `index.md` in the theme's root directory to include a one-sentence description of the theme:** Feel free to steal it from the `demo/description.md` file you just created. 

* You'd first edit this something like these lines:

```
![Screen shot of Future theme](theme-future-left-1280x1024.png)
```

* Then you'd replace the description text.

```
The Future theme is an eye-catching and simple [pillar-style](../pillar/index.html) theme, drawing much of its impact from the striking [Nova Mono](https://fonts.google.com/specimen/Nova+Mono) font by Wojciech Kalinowski.
```

Your theme doesn't have to support sidebars. Or it can support just one side but not the other. 

* Depending on how your theme supports sidebars, remove any of the text below from your theme's root directory `index.md`:

```
### Sidebar support
  Light theme: [Left](demo/light-sidebar-left.html) [Right](demo/light-sidebar-right.html)  
  Dark theme: [Left](demo/dark-sidebar-left.html) [Right](demo/dark-sidebar-right.html) 
```


### Add images for each mode/sidebar combination supported

Create the following images and name the exactly as shown, then put them in the demo root directory.

| Name of screenshot file  |  What it shows       | Dimensions   |
| ------------------------ |:--------------------:| --------------------------:|
| theme-dark-1024x1280.png             | Row 1, Col 2         | Row 1, Col 3               |
| Row 2, Col 1             | Row 2, Col 2         | Row 2, Col 3               |



### Replace the following Markdown files in the demo directory

* **Replace demo/intro.md:** `demo/intro.md` contains the first part of the article people see when they view the live demo theme. Add text and image links to the file as usual. Put images in the root directory of the demo theme.
* **Replace demo/left-sidebar-example.md:**  `demo/left-sidebar-example.md` contains the material you'd normally put in `sidebar.md`. This lets you create example text without affecting the production theme's `sidebar.md`. Put images in the root directory of the demo theme. 
* **Replace demo/right-sidebar.example.md:** The same goes for `demo/right-sidebar.example.md:`. Depending on how you use the theme `demo/right-sidebar.example.md` and `demo/left-sidebar.example.md` could be identical, as is the example with the Journey theme demo. Or they could be different, as is the example with the Reference theme demo. Put images in the root directory of the demo theme. 


## How a 
