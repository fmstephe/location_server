function new Connecter(id, msgListeners, locListeners) {
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

function introFunc() {
	users = "";
	userList.forEach(function(u) {users += "<li>"+JSON.stringify(u)+"</li>"});
	document.getElementById("player-list").innerHTML = users;
}

/*
	console.log(msg);
	var m = JSON.parse(msg.Msg);
	var p = new Player(m.x, m.name, turretLength, m.power, minPower, maxPower, powerInc, expRadius, null);
	playerList.append(p);
}

function handleLoc(msg) {
	var op = msg.Op;
	console.log(op);
	var usrInfo = msg.Usr;
	if (op == "sAdd" || op == "sNearby" || op == "sVisible") {
		playerMsg = new Msg(usrInfo.Id, JSON.stringify(new PlayerMsg(localPlayer)));
		msgService.jsonsend(playerMsg);
		userList.append(usrInfo);
	} else if (op == "sRemove" || op == "sNotVisible") {
		userList.filter(function(u) {return usrInfo.Id == u.Id});
	}   
}*/
