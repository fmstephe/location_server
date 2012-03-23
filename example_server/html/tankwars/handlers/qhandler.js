function QHandler(filterFun) {
	this.filterFun = filterFun;
	this.q = new LinkedList();
}

function handleMsg(msg) {
	if (this.filterFun(msg)) {
	       this.q.append(msg);
	}
}
