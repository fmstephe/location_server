var selectionUsers = new LinkedList();
var gameStarted = false;
var xPosMe, xPosYou;
var divs;

var locHandler = {
	handleLoc: function(loc) {
		var op = loc.Op;
		var usrInfo = loc;
		if (op == "sAdd" || op == "sNearby" || op == "sVisible") {
			selectionUsers.append(usrInfo);
		} else if (op == "sRemove" || op == "sNotVisible") {
			selectionUsers.filter(function(u) {return usrInfo.Id == u.Id});
		}
		users = "";
		selectionUsers.forEach(function(u) {users += userLiLink(u)});
		document.getElementById("player-list").innerHTML = users;
	}
};

var turnQHandler = new QHandler(function(msg) {return msg.Content.isPlayerMsg;});

function mkStartHandler(connect) {
	return {
		handleMsg: function(msg) {
			if (msg.Content.isStart) {
				var from = msg.From;
				var startOp = msg.Content.startOp;
				if (startOp == "start") {
					if (gameStarted) {
						// If the start msg is from the same person we are currently inviting this will cause deadlock
						// Need to break the deadlock by ordering user-ids and breaking the tie
						connect.sendMsg(from, startEnaged());
					} else {
						connect.sendMsg(from, startAccept());
						idYou = from;
						xPosMe = content.defs.pos[1];
						xPosYou = content.defs.pos[0];
						divs = content.defs.divs;
						gameStarted = true;
						initGame(idMe, from, xPosMe, xPosYou, divs, turnQHandler);
					}
				}
				if (startOp == "engaged") {
					gameStarted = false;
				}
				if (startOp == "accept") {
					initGame(idMe, from, xPosMe, xPosYou, divs, turnQHandler);
				}
			}
		}
	}
}

function main() {
	var connect = new Connect(handlers, handlers);
	var locHandlers = new LinkedList();
	locandlers.append(locHandler);
	var msgHandlers = new LinkedList();
	msgHandlers.append(turnQHandler);
	msgHandlers.append(mkStartHandler(connect));
	idMe = connect.usrId;
}

function userLiLink(user) {
	return "<li><a href=\"javascript:void(0)\" onclick=\"startGame('"+user.Id+"')\">"+JSON.stringify(user)+"</a></li>";
}

function startGame(otherId) {
	idYou = otherId;
	var pair = positionPair(canvasWidth);
	xPosMe = pair[0];
	xPosYou = pair[1];
	divs = genDivisors();
	connect.sendMsg(idYou, {op:"start", defs: {divs: divs, pos: pair}});
	gameStarted = true;
}
