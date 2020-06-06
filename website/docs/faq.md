# Metabuzz Frequently Asked Questions (Metabuzz FAQ)

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


