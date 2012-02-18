function new Connecter(id, lat, lng, handleMsg, handleLoc) {
	var addMsg = new Add(id);
	var locMsg = new InitLoc(lat, lng);
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
