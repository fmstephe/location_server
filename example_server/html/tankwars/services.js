function handleMsg(msg) {
	console.log(msg);
	var m = JSON.parse(msg.Msg);
	var p = new Player(m.x, m.name, turretLength, m.power, minPower, maxPower, powerInc, expRadius, null);
	playerList.append(p);
}

function handleLoc(msg) {
	var op = msg.Op;
	var usrInfo = msg.Usr;
	if (op == "sAdd" || op == "sNearby" || op == "sVisible") {
		playerMsg = new Msg(usrInfo.Id, JSON.stringify(new PlayerMsg(localPlayer)));
		msgService.jsonsend(playerMsg);
	} else if (op == "sMoved" || op == "sRemove" || op == "sNotVisible") {
		// Noop
	}
}
