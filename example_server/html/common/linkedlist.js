function LinkedList() {
	this.size = 0;
	this.first = null;
	this.last = null;
}

function Item(val) {
	this.val = val;
	this.next = null;
	this.prev = null;
}

LinkedList.prototype.getFirst = function() {
	return this.first.val;
}

LinkedList.prototype.getLast = function() {
	return this.last.val;
}

LinkedList.prototype.length = function() {
	return this.size;
}

LinkedList.prototype.clear = function() {
	this.first = null;
	this.last = null;
	this.size = 0;
}

LinkedList.prototype.append = function(val) {
	this.size++;
	var item = new Item(val);
	if (this.first == null) {
		this.first = item;
		this.last = item;
	} else {
		this.last.next = item;
		item.prev = this.last;
		item.next = null;
		this.last = item;
	}
}

LinkedList.prototype.forEach = function(fun) {
	var item = this.first;
	while (item != null) {
		fun(item.val);
		item = item.next;
	}
}

LinkedList.prototype.filter = function(pred) {
	var item = this.first;
	while (item != null) {
		if (pred(item.val)) {
			this.size--;
			item = remove(this, item);
		} else {
			item = item.next;
		}
	}
}

LinkedList.prototype.satOne = function(pred) {
	if (this.length() == 0) {
		return false;
	}
	var sat = false;
	this.forEach(function(e){sat || pred(e)});
	return sat;
}

LinkedList.prototype.satAll = function(pred) {
	var sat = true;
	this.forEach(function(e){sat = (sat && pred(e));});
	return sat;
}

remove = function(list, item) {
	if (item.prev != null) {
		item.prev.next = item.next;
	} else {
		list.first = item.next;
	}	
	if (item.next != null) {
		item.next.prev = item.prev;
	} else {
		list.last = item.prev;
	}
	return item.next;
}

LinkedList.prototype.circularNext = function(val) {
	var item = this.first;
	while (item != null) {
		if (item.val === val) {
			if (item.next == null) {
				return this.first.val;
			} else {
				return item.next.val;
			}
		}
		item = item.next;
	}
	return null;
}

LinkedList.prototype.contains = function(val) {
	var item = this.first;
	while (item != null) {
		if (item.val === val) {
			return true;
		} else {
			item = item.next;
		}
	}
	return false;
}
