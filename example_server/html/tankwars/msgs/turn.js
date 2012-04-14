// Message used to communicate a player's move each turn
function TurnMsg(turn, data) {
	this.isTurnMsg = true;
	this.turnCount = turn;
	this.data = data;
}
