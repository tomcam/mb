{{- /*  This "includes" the file showtheme.md from the  
        current markdown file's directory. It copies in the
        contents as if you'd typed them yourself. 
        
*/ -}}

{{ inc "common|showtheme.md" }}

{{ inc "intro.md" }}

{{ inc "description.md" }}

{{ inc "common|variations.md" }}

{{- /*  This "includes" the file named kitchen.md from the .common 
        directory. 
*/ -}}
{{ inc "common|kitchen.md" }}



