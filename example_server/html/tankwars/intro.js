var selectionUsers = new LinkedList();
var connect;
var gameStarted = false;
var xPosMe, xPosYou;
var divs;
var nickname;

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
}

var turnQHandler = new QHandler(function(msg) {return msg.Content.isPlayerMsg;});

var startHandler = {
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
						   xPosMe = msg.Content.defs.xPosYou;
						   xPosYou = msg.Content.defs.xPosMe;
						   divs = msg.Content.defs.divs;
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

function main() {
	var locHandlers = new LinkedList();
	var msgHandlers = new LinkedList();
	locHandlers.append(locHandler);
	msgHandlers.append(turnQHandler);
	msgHandlers.append(startHandler);
	connect = new Connect(msgHandlers, locHandlers);
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
	connect.sendMsg(idYou, new StartMsg("start", {divs: divs, xPosMe: xPosMe, xPosYou: xPosYou}));
	gameStarted = true;
}
