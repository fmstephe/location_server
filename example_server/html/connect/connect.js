function Connect(msgListeners, locListeners) {
	var handleLoc = function(loc) {
		locListeners.foreach(function(listener) {listener.handleLoc(loc)});
	}
	var handleMsg = function(msg) {
		msgListeners.foreach(function(listener) {listener.handleMsg(msg)});
	}
	this.msgListeners = msgListeners;
	this.locListeners = locListeners;
	this.msgService = new WSClient("Message", "ws://178.79.176.206:8003/msg", handleMsg, function(){}, function() {});
	this.locService = new WSClient("Location", "ws://178.79.176.206:8002/loc", handleLoc, function(){}, function() {});
	this.msgService.connect();
	this.locService.connect();
	var id = getId();
	var addMsg = new Add(id);
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

var 2syncRequest = "2sync-request";
var 2syncResponse = "2sync-response";

Connect.prototype.2sync = function(idMe, idYou, syncFun) {
	var synced = false;
	var cnct = this;
	// NB: The correctness of this approach relies on the interval function being unable to run even once before this function has completed
	intervalId = setInterval(function() {cnct.sendMsg(new Msg(idYou, JSON.stringify({op: 2syncRequest, id: idYou})));}, 300);
	var syncListener = function(msg) {
		var from = msg.From;
		var content = JSON.parse(msg.Msg.Content);
		if (content.op == 2syncRequest) {
			var id = content.id;
			clearInterval(intervalId);
			cnct.rmvMsgListener(syncListener);
			if (id == idMe && from == idYou) {
				cnct.sendMsg(idYou, JSON.stringify({op: 2syncResponse, id: idYou}));
				syncFun();
			} else {
				console.log("Received 2sync request with unexpected id " + id + " from " + from);
			}
		} else if (content.op == 2syncResponse) {
			var id = content.id;
			clearInterval(intervalId);
			cnct.rmvMsgListener(syncListener);
			if (id == idMe && from == idYou) {
				syncFun();
			} else {
				console.log("Received 2sync response with unexpected id " + id + " from " + from);
			}
		}
	}
	this.addMsgListener(syncListener);
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
