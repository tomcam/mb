<script>
var scrollTimer = -1;
function bodyScroll(){if(scrollTimer != -1){clearTimeout(scrollTimer);}scrollTimer = window.setTimeout(sidebarHeight, 100);}window.onresize=sidebarHeight;
function sidebarHeight() {
s=document.getElementById('sidebar');
a=document.getElementById('article'); 
if (s != null && a != null) {
  // If there's a sidebar,
  ha=a.offsetHeight;
  hs=s.offsetHeight;
  // If the article is longer, make
  // the sidebar as long as the article.
  if (ha>hs) {
    h=a.offsetHeight+'px';s.style.height=h;
  }
  // If the sidebar is longer than the
  // article, make the article as long
  // as the sidebar.
  else{
    h=s.offsetHeight+'px';a.style.height=h;
  }
}
document.onreadystatechange = function () {
if (document.readyState == "interactive") {
  // Init or start code here
  bodyScroll();
}
} 
  
  </script>
<body onscroll="bodyScroll();">

