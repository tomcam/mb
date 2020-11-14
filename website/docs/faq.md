# Metabuzz Frequently Asked Questions (Metabuzz FAQ)


## Code listings and fenced code in Metabuzz

### How do I change highlighting styles in Metabuzz fenced code listings?

Metabuzz comes with dozens of presets for its
[Code highlighting](code-highlighting.html)
feature thanks to the [Chroma](https://github.com/alecthomas/chroma) package. You can see an example of each one in the 
[Style gallery](https://xyproto.github.io/splash/docs/)

You change the global highlighting style using the
[site file](site-file.html)
* To set the code highlighting style add the following
to your site file, which you can find in the project's 
root directory at `.mb/site/site.toml`:

```
# Replace "github" with any other style found at
at the [Style gallery](https://xyproto.github.io/splash/docs/)
k
jkjkjkl
jkjkPjklP jkl0kk
[[MarkdownOptions]
  HighlightStyle = "github"
  HeadingIDs = true

<Paste>MarkdownOptions]
  HighlightStyle = "github"
  HeadingIDs = true



https://xyproto.github.io/splash/docs/

Source
https://github.com/alecthomas/chroma/tree/master/styles

## Formatting Markdown

### How do I get two lines to display without a blank space between them?
How do I get these two lines to be displayed as shown? They end up all on line.

Normally if your Markdown looks like this:

```
Line 1
Line 2
```

The output is this:

Line 1 
Line 2

But if you put a blank space between the lines, the output leaves a blank line between them:

```
Line 1

Line 2
```

Line 1

Line 2

So how do you get them to do *this?*

Line1  
Line2

#### Answer 

Markdown has a special property: if you end a line with 2 spaces, the following line will appear under it. **In the following example, replace each dot character (·) with a space.**

```
Line 1··
Line 2 
```
And you will indeed get what you hoped for:

Line1  
Line2


