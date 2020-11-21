===
#theme="new-wide"
#theme="new-pillar"
#pagetype="press-release"
#theme="_"
theme="textual"
sidebar="none"
#mode="dark"
#mode="light"
===
<svg xmlns="http://www.w3.org/2000/svg" width="32px" height="32px"  viewBox="0 0 20 20" fill="currentColor">
  <path fill-rule="evenodd" d="M12 1.586l-4 4v12.828l4-4V1.586zM3.707 3.293A1 1 0 002 4v10a1 1 0 00.293.707L6 18.414V5.586L3.707 3.293zM17.707 5.293L14 1.586v12.828l2.293 2.293A1 1 0 0018 16V6a1 1 0 00-.293-.707z" clip-rule="evenodd" />
</svg>

Run `themevars` with the name of the them:

```
./themevars textual
```



(Redmond, WA, November 17, 2020) Hey people!

Metabuzz automatically generates an id attribute for each header from h1 to h6 by taking the text of the link itself, reducing it to lowercase, and either replacing spaces and other non-letter characters with hyphens, or removing them altogether.

# Generating these tests

## Bug: child TOML pagetype is being ignored

{{ inc "theme-and-variations.md" }}
{{ inc "mdemo.md" }}

