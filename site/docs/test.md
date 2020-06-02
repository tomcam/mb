===
sidebar="none"
list = [{theme = "ned"}]
===

# Theme is {{.FrontMatter.Theme}}

Sidebar? {{.FrontMatter.Sidebar}}

The time is {{ftime}}

The list is: {{.FrontMatter.List}}

[Right sidebar](test-rside.html)
[Left sidebar](test-lside.html)
[No sidebar](test.html)


The maximum width of text is of increasing concern as screens get to higher resolutions. While a phone might have a large number of dots horizontally, they're quite small.


## Let's take a look at the kitchen sink

Normal body text, with **strong**, and with *emphasis*.

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
|  Left-justified Header|  Centered header   | Right-justified    |
| --------------------- |:------------------:| ------------------:|
| Row 1, Col 1          | Row 1, Col 2       | Row 1, Col 3       |
| Row 2, Col 1          | Row 2, Col 2       | Row 2, Col 3       |

```

And here's what results in this theme:

|  Left-justified Header|  Centered header   | Right-justified    |
| --------------------- |:------------------:| ------------------:|
| Row 1, Col 1          | Row 1, Col 2       | Row 1, Col 3       |
| Row 2, Col 1          | Row 2, Col 2       | Row 2, Col 3       |

## Block quote

> "Drive, she said," went the block quote.



