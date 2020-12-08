# CSS Snippets

article > h3:before { content: 'foo'; }

See css-snippets.txt


## Show unordered list in a sidebar as numbers in circles For a news-type section


```
/*
 * --------------------------------------------------
 * Sidebar unordered list for breaking news
 * shows:
 * - Item preceded by number in circle
 * - Item is bold
 * - Indented item in normal text with bottom padding
 *
 * Example usage (note 2nd level of indentation):
 *
 *  * Item 1
 *     - More about item 1
 *   * Item 2
 *     - More about item 2
 *   * Item 3
 *     - More about item 3
 *
 * --------------------------------------------------
 */

aside > ul {
  counter-reset:li;
  list-style-type:none;
  font-size:1rem;
  line-height:1.2rem;
  padding-left:1em;
  border-top:none;
}

aside > ul li {
  list-style-type:none;
  border:none;
  font-weight:normal;
}

aside > ul > li > ul > li {
  line-height:1.5em;
  padding-bottom:1em;
}
aside > ul >li {
  font-weight:bold;
  position:relative;
  padding: .5em 1em 0rem 3em;
  border:none;
} 
aside > ul > li:before {
  background-color:var(--fg);
  content:counter(li);
  counter-increment:li;
  height:2rem;
  line-height:2rem;
  width:2em;
  border-radius:50%;
  color: var(--bg);
  text-align:center;
  position:absolute;
  left:0;
}

``

## Under an h3, show 3 pictures per row with a caption area underneath


```
/*
 * --------------------------------------------------
 * Special feature: "Featured Posts" section showing
 * 3 images in a row with captions underneath.
 * Start with an h3 ("Featured Posts, for example)
 * then a ul with 3 items consisting of an image
 * and text for it.
 *
 *
 * Example:
 *   ### Featured stories
 *   * ![Picture of Oasis](oasis.png)
 *     #### Oasis reunites
 *     [MORE](oasis-reunites.html)
 *   * ![Picture of bookcase](bookcase.png)
 *     ##### Man, what a bookcase 
 *     [MORE](bookcase.html)
 *   * ![Picture of yomama](yomama.png)
 *     #### I won't say it 
 *     [MORE](yomama.html)
 *
 * Notes:
 * - Remember to indent after the bullet.
 * If you have less than 3 items, just put an empty
 * header 3 after the first 1 or 2:
 *
 * --------------------------------------------------
 */

article > h3 {clear:left;}
article > h3 + ul {
	padding:0;
	margin:0;
	list-style-type:none;
}

/* Each li is 1/3 as wide as the container. */
article > h3 + ul li {
  background-color:white;
	font-family:var(--informal);
	width:28.3%;
	float:left;
	line-height:1em;
  margin-right:5%;
  border: 1px solid whitesmoke;
}

/* Leave some room on the right side of each image */
article > h3 + ul li > img {
  /* Not necessary */
  display:inline;
  width:100%;
 }

article > h3 + ul > li > h4 {
 font-size:.8em;
 font-family:var(--times);
 padding-left:1em;
}

article > h3 + ul > li > h4 + p{
 padding-left:1em;
 padding-bottom:2em;
}
```

### Example usage:

```
### News Roundup
* ![2](2.jpg) 
  #### h4 here. Let's see what a long caption looks like
  [hello](/)
* ![2](2.jpg) 
  #### h4 here
  [hello](/)
* ![2](2.jpg) 
  #### h4 here
  [hello](/)
```



## 12/7/20: Create table of contents in sidebar that shows as boxes of text - maybe a duplicaet of next one, not sure. Good for toc

```
aside {padding-left:1em;}
/*
 * --------------------------------------------------
 * Sidebar unordered list shows as boxes, without
 * indentation--it's for table of contents
 * --------------------------------------------------
 */
aside > ul {
  background-color:whitesmoke;
  margin-right:1em;
  margin-left:var(--left-margin);
  border-collapse:collapse;
}


aside > ul li {
  list-style-type:none;
  /* Border bottom stretches across column at all levels 
   * but produces thicker border for cases in which
   * a header (say, an h4) is followed by a higher
   * level header (say, an h2).*/
  border-top: 1px solid gray;
  box-shadow: rgb(128,128,128) 1px 1px 1px 0px;
}

aside > ul li a {
  padding-left:.5em;padding-right:.5em;
  text-decoration:none;
  line-height:1em;
  /* Only underlines text. Doesn't stretch across the whole
   * sidebar width. 
   * border-bottom: 1px solid gray;
   */
}

aside > ul li a:active, aside > ul li a:hover {
  font-weight:bold;
}

/* Special case: h1 for table of contents
 * gets distinct look */
aside > ul > li > a {
  font-weight:bold;
}


```



## Create table of contents in sidebar that shows as boxes of text

```
/*
 * --------------------------------------------------
 * Sidebar unordered list shows as boxes, without
 * indentation--it's for table of contents
 * --------------------------------------------------
 */
aside > ul {
  background-color:whitesmoke;
  margin-right:1em;
  border-collapse:collapse;
}

aside > ul li {
  list-style-type:none;
  margin-left:0;
  /* Border bottom stretches across column at all levels */
  border-bottom: 1px solid gray;
}

aside > ul li a {
  padding-left:.5em;padding-right:.5em;
  text-decoration:none;
  line-height:1em;
}

aside > ul li a:active, aside > ul li a:hover {
  /* For illustrative purposes only */
  color:blue;
}


```



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

