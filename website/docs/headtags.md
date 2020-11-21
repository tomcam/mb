# The headtags directory and code injection

The `.mb/headtags` subdirectory inside the [.mb directory](mb-directory.html)
contains a list of text files that get copied
to your HTML output file in the form of `<meta>` tags. 
This sometimes goes by the impressive-sounding name of *code injection*.

For example, the file named `generator` contains these contents
(the version number may differ from this one).
```
<meta name="Generator" content="Metabuzz 1.0.1">
```

## Note

All tags in this directory get copied verbatim 
into each generated HTML, so be sure they're well-formed.


