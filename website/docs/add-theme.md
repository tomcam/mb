# Adding a theme to the Metabuzz gallery

Most of the theme gallery is automated. You need to provide a few assets 
that are unique to your theme, and a standardized test page is created 
with sidebars, without sidebars, and in dark mode. (You can omit any of these except the light mode theme and the `theme-1280x1280.png` file.)

Apart from the theme, you'll need:

1. A one-paragraph description of the theme
2. A 1280x1280 .png screenshot of your favorite version of the theme named, literally, `theme-1280x1280.png`.
3. Thumbnail screenshots, all PNG files at a resolution of 256x256. If you don't have one, you don't need to include it.
   * Light theme version showing the left sidebar named `light-sidebar-left-256x256.png`
   * Light theme version showing the right sidebar named `light-sidebar-right-256x256.png`
   * Light theme version with no sidebar named light-sidebar-none-256x256.png
   * Dark theme version showing the left sidebar named `dark-sidebar-left-256x256.png`
   * Dark theme version showing the right sidebar named `dark-sidebar-right-256x256.png`
   * Dark theme version with no sidebar named `dark-sidebar-none-256x256.png`
4. A file named `intro.md` that contains the Markdown for a brief article, typically preceded with an image. The `intro.md` file can be empty.
5. A file named `right-sidebar-example.md` containing the markdown for a right sidebar example fitting that theme.
6. A file named `left-sidebar-example.md` containing the markdown for a left sidebar example fitting that theme. It' okay if the left and right sidebars are identical, but if they're designed differently in the style sheet this is a way to show them both off.

