function Connect(userId, msgHandlers, locHandlers, clsFun, locatedFun) {
	var thisConn = this;
	var url = document.URL;
	url = url.replace("http", "ws");
	var handleLoc = function(loc) {
		locHandlers.forEach(function(handler) {handler.handleLoc(loc)});
	}
	var handleMsg = function(msg) {
		msg.content = JSON.parse(msg.content);
		msgHandlers.forEach(function(handler) {handler.handleMsg(msg)});
	}
	this.msgHandlers = msgHandlers;
	this.locHandlers = locHandlers;
	this.msgService = new WSClient("Message", url + "msg", handleMsg, function(){}, clsFun);
	this.locService = new WSClient("Location", url + "loc", handleLoc, function(){}, clsFun);
	this.handleMsgLocal = handleMsg;
	this.msgService.connect();
	this.locService.connect();
	this.unackedMsgs = new LinkedList();
	this.usrId = userId;
	var addMsg = new Add(this.usrId);
	this.msgService.jsonsend(addMsg);
	this.locService.jsonsend(addMsg);
	var lsvc = this.locService;
	var initLoc = function(position) {
		lat = position.coords.latitude;
		lng = position.coords.longitude;
		var locMsg = new InitLoc(lat, lng);
		lsvc.jsonsend(locMsg);
		locatedFun();
	}
	setInitCoords(initLoc);
}

Connect.prototype.sendMsg = function(to, content) {
	var msg = new Msg(to, JSON.stringify(content));
	this.msgService.jsonsend(msg);
}

Connect.prototype.sendLoc = function(loc) {
	this.locService.jsonsend(loc);
}

Connect.prototype.addMsgHandler = function(handler) {
	this.msgHandlers.append(handler);
}

Connect.prototype.rmvMsgHandler = function(handler) {
	this.msgHandlers.filter(function(l) {return handler === l;});
}

Connect.prototype.addLocHandler = function(handler) {
	this.locHandlers.append(handler);
}

Connect.prototype.rmvLocHandler = function(handler) {
	this.locHandlers.filter(function(l) {return handler === l;});
}

Connect.prototype.close = function() {
	this.msgService.close();
	this.locService.close();
}

var syncRequest = {isSyncRequest: true};
var syncResponse = {isSyncResponse: true};

Connect.prototype.sync = function(idMe, idYou, fun) {
	var thisConn = this;
	// NB: The correctness of this approach relies on the interval function being unable to run even once before this function has completed
	// Otherwise the SyncRequest might be sent, and responded to, before the syncHandler is registered (just echos of threading paranoia)
	var intervalId = setInterval(function() {thisConn.sendMsg(idYou, syncRequest);}, 300);
	var syncHandler = {}; // Predeclaration so we can refer to this object inside syncHandler
	var handle = function(msg) {
		var from = msg.from;
		var content = msg.content;
		if (content.isSyncRequest) {
			var name = content.name;
			if (from == idYou) {
				clearInterval(intervalId);
				thisConn.rmvMsgHandler(syncHandler);
				thisConn.sendMsg(idYou, syncResponse);
				fun();
			} else {
				console.log("Received sync request from unexpected user:" + from);
			}
		} else if (content.isSyncResponse) {
			var name = content.name;
			if (from == idYou) {
				clearInterval(intervalId);
				thisConn.rmvMsgHandler(syncHandler);
				fun();
			} else {
				console.log("Received sync response from unexpected user:" + from);
			}
		}
	}
	syncHandler.handleMsg = handle;
	this.addMsgHandler(syncHandler);
}

function setInitCoords(initLoc) {
	if (navigator.geolocation) {
		navigator.geolocation.getCurrentPosition(initLoc, function(error) { console.log(JSON.stringify(error)), initLoc({"coords": {"latitude":1, "longitude":1}}) }); 
	} else {
		initLoc({"coords": {"latitude":1, "longitude":1}});
	}
}

function init(position) {
	lat = position.coords.latitude;
	lng = position.coords.longitude;
	var locMsg = new InitLoc(lat, lng);
}
