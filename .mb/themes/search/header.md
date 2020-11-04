{{- /*  Automatically name first item in header    
        based on company name, author name name
        if no company was specified, or just 
        the name of the theme if neither of those
        was specified.
        
*/ -}}
{{- if .Site.Company.Name -}}
{{- $name := .Site.Company.Name -}}
* [{{ $name -}}](/)
{{- else if .Site.Author.FullName -}}
{{- $name := .Site.Author.FullName -}}
* [{{ $name -}}](/)
{{- else }}
* [{{.FrontMatter.Theme}}](/)
{{- end }} 
* [Events](/)
* [Subscribe](/)
<span id="demo">x</span>
<script>

//url = '/.pub/.indexing/docs.json'; 
url = 'metabuzz-search.json'
//https://stackoverflow.com/questions/7346563/loading-local-json-file
// https://stackoverflow.com/questions/48594581/asynchronous-callback-in-javascript

// https://stackoverflow.com/questions/7346563/loading-local-json-file
/*
function loadJSON(ex) {
  xobj = new XMLHttpRequest()
  xobj.addEventListener("load", reqListener)
  xobj.overrideMimeType("application/json")
  xobj.open('GET', url,true)
  alert('loadJSON(): xobj = ' + xobj)
  xobj.onreadystatechange = function() {
    if (xobj.readyState === 4 && xobj.status === 200) {
      ex(xobj.responseText)
    }
  };
  xobj.send(null)
  XmlHttpRequest.send(null)
}
*/

function loadJson(ex) {
  //alert('loadJson()')
  var XmlHttpRequest = new XMLHttpRequest();
  XmlHttpRequest.onreadystatechange = function () {
    //alert('onreadystatechange')
    if (XmlHttpRequest.readyState == 4 && XmlHttpRequest.status == "200") {
      // .open will NOT return a value 
      // but simply returns undefined in async mode so use a callback
      //alert('loadJson callback happening')
      document.getElementById("demo").innerHTML = XmlHttpRequest.name;
      //ex(XmlHttpRequest.responseText);
    }
  };
  XmlHttpRequest.overrideMimeType("application/json");
  XmlHttpRequest.open('GET', 'docs.json', true);
  XmlHttpRequest.send(null);
  //alert('I hope it is: ' + XmlHttpRequest.response)
  return (XmlHttpRequest.responseText)
}

/*
var callback = function(){
};
*/


function show(f) {
  alert(f)
 } 

function readTextFile(file, cb) {
    //alert('readTextFile() ' + file )
    var rawFile = new XMLHttpRequest();
    rawFile.overrideMimeType("application/json");
    rawFile.open("GET", file, true);
    rawFile.onreadystatechange = function() {
        if (rawFile.readyState == 4 && rawFile.status == "200") {
            document.getElementById("demo").innerHTML = rawFile.name;
            //cb(rawFile.responseText);
        }
    }
    rawFile.send(null);
}

function callback(){
  console.log('1')
  document.getElementById("demo").innerHTML = 'y';
  //alert('callback()')
  var xmlhttp = new XMLHttpRequest();
  xmlhttp.onreadystatechange = function() {
      console.log('this.readyState: ' + this.readyState + '. this.status: ' + this.status)
      if (this.readyState == 4 && this.status == 200) {
          var myObj = JSON.parse(this.responseText);
          document.getElementById("demo").innerHTML = myObj.name;
      }
  };
  xmlhttp.open("GET", "docs.json", true);
  xmlhttp.send();
}




//usage:


function lsearch(){
  alert('Searching for: ' + document.searchForm.search.value);
  return false;
}

// Ensure document's loaded before running Javascript
if (
    document.readyState === "complete" ||
    (document.readyState !== "loading" && !document.documentElement.doScroll)
) {
  callback()
} else {

  document.addEventListener("DOMContentLoaded", callback);
}
</script>
<form name="searchForm" onSubmit="lsearch()"><input type="text" id="search" name="search"><span class='icn icn-find'> </span>
</form>


