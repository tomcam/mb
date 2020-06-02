## 0103: "Error reading front matter"
Probably the result of a malformed key/value pair in the front matter. Real-life example: the
missing equals sign in this front matter. It read like this:

```
===
title"Debut Theme"
===
```

But there should have been an equals sign separating the key
from the value, like this:

```
===
title = "Debut Theme"
===
```

Another example is a missing double-quote mark, for example:

```
===
title = Debut Theme"
===
```

It's hard to notice, but there should be a double-quote in front of the `D` as shown:

```
===
title = "Debut Theme"
===
```


0108: "Error reading theme file"

This can happen when there's a syntax error in the .toml file, for example, the missing close quote on `"responsive.css"` 

```
Stylesheets = ["reset.css", "fonts.css", "layout.css", "theme-light.css", "home.css", "sizes.css", "responsive.css]
```
 Error loading /Users/tom/code/d2/.d2/themes/journey/home/home.toml
