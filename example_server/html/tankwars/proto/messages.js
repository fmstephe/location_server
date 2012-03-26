// Message used to negotiate the start of a game
function StartMsg(startOp, defs) {
	this.isStartMsg = true;
	this.startOp = startOp;
	this.defs = defs;
}

function mkStartReq(defs) {
	return new StartMsg("start", defs);
}

function mkStartAccept() {
	return new StartMsg("accept");
}

function mkStartEngaged() {
	return new StartMsg("engaged");
}

// Message used to communicate a player's move each turn
function PlayerMsg(player) {
	this.isPlayerMsg = true;
	this.player = player;
}

function NameReq() {
	this.isNameReq = true;
}

function NameResp(nick) {
	this.isNameResp = true;
	this.nick = nick;
}

function BusyMsg(isBusy) {
	this.isBusyMsg = true;
	this.isBusy = isBusy;
}
