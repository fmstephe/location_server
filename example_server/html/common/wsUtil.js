function WSClient(name, url, msgFun, opnFun, clsFun) {
	
	this.jsonsend = jsonsend;
	this.name = name;
	this.msgFun = msgFun;
	this.opnFun = opnFun;
	this.clsFun = clsFun;
	
	this.connect = function() {
		this.ws = new WebSocket(url);
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
		if (this.readState == 0) { // in opening state
			this.earlyMsgs.append(obj);
			console.log("early message stored: "+JSON.stringify(obj));
		} else {
			msg = JSON.stringify(obj);
			this.send(msg);
			console.log("json message sent: "+msg);
		}
	}
}

function onmessage(m) { 
	if (m.data) {
		console.log(m.data);
		var msg = JSON.parse(m.data);
		this.msgFun(msg);
	}   
}

function onclose(m) {
	console.log(this.name+" Websocket Connection Closed!");
	this.clsFun();
	this.ws=null;
}
