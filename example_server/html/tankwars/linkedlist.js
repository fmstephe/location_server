function LinkedList() {
	this.size = 0;
	this.first = null;
	this.last = null;
	this.getFirst = function() {return this.first.val;};
	this.getLast = function() {return this.last.val;};
	this.clear = clear;
	this.append = append;
	this.forEach = forEach;
	this.filter = filter;
	this.circularNext = circularNext;
	this.contains = contains;
	this.length = function() {return this.size;};
}

function clear() {
	this.first = null;
	this.last = null;
	this.size = 0;
}

function Item(val) {
	this.val = val;
	this.next = null;
	this.prev = null;
}

function append(val) {
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

function forEach(fun) {
	var item = this.first;
	while (item != null) {
		fun(item.val);
		item = item.next;
	}
}

function filter(pred) {
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

function remove(list, item) {
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

function circularNext(val) {
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

function contains(val) {
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
