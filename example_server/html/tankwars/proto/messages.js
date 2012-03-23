// Message used to negotiate the start of a game
function StartMsg(startOp) {
	this.isStart = true;
	this.startOp = startOp;
}

function startReq() {
	return new StartMsg("start");
}

function startAccept() {
	return new StartMsg("accept");
}

function startEngaged() {
	return new StartMsg("engaged");
}

// Message used to communicate a player's move each turn
function PlayerMsg(player) {
	this.isPlayerMsg = true;
	this.player = player;
}
