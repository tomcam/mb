===
theme="new-pillar"
sidebar="right"
mode="light"
===


# hello, world.

This is what the empty theme looks like.

{{ inc "theme-and-variations.md" }}


**Table of contents** 

* [Common text formatting](#common-text-formatting)
* [Links](#links)

## Markdown syntax

Here's how markdown appears in the **{{.FrontMatter.Theme }}** theme
{{- if .FrontMatter.PageType }}
with the PageType **{{ .FrontMatter.PageType }}**
{{ end }}:
## Common text formatting

#### You type:
```
Normal body text, **strong**, ~~strikethrough~~, and with *emphasis*.
```

#### It shows as:
Normal body text, **strong**, ~~strikethrough~~, and with *emphasis*.

Horizontal rule:

#### You type:
```
---
```

#### It shows as:
---

## Links

#### You type:
```
[link text](https://appscripting.com)
```

#### It shows as:
[link text](https://appscripting.com)

![screenshot](theme-1280x1024.png)

