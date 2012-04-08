// Message used to negotiate the start of a game
function StartMsg(startOp, divs, xPosMe, xPosYou, initWind) {
	this.isStartMsg = true;
	this.startOp = startOp;
	this.divs = divs;
	this.xPosMe = xPosMe;
	this.xPosYou = xPosYou;
	this.initWind = initWind;
}

function mkInvite(divs, xPosMe, xPosYou, initWind) {
	return new StartMsg("invite", divs, xPosMe, xPosYou, initWind);
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
function PlayerMsg(player, newWind) {
	this.isPlayerMsg = true;
	this.player = player;
	this.newWind = newWind;
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
