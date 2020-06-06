# Using the command line

[Create a new project]  
[Build a site](#build-a-site)  
[Display configuration](#config)  
[Delete me](#del)  

## Build a site

```
mb build
```

### Build a site using global themes instead of local

When you create a site all the Metabuzz themes get copied to the site's `.themes` directory. That way if you change a theme it doesn't affect other sites created using that theme. If you want to use the original Metabuzz thems, use the `-d` option like this:

```
mb -d build
```



## Create a theme

To create a theme you copy an existing one:

```
d2 copytheme
```

You're asked which theme you'd like to copy:

```
Name of theme to copy?
```


## Create a new project

### init

#### Create a new project in the current directory

`init` by itself assumes you've created a new directory 
and changed to it, and that you want your project
to be based there. If it looks like there's
already a project, nothing will happen.

Example:

``` bash
$ ./docset init
```

### init -sitename

#### Create a new project directory

Example: Create a new project named foo

Create a directory named `foo` in the current directory.
and populate it with the required fies.

``` bash
$ ./docset init -sitename=foo
```

If there's already a project directory by that name,
it stays in place and the command returns an error. Nothing gets deleted.

## new

#### Create a new theme

To create a new theme, use `new`, followed by `theme=` and  
a valid directory (and file) name. In this example,
the new theme is named `foo2`.
```
docset new -theme=foo2 
```

If there's already a theme directory by that name,
it stays in place and the command returns an error. Nothing gets deleted.

Using `-theme=` without the `-from` flag creates a theme based on the default
themes. You can choose any theme to be the [default](product-config.html#defaulttheme).

### Create a theme from any specified theme

Using the command line you can also create a theme from any
theme you choose, not just the default. Specify that theme's name
using the `-from=` flag:

```
docset new -theme=foo2 -from=clean
```
If there's already a theme directory by that name,
it stays in place and the command returns an error. Nothing gets deleted.


## -showsettings

## Display configuration {: #display-configuration }

## DELETE ME {#del}

Shows global configuration settings such as where the configuration
file lives, the location of the user's home directory, etc.

Example:

``` bash
$ mb info
```

## Verbose -v flag

Metabuzz can process files much faster the standard output displays text, so it doesn't say much when you're building a site. If you get errors that you can't figure out, it may help to turn on Verbose mode during a build. Here's an example:

```
mb new site mysite
cd mysite
mb -v build
```

The output will look something like this:

```
[/Users/jared/html/mysite/index.md]
	Theme: [wide]
	Mode: [light]
	Created file [/Users/jared/html/mysite/.pub/index.html]
```


