# Creating a Metabuzz production site from scratch


* Copy the theme you like because it's likely you'll want to change the
headers or sidebar
* Fill out company, author in site.toml
* Fill out social media URLs in site.toml "[social]" section
* Add favicons or whatever to the Header directory


## Add header tags as appropriate

The [headers directory](headers-dir.html) 

### Example: adding a header tag for favicon support

```html
<link rel="icon" type="image/png" href="favicon-32x32.png" sizes="32x32" />
<link rel="icon" type="image/png" href="favicon-16x16.png" sizes="16x16" />
```

# For each page

* Try to have a unique Description
* May want to stick to a single H1
* Keep semantically important headers to H1-H3
