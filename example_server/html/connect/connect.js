function new Connecter(msgListeners, locListeners) {
	var handleLoc = function(loc) {
		for (var i in locListeners) {
			locListeners[i].handleLoc(loc);
		}
	}
	var handleMsg = function(msg) {
		for (var i in msgListeners) {
			msgListeners[i].handleMsg(msg);
		}
	}
	var initLoc = function(position) {
		lat = position.coords.latitude;
		lng = position.coords.longitude;
		var locMsg = new InitLoc(lat, lng);
		locService.jsonsend(locMsg)
	}
	this.msgService = new WSClient("Message", "ws://178.79.176.206:8003/msg", handleMsg, function(){}, function() {});
	this.locService = new WSClient("Location", "ws://178.79.176.206:8002/loc", handleLoc, function(){}, function() {});
	this.msgService.connect();
	this.locService.connect();
	var id = getId();
	var addMsg = new Add(id);
	this.msgService.jsonsend(addMsg);
	this.locService.jsonsend(addMsg);
	setInitCoords(initLoc);
}

Connector.prototype.sendMsg = function(msg) {
	this.msgService.jsonsend(msg);
}

Connector.prototype.sendLoc = function(loc) {
	this.locService.jsonsend(loc);
}

function setInitCoords(initLoc) {
	if (navigator.geolocation) {
		navigator.geolocation.getCurrentPosition(initLoc, function(error) { console.log(JSON.stringify(error)), initLoc({"coords": {"latitude":1, "longitude":1}}) }); 
	} else {
		alert("Your browser does not support websockets");
	}
}

function init(position) {
	lat = position.coords.latitude;
	lng = position.coords.longitude;
	var locMsg = new InitLoc(lat, lng);
	
}
