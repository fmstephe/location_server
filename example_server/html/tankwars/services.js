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
