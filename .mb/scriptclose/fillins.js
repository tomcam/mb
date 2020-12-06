<script>
var scrollTimer = -1;
function bodyScroll(){if(scrollTimer != -1){clearTimeout(scrollTimer);}scrollTimer = window.setTimeout(sidebarHeight, 100);}window.onresize=sidebarHeight;
function sidebarHeight() {
s=document.getElementById('sidebar');
a=document.getElementById('article'); 
h=a.offsetHeight+'px';
s.style.height=h;
//s.style.height='100vh';
}
document.onreadystatechange = function () {
if (document.readyState == "interactive") {
  // Init or start code here
  bodyScroll();
}
} 
  
  </script>
<body onscroll="bodyScroll();">

