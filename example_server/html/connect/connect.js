function Connect(msgListeners, locListeners) {
	var handleLoc = function(loc) {
		locListeners.forEach(function(listener) {listener.handleLoc(loc)});
	}
	var handleMsg = function(msg) {
		msg.Msg.Content = JSON.parse(msg.Msg.Content);
		msgListeners.forEach(function(listener) {listener.handleMsg(msg)});
	}
	this.msgListeners = msgListeners;
	this.locListeners = locListeners;
	this.msgService = new WSClient("Message", "ws://178.79.176.206:8003/msg", handleMsg, function(){}, function() {});
	this.locService = new WSClient("Location", "ws://178.79.176.206:8002/loc", handleLoc, function(){}, function() {});
	this.msgService.connect();
	this.locService.connect();
	this.id = getId();
	var addMsg = new Add(this.id);
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

Connect.prototype.sendMsg = function(msg) {
	msg.Msg.Content = JSON.stringify(msg.Msg.Content);
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

var requestCode = "sync-request";
var responseCode = "sync-response";

function SyncRequest(to, syncName) {
	return new Msg(to, {sync: requestCode, name: syncName});
}

function SyncResponse(to, syncName) {
	return new Msg(to, {sync: responseCode, name: syncName});
}

Connect.prototype.sync = function(idMe, idYou, syncName, fun) {
	var synced = false;
	var cnct = this;
	// NB: The correctness of this approach relies on the interval function being unable to run even once before this function has completed
	intervalId = setInterval(function() {cnct.sendMsg(SyncRequest(idYou, syncName));}, 300);
	var syncListener = function(msg) {
		var from = msg.Msg.From;
		var content = msg.Msg.Content;
		if (content.sync == requestCode) {
			var name = content.name;
			if (name == syncName && from == idYou) {
				clearInterval(intervalId);
				cnct.rmvMsgListener(syncListener);
				cnct.sendMsg(SyncRequest(idYou, syncName));
				fun();
			} else {
				console.log("Received 2sync request with unexpected id " + id + " from " + from);
			}
		} else if (content.sync == responseCode) {
			var name = content.name;
			if (name == syncName && from == idYou) {
				clearInterval(intervalId);
				cnct.rmvMsgListener(syncListener);
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
