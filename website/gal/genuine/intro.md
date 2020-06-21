## Abusing classless CSS for richer HTML presentation

CSS without classes sounds a bit horrifying, but it's fact
of life for any [Markdown/CommonMark](https://commonmark.org) converter.  Markdown has no provisions for classes, so your websites are doomed to a life of montony, right?

### Classless CSS Hack #1: page regions matter

The HTML standard (MDN's [Basic sections of a document](https://developer.mozilla.org/en-US/docs/Learn/HTML/Introduction_to_HTML/Document_and_website_structure#Basic_sections_of_a_document) describes an HTML page as broken up into these sections, called page regions:

* Header
* Navigation bar
* Main content
* Sidebar
* Footer

Which means that you can create completely different text using the same markup. For example,
your main content area (normally within `<article>` tags) might show text as black over whitesmoke:

```css
article > p {color:black;background-color:whitesmoke;}
```

Whereas your sidebar may show them reversed:

```css
aside > p {color:whitesmoke;background-color:black;}
```

Which means the text `hello, world.` would appear as dark on light in the article,
but light on dark in the sidebar, even though the Markdown was the same.

### Classless CSS Hack #2: push character formatting to its limits

You probably know that to create **strong** word in markdown you decorate with a pair of `**` 
asterisk characters on each side of the phrase to embolden. Suppose your logo has two colors. 
Neither is either in italic or bold.



