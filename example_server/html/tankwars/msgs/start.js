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

