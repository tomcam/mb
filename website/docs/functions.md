# Built-in Functions

## ftime

Displays many possible combinations of date and time.

Explain that example format is based exactly on: 

```
Mon Jan 2 15:04:05 -0700 MST 2006 
```

| Format string                                        | Result                                           |
|------------------------------------------------------------|--------------------------------------------|
| \{\{ftime "Mon Jan 2 15:04:05 -0700 MST 2006"\}\}    | {{ ftime "Mon Jan 2 15:04:05 -0700 MST 2006" }}  |
| \{\{ftime "3pm"\}\}    | {{ ftime "3pm" }}  |
| \{\{ftime "15"\}\}    | {{ ftime "15" }}  |
| \{\{ftime "3:04pm"\}\}    | {{ ftime "3:04pm" }}  |
| \{\{ftime "15:04"\}\}    | {{ ftime "15:04" }}  |
| \{\{ftime "3:04:05"\}\}    | {{ ftime "3:04:05" }}  |
| \{\{ftime "15:04:05"\}\}    | {{ ftime "15:04:05" }}  |
| \{\{ftime "3:04:05 -0700 MST 2006"\}\}    | {{ ftime "3:04:05 -0700 MST 2006" }}  |
| \{\{ftime "15:04:05 -0700 MST 2006"\}\}    | {{ ftime "15:04:05 -0700 MST 2006" }}  |
| \{\{ftime "Jan 2 15:04:05 -0700 MST 2006"\}\}    | {{ ftime "Jan 2 15:04:05 -0700 MST 2006" }}  |
| \{\{ftime "Jan 2 2006"\}\}    | {{ ftime "Jan 2 2006" }}  |
| \{\{ftime "Jan 2, 3:04pm"\}\}    | {{ ftime "Jan 2,  3:04pm" }}  |
| \{\{ftime "January 2006"\}\}    | {{ ftime "January 2006" }}  |
| \{\{ftime "Monday, January 2006"\}\}    | {{ ftime "Monday, January 2006" }}  |
| \{\{ftime "Monday, January 2"\}\}    | {{ ftime "Monday, January 2" }}  |
| \{\{ftime "Monday, January 2, 2006"\}\}    | {{ ftime "Monday, January 2, 2006" }}  |
| \{\{ftime "Jan"\}\}    | {{ ftime "Jan" }}  |
| \{\{ftime "January"\}\}    | {{ ftime "January" }}  |
| \{\{ftime "2"\}\}    | {{ ftime "2" }}  |
| \{\{ftime "2006"\}\}    | {{ ftime "2006" }}  |
| \{\{ftime "1/2/06"\}\}    | {{ ftime "1/2/06" }}  |
| \{\{ftime\}\}    | {{ ftime }}  |
| \{\{ftime "Mon Jan 2 15:04:05 -0700 MST 2006"\}\}    | {{ ftime "Mon Jan 2 15:04:05 -0700 MST 2006" }}  |

## hostname

Returns the name of the server it's running on (or URL).

## inc

Lets you include "boilerplate" text from a common directory in cases where you would otherwise need to copy and paste. 

Treats it as a Go template, so either HTML or Markdown
work fine.
The location of the file appears first, before a pipe character.
It can be one of:

"article" for the current markdown file's directory
"common" for the Site.CommonSubDir subdirectory

So it might look like inc "articles|kitchen.md"


See [inc](functions/inc.html) in the functions reference

"foo.md"
"article|foo.md"
"common|foo.md"


## scode

Can pass it the name of a markdown file too, not just HTML

## toc \{\{levels\}\} \{\list type\}\} 

Generates a table of contents at the current
location in the document. It's generated from
headers.

`{levels}` is the header depth. Defaults to 6 and can be
specified as "1" to "6".

`{list type}` is either "ul" or "ol". "ul" is an undordered
(bullet) list. "ol" is an ordered (numbered) list.

### Example: Generate a table of contents for all headings level 1 to 6:

```
## Table of contents

\{\{ toc \}\}

```

This is equivalent to:
```
\{\{ toc "6" "ul" \}\}

```


### Example: Generate a table of contents headings level 1 to 3:

```
## Table of contents

{{ toc "3" }}

```
### Example: Generate a table of contents headings level 1 to 3,
preceding the titles with numbers:

```
## Table of contents

{{ toc "6" "ol" }}

```





```
toc
```



