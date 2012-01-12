function render(ctxt, obj) {
	var list = obj.clearings();
	list.forEach(function(elem) {clear(ctxt,elem.clr)} );
	var list = obj.renderings();
	list.forEach(function(elem) {elem.render(ctxt);} );
}

function clear(ctxt, clr) {
	ctxt.clearRect(clr.x,clr.y,clr.width,clr.height);
}
