// Message used to negotiate the start of a game
function StartMsg(startOp, defs) {
	this.isStartMsg = true;
	this.startOp = startOp;
	this.defs = defs;
}

function mkInvite(defs) {
	return new StartMsg("start", defs);
}

function mkAccept() {
	return new StartMsg("accept");
}

function mkEngaged() {
	return new StartMsg("engaged");
}

function mkDecline() {
	return new StartMsg("decline");
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

function BusyReq() {
	this.isBusyReq = true;
}
