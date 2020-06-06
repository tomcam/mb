===
theme="lanyon"
sidebar="left"
===

* [Try the dark theme](index-dark.html)
* [Try the home page](home.html)
                                                                   
## Sidebar support

Sidebar support for this theme:

* [Left sidebar](index.html)

### Consider revamping your website.        

Or just replacing it completely with Debut. You know you like it.     


[Google Apps Script](https://appscripting.com) is a big subject with a long history and very few complete tutorials. I had a deceptively simple question that ended up as this article: what does a minimum, typical Google Apps Script program look like and what are current best practices? All the existing tutorials I know about leave a huge amount unsaid, and work best if you already know Apps Script.

Most Apps Script programs start with some kind of user input. That normally requires a custom form. Custom forms use HTML and are translated into Javascript objects in Google Apps, and they function very differently from the branded Google Forms. This app will appear as a sidebar on a spreadsheet. It gets built automatically by the built-in onOpen() function, whch is triggered when the spreadsheet is first loaded. When you choose it from the menu, the sidebar itself is created and displayed.

This tutorial shows you step by step how to create the simple activity tracker shown below, which lets you jot down something you did, then tiestamps it and appends both items to the end of a spreadsheet.

## Let's take a look at the kitchen sink

Normal body text, [link text](https://appscripting.com) with **strong**, and with *emphasis*.

Horizontal rule:

---

# Header style h1
## h2
### h3
#### h4
##### h5
###### h6

## Coding styles

You can format text inline as `code`, or go block style:

``` python
print ("This is a code block")
```

## There are 2 or 3 kinds of list types

### Ordered lists

1. Ordered lists have numeric sequences
1. Even though you write `1` in Markdown,
1. The numbers display properly on output


### Unordered, or bullet lists

Reasons people hate bullet lists

* It reminds them of bad PowerPoint
* I can't really think of another reason
  + You can indent bullet lists
  + Just go in a couple spaces, and use `+` instead of bullet
  + The `+` isn't required. It's just for clarity

### The "third" list type: definition lists

A definition list lets you display things like an item
and its meaning in a distinct way:

Definition list
: A way to show a visual relationship between a word or phrase
and an explanation for it

Markdown
: A convention for generating HTML from a more human-readable 
source format.

### Markdown trick: Creating clickable image links in Markdown
```
[![Twitter logo](twitter-32x32-black.png)](https://twitter.com)
```


## Tables

You can just sort of draw tables:

```
|  Left-justified Contents |  Centered Contents   | Right-justified Contents   |
| ------------------------ |:--------------------:| --------------------------:|
| Row 1, Col 1             | Row 1, Col 2         | Row 1, Col 3               |
| Row 2, Col 1             | Row 2, Col 2         | Row 2, Col 3               |

```

And here's what results in this theme:

|  Left-justified Contents |  Centered Contents   | Right-justified Contents   |
| ------------------------ |:--------------------:| --------------------------:|
| Row 1, Col 1             | Row 1, Col 2         | Row 1, Col 3               |
| Row 2, Col 1             | Row 2, Col 2         | Row 2, Col 3               |

## Block quote

> "Drive, she said," went the block quote.




