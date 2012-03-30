var findPlayers = (function() {

	var phrases = ["Kill", "Destroy", "Annihilate", "Devastate", "Massacre", "Defeat", "Bludgen", "Explode", "Defile", "Humiliate", "Crush", "Murder", "Smash", "Assassinate"];
	var selectionUsers = new LinkedList();
	var connect;
	var gameStarted = false;
	var xPosMe, xPosYou;
	var divs;
	var nickname = "anonymous";
	var tankGame;

	var locHandler = {
		handleLoc: function(loc) {
				   var op = loc.Op;
				   var usrInfo = loc;
				   if (op == "sAdd" || op == "sNearby" || op == "sVisible") {
					   usrInfo.nick = "anonymous";
					   usrInfo.isBusy = false;
					   selectionUsers.append(usrInfo);
					   connect.sendMsg(usrInfo.Id, new NameReq());
					   connect.sendMsg(usrInfo.Id, new BusyReq());
				   } else if (op == "sRemove" || op == "sNotVisible") {
					   selectionUsers.filter(function(u) {return usrInfo.Id == u.Id});
					   refreshUsers();
				   }
			   }
	}

	var turnQHandler = new QHandler(function(msg) {return msg.Content.isPlayerMsg;});

	var startHandler = {
		handleMsg: function(msg) {
				   if (msg.Content.isStartMsg) {
					   var from = msg.From;
					   var startOp = msg.Content.startOp;
					   if (startOp == "start") {
						   if (gameStarted) {
							   // If the start msg is from the same person we are currently inviting this will cause deadlock
							   // Need to break the deadlock by ordering user-ids and breaking the tie
							   connect.sendMsg(from, mkEnaged());
						   } else {
							   xPosMe = msg.Content.defs.xPosYou;
							   xPosYou = msg.Content.defs.xPosMe;
							   divs = msg.Content.defs.divs;
							   selectionUsers.forEach(function(u) {if (u.Id == from) u.hasInvited = true});
							   refreshUsers();
						   }
					   }
					   if (startOp == "engaged") {
						   gameStarted = false;
						   selectionUsers.forEach(function(u) {if (u.Id == from) u.isBusy = true;});
						   refreshUsers();
					   }
					   if (startOp == "declined") {
						   gameStarted = false;
					   }
					   if (startOp == "accept") {
						   playGameState();
						   selectionUsers.forEach(function(u) {if (u.Id != from)connect.sendMsg(u.Id, new BusyMsg(true));});
						   tankGame = mkTankGame();
						   tankGame.initGame(idMe, from, xPosMe, xPosYou, connect, divs, turnQHandler);
					   }
				   }
			   }
	}

	var busyHandler = {
		handleMsg: function(msg) {
				   if (msg.Content.isBusyMsg) {
					   var from = msg.From;
					   selectionUsers.forEach(function(u) {if (u.Id == from) u.isBusy = msg.Content.isBusy});
					   refreshUsers();
				   }
			   }
	}

	var busyReqHandler = {
		handleMsg: function(msg) {
				   if (msg.Content.isBusyReq) {
					   var from = msg.From;
					   connect.sendMsg(from, new BusyMsg(gameStarted));
				   }
			   }
	}

	var nameHandler = {
		handleMsg: function(msg) {
				   var from = msg.From;
				   if (msg.Content.isNameReq) {
					   connect.sendMsg(from, new NameResp(nickname));
				   } else if (msg.Content.isNameResp) {
					   selectionUsers.forEach(function(u) {if (u.Id == from) u.nick = msg.Content.nick; u.phrase = phrases[r(phrases.length-1)];});
				   }
				   refreshUsers();
			   }
	}

	var refreshUsers = function() {
		var users = "";
		if (selectionUsers.length() > 0) {
			selectionUsers.forEach(function(u) {users += userLiLink(u)});
		} else {
			users = "<div class='player-column'>There's no one nearby to play with :(</div>";
		}
		document.getElementById("player-div").innerHTML = users;
	}

	var userLiLink = function(usr) {
		var inviteClass = usr.isBusy ? "busybutton" : "activebutton";
		var responseVis = usr.hasInvited ? "visible" : "hidden";
		var firstLine = "<button class='"+inviteClass+"' onclick=\"findPlayers.invite('"+usr.Id+"')\">Play!</button>"+usr.nick;
		var secondLine = "<button class='activebutton' onclick=\"findPlayers.accept('"+usr.Id+"')\">Accept</button> <button class='activebutton' onclick=\"findPlayers.decline('"+usr.Id+"')\">Decline</button>";
		return "<div class='player-column'><div>"+firstLine+"</div><div style='visibility: "+responseVis+"'>"+secondLine+"</div></div>";
	}

	// Public functions
	return {
		main: function() {
			      nickname = document.getElementById('nickname').value;
			      var locHandlers = new LinkedList();
			      var msgHandlers = new LinkedList();
			      locHandlers.append(locHandler);
			      msgHandlers.append(nameHandler);
			      msgHandlers.append(turnQHandler);
			      msgHandlers.append(startHandler);
			      msgHandlers.append(busyHandler);
			      msgHandlers.append(busyReqHandler);
			      connect = new Connect(msgHandlers, locHandlers);
			      console.log("User Id: "+connect.usrId);
			      idMe = connect.usrId;
			      refreshUsers();
		      },
		invite: function(otherId) {
				// Flush game q
				turnQHandler.q.clear();
				idYou = otherId;
				var terrainCanvas = document.getElementById("terrain");
				var pair = positionPair(terrainCanvas.width);
				xPosMe = pair[0];
				xPosYou = pair[1];
				divs = genDivisors();
				connect.sendMsg(idYou, mkInvite({divs: divs, xPosMe: xPosMe, xPosYou: xPosYou}));
				gameStarted = true;
			},
		accept: function(otherId) {
				// Flush game q
				turnQHandler.q.clear();
				idYou = otherId;
				gameStarted = true;
				connect.sendMsg(otherId, mkAccept());
				playGameState();
				selectionUsers.forEach(function(u) {if (u.Id != idYou) connect.sendMsg(u.Id, new BusyMsg(true));});
				tankGame = mkTankGame();
				tankGame.initGame(idMe, idYou, xPosMe, xPosYou, connect, divs, turnQHandler);
			},
		decline: function(otherId) {
				 connect.sendMsg(otherId, mkDecline());
				 selectionUsers.forEach(function(u) {if (u.Id == otherId) u.hasInvited = false;});
				 refreshUsers();
			 }
	}
})();
