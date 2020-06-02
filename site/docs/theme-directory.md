# Theme directory

# TODO start this bad boy


## How to find out what theme you're using

If you're not too familiar with knowing where your [site file](site-file.html) or [config file](config-file.html) are, you can find what theme you're using very simply. Create a file something like this:

```
```

## Location of themes directory

By default, the `.mb/themes` directory is stored where the host operating
system looks for application configuration data. You can find 
out where that directory is at the command line using the
`mb info` command, like this:

```
mb info
```

Among the information it displays you'll see something like
this. Note that it's operating system dependent. This is a result you
might see on MacOS for a user named Jared:

```
Theme directory /Users/jared/Library/Application Support/metabuzz/.mb: (present)
```

You'll need to know where your themes directory is before you can
complete this tutorial.

## Changing the location of the themes directory (helpful for Git users)

The `.mb` directory is for global configuration, meaning it tells Metabuzz
where to find configuration data for all Metabuzz sites you run, 
including themes. 

You can change where Metabuzz looks for themes by moving `.mb` itself. If you 
use Git this can make sense when you want to ensure the theme files
are part of the same repository as the rest of your site. 

See [Global Configuration File](config-file.html) for details on how
to move the `.mb` directory.


