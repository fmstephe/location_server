function WSClient(name, url, msgFun, opnFun, clsFun) {
	
	this.jsonsend = jsonsend;
	this.name = name;
	this.msgFun = msgFun;
	this.opnFun = opnFun;
	this.clsFun = clsFun;
	
	this.connect = function() {
		if ("WebSocket" in window) { this.ws = new WebSocket(url); }
		else if ("MozWebSocket" in window) { this.ws = new MozWebSocket(url); }
		this.ws.onopen = onopen;
		this.ws.onmessage = onmessage;
		this.ws.onclose = onclose;
		this.ws.jsonsend = jsonsend;
		this.ws.name = this.name;
		this.ws.msgFun = this.msgFun;
		this.ws.opnFun = this.opnFun;
		this.ws.clsFun = this.clsFun;
		this.ws.earlyMsgs = new LinkedList();
	}
}

function onopen() {
	console.log(this.name+" Websocket Connection Open!");
	this.opnFun();
	var wsClosure = this;
	this.earlyMsgs.forEach(function(obj) {wsClosure.jsonsend(obj)});
}

function jsonsend(obj) {
	if (this.ws) {
		this.ws.jsonsend(obj);
	} else {
		if (this.readyState == undefined || this.readyState == 0) { // in opening state
			this.earlyMsgs.append(obj);
			console.log(this.name + ": early message stored: "+JSON.stringify(obj));
		} else {
			msg = JSON.stringify(obj);
			this.send(msg);
			console.log(this.name + ": json msg delivered: "+msg);
		}
	}
}

function onmessage(m) { 
	if (m.data) {
		console.log(this.name + ": json msg received: "+m.data);
		var msg = JSON.parse(m.data);
		if (msg.Op == "sAck") {
			this.unackedMsgs.filter(function(tsMsg) {return tsMsg.msg.Id == msg.Id});
		} else {
			this.msgFun(msg);
		}
	}   
}

function onclose(m) {
	console.log(this.name+" Websocket Connection Closed!");
	this.clsFun();
	this.ws=null;
}
