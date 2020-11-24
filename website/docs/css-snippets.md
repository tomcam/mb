# CSS Snippets

article > h3:before { content: 'foo'; }

See css-snippets.txt

## Remove bullet characters from nested ul li for table of contents generation
article > ul > li, 
  article > ul > li > ul > li, 
  article > ul > li > ul > li > ul > li,
  article > ul > li > ul > li > ul > li > ul > li,
  article > ul > li > ul > li > ul > li > ul > li > ul > li,
  article > ul > li > ul > li > ul > li > ul > li > ul > li, ul > li
  {list-style-type:none;}

## Last item shows as button-type thing

```
header > ul > li:last-child > a {border: 2px solid var(--header-fg);color:var(--header-fg);}  
```
## Gallery thingie

article > h2 {margin-top:4rem;clear:left;}
article > h2 + ul {
	padding:0;
	margin:0;
	list-style-type:none;
}

/* Leave some room on the right side of each image */
article > h2 + ul li > img {
	width:95%;
  box-shadow: rgb(128,128,128) 1px 1px 3px 0px;
  margin-bottom:.75rem;
}

/* Each li is 1/2 as wide as the container. */
article > h2 + ul li {
	font-family:var(--informal);
	max-width:50%;
	width:50%;
	float:left;
	line-height:1em;
  font-size:.8em;
  margin-bottom:2rem;
  margin-top:2rem;
}

/* Start a new line after the 2nd column */
article > h2 + ul li:nth-child(2) {
	display:block;
}


## Backup cpy of ballery thingie Iguess

article > h2 {margin-top:4rem;clear:left;}
article > h2 + ul {
	padding:0;
	margin:0;
	list-style-type:none;
}

/* Leave some room on the right side of each image */
article > h2 + ul li > img {
	width:95%;
  box-shadow: rgb(128,128,128) 1px 1px 3px 0px;
  margin-bottom:.75rem;
}

/* Each li is 1/2 as wide as the container. */
article > h2 + ul li {
	font-family:var(--informal);
	max-width:50%;
	width:50%;
	float:left;
	line-height:1em;
  font-size:.8em;
  margin-bottom:2rem;
  margin-top:2rem;
}

/* Start a new line after the 2nd column */
article > h2 + ul li:nth-child(2) {
	display:block;
}





## Card

See https://codepen.io/edeesims/pen/iGDzk

```
/* Features specific to this theme */
h1:hover{opacity:0}
h1:hover{transform:rotate(180eg);transition:transform 0.5s;}
h1 {
  position:absolute;
  top:50%;
  left:50%;
  width:15rem;
  height:20rem;
  margin:-25%;
  float:left;
  perspective:50rem;
  background-color:whitesmoke;
  padding:15% 5% 5% 5%;
  border:1px solid lightgray;
  box-shadow: 0 0 3px rgba(0,0,0,0.1);
  text-align:center;
  border-radius:1rem;
  backface-visibility:hidden;
}

h2 {
  position:absolute;
  top:50%;
  left:50%;
  width:15rem;
  height:20rem;
  margin:-25%;
  float:left;
  perspective:50rem;
  background:whitesmoke;
  padding:15% 5% 5% 5%;
  border:1px solid lightgray;
  box-shadow: 0 0 3px rgba(0,0,0,0.1);
  text-align:center;
  border-radius:1rem;
  backface-visibility:hidden;
}

h5 {
  position:absolute;
  top:50%;
  left:50%;
  width:15rem;
  height:20rem;
  margin:-25%;
  transform:rotate(2deg);
  float:left;
  perspective:50rem;
  background:whitesmoke;
  padding:15% 5% 5% 5%;
  border:1px solid lightgray;
  box-shadow: 0 0 3px rgba(0,0,0,0.1);
  text-align:center;
  border-radius:1rem;
  backface-visibility:hidden;
}

h4 {
  position:absolute;
  top:50%;
  left:50%;
  width:15rem;
  height:20rem;
  margin:-25%;
  transform:rotate(4deg);
  float:left;
  perspective:50rem;
  background:whitesmoke;
  padding:15% 5% 5% 5%;
  border:1px solid lightgray;
  box-shadow: 0 0 3px rgba(0,0,0,0.1);
  text-align:center;
  border-radius:1rem;
  backface-visibility:hidden;
}
```

