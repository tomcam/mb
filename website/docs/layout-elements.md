# Layout elements

Metabuzz follows 
[standard layout conventions for HTML content](https://developer.mozilla.org/en-US/docs/Learn/HTML/Introduction_to_HTML/Document_and_website_structure#HTML_layout_elements_in_more_detail) 
as described by MDN, and as implied by the format of an HTML document.
The standard format has a site broken up into sections, based on the HTML tags
`<header>`, `<nav>`, `<aside>`, `<article>`, and `<footer>`.

These sections are called *layout elements* because they directly affect
both the physical look of the page and the semantic menaing of each
layout element.

| Layout element  | Theme TOML  | Function   |
| :-------------- | :--------   |:-----------|
| `<header>`      | `[header]`  | Header     |
| `<nav>`         | `[nav]`     | Navbar     |
| `<aside>`       | `[sidebar]` | One sidebar (an HTML document can only have 1 `article` tag) |
| `<article>`    | `[article]`  | Body of document |
| `<footer>`     | `[footer]`  | Footer     |

The layout elements are generated from descriptions in your `[theme.toml](#theme-toml)`
file.

For example, you may see a section like this in a `theme.toml` file.

##### file mytheme.toml
```
[Sidebar]
  HTML = ""
  File = "sidebar.md"
```

This directs Metabuzz to create an `<aside>` tag with the HTML generated from the
file `sidebar.md` and using that as its contents (inner HTML). (You can use
any filename, by the way. It could be `aside.md` or `foo.md`.)

Here's a simple example.

##### file sidebar.md
```
**Sale** in the [shop](shop.html) today only!
```

#### Resuling HTML
```html
<aside><p><strong>Sale</strong> in the <a href="shop.html">shop</a> today only!</p>          
</aside>  
```


## Layout elements defined in the theme TOML file

The relevant parts of your theme TOML file for 
creating web content look like this:

```
[header]
  html = ""
  file = "header.md"

[nav] 
  html = ""
  file = "nav.md"

[article]
  html = ""
  file = ""

[sidebar]
  html = ""
  file = "sidebar.md"

[footer]
  html = ""
  file = "footer.md"
```

### file can be either HTML or Markdown

While most of the examples shown use Markdown files
