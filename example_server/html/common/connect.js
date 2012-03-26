function Connect(msgListeners, locListeners) {
	var thisConn = this;
	var handleLoc = function(loc) {
		locListeners.forEach(function(listener) {listener.handleLoc(loc)});
	}
	var handleMsg = function(msg) {
		msg.Content = JSON.parse(msg.Content);
		msgListeners.forEach(function(listener) {listener.handleMsg(msg)});
	}
	this.msgListeners = msgListeners;
	this.locListeners = locListeners;
	this.msgService = new WSClient("Message", "ws://178.79.176.206:8003/msg", handleMsg, function(){}, function() {});
	this.locService = new WSClient("Location", "ws://178.79.176.206:8002/loc", handleLoc, function(){}, function() {});
	this.msgService.connect();
	this.locService.connect();
	this.unackedMsgs = new LinkedList();
	this.usrId = getId();
	var addMsg = new Add(this.usrId);
	this.msgService.jsonsend(addMsg);
	this.locService.jsonsend(addMsg);
	var lsvc = this.locService;
	var initLoc = function(position) {
		lat = position.coords.latitude;
		lng = position.coords.longitude;
		var locMsg = new InitLoc(lat, lng);
		lsvc.jsonsend(locMsg)
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

Connect.prototype.addMsgListener = function(listener) {
	this.msgListeners.append(listener);
}

Connect.prototype.rmvMsgListener = function(listener) {
	this.msgListeners.filter(function(l) {return listener == l;});
}

Connect.prototype.addLocListener = function(listener) {
	this.locListeners.append(listener);
}

Connect.prototype.rmvLocListener = function(listener) {
	this.locListeners.filter(function(l) {return listener == l;});
}

Connect.prototype.close = function() {
	this.msgService.close();
	this.locService.close();
}

var requestCode = "sync-request";
var responseCode = "sync-response";

function SyncRequest(syncName) {
	return {sync: requestCode, name: syncName};
}

function SyncResponse(syncName) {
	return {sync: responseCode, name: syncName};
}

Connect.prototype.sync = function(idMe, idYou, syncName, fun) {
	var synced = false;
	var thisConn = this;
	// NB: The correctness of this approach relies on the interval function being unable to run even once before this function has completed
	// Otherwise the SyncRequest might be sent, and responded to, before the syncListener is registered (just echos of threading paranoia)
	var intervalId = setInterval(function() {thisConn.sendMsg(idYou, SyncRequest(syncName));}, 300);
	var syncListener = function(msg) {
		var from = msg.From;
		var content = msg.Content;
		if (content.sync == requestCode) {
			var name = content.name;
			if (name == syncName && from == idYou) {
				clearInterval(intervalId);
				thisConn.rmvMsgListener(syncListener);
				thisConn.sendMsg(idYou, SyncResponse(syncName));
				fun();
			} else {
				console.log("Received 2sync request with unexpected id " + id + " from " + from);
			}
		} else if (content.sync == responseCode) {
			var name = content.name;
			if (name == syncName && from == idYou) {
				clearInterval(intervalId);
				thisConn.rmvMsgListener(syncListener);
				fun();
			} else {
				console.log("Received 2sync response with unexpected id " + id + " from " + from);
			}
		}
	}
	this.addMsgListener({handleMsg: syncListener});
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
