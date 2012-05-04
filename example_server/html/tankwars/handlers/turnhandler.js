function TurnHandler(completeFun, idMe, idYou) {
	this.completeFun = completeFun;
	this.turnCount = 0;
	this.idMe = idMe;
	this.idYou = idYou;
	this.msgs = new LinkedList();
}

TurnHandler.prototype.isComplete = function() {
	var tm = this;
	var turnMsgs = new LinkedList();
	this.msgs.forEach(function(m) {if (m.content.turnCount == tm.turnCount) turnMsgs.append(m);});
	this.complete = this.completeFun(turnMsgs);
	return this.complete;
}

TurnHandler.prototype.getTurn = function() {
	if (!this.isComplete()) {
		return null;
	}
	var tm = this;
	var turn = this.msgs.filter(function(msg) {return msg.content.turnCount == tm.turnCount});
	this.turnCount++;
	return turn;
}

TurnHandler.prototype.handleMsg = function(msg) {
	if (msg.content.isTurnMsg && (msg.from == this.idYou || msg.from == this.idMe)) {
		this.msgs.append(msg);
	}
}
