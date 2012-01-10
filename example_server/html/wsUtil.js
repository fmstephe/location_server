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
	}
}

function onopen() {
	console.log(this.name+" Websocket Connection Open!");
	this.opnFun();
}
function jsonsend(obj) {
	msg = JSON.stringify(obj);
	console.log(msg);
	if (this.ws) {
		this.ws.send(msg);
	} else {
		this.send(msg);
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
