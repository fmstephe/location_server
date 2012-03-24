var findPlayers = (function() {

	var selectionUsers = new LinkedList();
	var connect;
	var gameStarted = false;
	var xPosMe, xPosYou;
	var divs;
	var nickname = "anonymous";

	var locHandler = {
		handleLoc: function(loc) {
				   var op = loc.Op;
				   var usrInfo = loc;
				   if (op == "sAdd" || op == "sNearby" || op == "sVisible") {
					   usrInfo.nick = "anonymous";
					   selectionUsers.append(usrInfo);
					   connect.sendMsg(usrInfo.Id, new NameReq());
				   } else if (op == "sRemove" || op == "sNotVisible") {
					   selectionUsers.filter(function(u) {return usrInfo.Id == u.Id});
				   }
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
							   connect.sendMsg(from, mkStartEnaged());
						   } else {
							   connect.sendMsg(from, mkStartAccept());
							   idYou = from;
							   xPosMe = msg.Content.defs.xPosYou;
							   xPosYou = msg.Content.defs.xPosMe;
							   divs = msg.Content.defs.divs;
							   gameStarted = true;
							   tankGame().initGame(idMe, from, xPosMe, xPosYou, connect, divs, turnQHandler);
						   }
					   }
					   if (startOp == "engaged") {
						   gameStarted = false;
					   }
					   if (startOp == "accept") {
						   tankGame().initGame(idMe, from, xPosMe, xPosYou, connect, divs, turnQHandler);
					   }
				   }
			   }
	}

	var nameHandler = {
		handleMsg: function(msg) {
				   var from = msg.From;
				   if (msg.Content.isNameReq) {
					   connect.sendMsg(from, new NameResp(nickname));
				   } else if (msg.Content.isNameResp) {
					   selectionUsers.forEach(function(u) {if (u.Id == from) {u.nick = msg.Content.nick;}});
					   users = "";
					   selectionUsers.forEach(function(u) {users += userLiLink(u)});
					   document.getElementById("player-list").innerHTML = users;
				   }
			   }
	}

	var userLiLink = function(usr) {
		return "<li><a href=\"javascript:void(0)\" onclick=\"findPlayers.startGame('"+usr.Id+"')\">"+usr.nick+"</a></li>";
	}

	// Public functions
	return {
		main: function() {
			      var locHandlers = new LinkedList();
			      var msgHandlers = new LinkedList();
			      locHandlers.append(locHandler);
			      msgHandlers.append(nameHandler);
			      msgHandlers.append(turnQHandler);
			      msgHandlers.append(startHandler);
			      connect = new Connect(msgHandlers, locHandlers);
			      idMe = connect.usrId;
		      },
		startGame: function(otherId) {
				   idYou = otherId;
				   var terrainCanvas = document.getElementById("terrain");
				   var pair = positionPair(terrainCanvas.width);
				   xPosMe = pair[0];
				   xPosYou = pair[1];
				   divs = genDivisors();
				   connect.sendMsg(idYou, new StartMsg("start", {divs: divs, xPosMe: xPosMe, xPosYou: xPosYou}));
				   gameStarted = true;
			   },
		setName: function(nick) {
				 nickname = nick;
			 }
	}
})();
