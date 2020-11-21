# The .mb directory

# TODO: Its location can be changed

To build a site, your project directory needs:
* A directory named `.mb` at its root.
* A file named `index.md` or `README.md` 


Note: In some operating systems, such as MacOS or Linux.
directories starting with a dot character are  invisible
by default, but you can still view their contents or
make them your current directory. Metabuzz uses
this to reduce visual clutter.

The `.mb` directory contains things that may influence the
output of your site, such as theme files or [header tags](headtags.html)
for code injection into your output file's `header` section.

If you have multiple projects then each will have its own `.mb` directory.

`.mb` directory's location is in your application data directory,
which varies by operaing system You can choose where it goes with the
[Global Configuration File](config-file.html).


To find out where your `.mb` file is currently, just enter at the
terminal prompt:

```
mb info
```

And you'll get output something like this:

```
Home dir: /Users/layla
Reported current dir: /Users/layla/code/mb/foo
Actual current dir: /Users/layla/code/mb/foo
scode path /Users/layla/code/mb/foo/.mb/scodes: (present)
a.Flags.Verbose false
Default config directory /Users/layla/Library/Application Support/metabuzz/.mb: (Not present)
Actual config directory /Users/layla/code/mb/.mb: (present)
Site file:  /Users/layla/code/mb/foo/.mb/site/site.toml: (present)
Theme directory /Users/layla/code/mb/.mb/themes: (present)
Code highlighting style:  github
Default theme:  wide
Highlight: 
This appears to be a project/site source directory
Site directory:  /Users/layla/code/mb/foo: (present)
Publish directory /Users/layla/code/mb/foo/.pub: (present)
Theme directory /Users/layla/code/mb/foo/.mb/themes: (present)
Headers directory /Users/layla/code/mb/foo/.mb/headtags: (present)
Shortcode directory:  /Users/layla/code/mb/foo/.mb/scodes: (present)
```j

