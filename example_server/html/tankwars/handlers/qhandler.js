function QHandler(filterFun) {
	this.filterFun = filterFun;
	this.q = new LinkedList();
}

QHandler.prototype.handleMsg = function(msg) {
	if (this.filterFun(msg)) {
	       this.q.append(msg);
	}
}
