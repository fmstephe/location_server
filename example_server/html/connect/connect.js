function new Connecter(id, lat, lng, msgListeners, locListeners) {
	var addMsg = new Add(id);
	var locMsg = new InitLoc(lat, lng);
	this.msgListeners = msgListeners;
	this.locListeners = locListeners;
	this.msgService = new WSClient("Message", "ws://"+host+":8003/msg", handleMsg, function(){}, function() {});
	this.locService = new WSClient("Location", "ws://"+host+":8002/loc", handleLoc, function(){}, function() {});
	this.msgService.connect();
	this.locService.connect();
	this.msgService.jsonsend(addMsg);
	this.locService.jsonsend(addMsg);
	this.locService.jsonsend(locMsg);
	this.users = new LinkedList();
}

function initIntro() {
	setInterval(introFunc, framePause*10);
}

function introFunc() {
	users = "";
	userList.forEach(function(u) {users += "<li>"+JSON.stringify(u)+"</li>"});
	document.getElementById("player-list").innerHTML = users;
}

function handleMsg(msg) {
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
}
