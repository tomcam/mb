# Global Configuration File

## Changing where application configuration data is stored

Upon installation themes and other customizable assets are stored whereever
your system expects user application data to be, in a directory named
`metabuzz/.mb`. For example, for a user named Rajiv on a Mac the directory containing theme files would be something like:

```
/Users/rajiv/Library/Application Support/metabuzz/.mb/themes

```

You can change that location easily. One reason to do this is that you might want your themes in the same Git repository as your project, which means they'd have to be in a subdirectory of your project. You can't fake Git by creating a symlink (also known as an alias), because Git will slap you down for attempting to use symlinks.

Metabuzz looks for a global configuration file named `.mb/metabuzz.toml` in  your home directory when it starts. If it finds such a file, it looks at the `configdir` setting to see if you've decided to move the the global config file. Here's how to take advantage of this and change the config file's directory location.

```
# Create a directory named .mb in 
# your home directory.
mkdir ~/.mb
# Create a config file in that directory
nvim ~/.mb/metabuzz.toml

```

Now edit the file to include the new location for theme files. In this example, your username is Rajiv and you want the theme files to be stored with your current project, which is called `mymanual` and is found in your home directory under `/docs`:

```
configdir="/Users/rajiv/docs/mymanual/.mb"
```

Don't use designations like `~` because those are expected to be resolved by the shell before Metabuzz is invoked.

Now all theme files are found not in `/Users/rajiv/Library/Application Support/metabuzz/.mb/themes` but inside of Rajiv's `mymanual` directory. Don't worry; this doesn't add more than a few hundred K of mostly text files to the project.



