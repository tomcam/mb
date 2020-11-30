**Sale** in the [shop](shop.html) today only!

{{- /* the sidebar element doesn't resize the way we'd like. This keeps it filling the height as it should */ -}}
<script>
function sidebarHeight() {
s=document.getElementById('sidebar');
a=document.getElementById('article'); 
h=a.offsetHeight+'px';
s.style.height=h;
}
window.onresize=sidebarHeight;
</script>

