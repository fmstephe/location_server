var findPlayers = (function() {

	var selectionUsers = new LinkedList();
	var connect;
	var committedToGame = false;
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
					   if (tankGame && idYou == usrInfo.Id) {
						   tankGame.kill();
						   escapeGame();
					   }
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
						   if (committedToGame) {
							   // If the start msg is from the same person we are currently inviting this will cause deadlock
							   // Need to break the deadlock by ordering user-ids and breaking the tie
							   connect.sendMsg(from, mkEnaged());
						   } else {
							   xPosMe = msg.Content.defs.xPosYou;
							   xPosYou = msg.Content.defs.xPosMe;
							   divs = msg.Content.defs.divs;
							   selectionUsers.forEach(function(u) {if (u.Id == from) u.inviteRcv = true});
							   refreshUsers();
						   }
					   }
					   if (startOp == "engaged") {
						   committedToGame = false;
						   selectionUsers.forEach(function(u) {if (u.Id == from) u.isBusy = true; u.inviteSent = false;});
						   selectionUsers.forEach(function(u) {connect.sendMsg(u.Id, new BusyMsg(false));});
						   refreshUsers();
					   }
					   if (startOp == "decline") {
						   committedToGame = false;
						   selectionUsers.forEach(function(u) {if (u.Id == from) {u.inviteSent = false; u.declined = true;}});
						   selectionUsers.forEach(function(u) {connect.sendMsg(u.Id, new BusyMsg(false));});
						   setTimeout(function(){selectionUsers.forEach(function(u) {if (u.Id == from) {u.declined = false}}); refreshUsers();}, 2000);
						   refreshUsers();
					   }
					   if (startOp == "accept") {
						   playGameState();
						   var nickYou;
						   selectionUsers.forEach(function(u) {if (u.Id == from) u.inviteSent = false;});
						   selectionUsers.forEach(function(u) {if (u.Id == from) {nickYou = u.nick;}});
						   tankGame = mkTankGame();
						   tankGame.init(idMe, from, nickname, nickYou, xPosMe, xPosYou, connect, divs, turnQHandler, escapeGame);
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
					   connect.sendMsg(from, new BusyMsg(committedToGame));
				   }
			   }
	}

	var nameHandler = {
		handleMsg: function(msg) {
				   var from = msg.From;
				   if (msg.Content.isNameReq) {
					   connect.sendMsg(from, new NameResp(nickname));
				   } else if (msg.Content.isNameResp) {
					   selectionUsers.forEach(function(u) {if (u.Id == from) u.nick = msg.Content.nick;});
				   }
				   refreshUsers();
			   }
	}

	function refreshUsers() {
		var users = "";
		if (selectionUsers.length() > 0) {
			selectionUsers.forEach(function(u) {users += userLiLink(u)});
		} else {
			users = "<div class='player-column' style='text-align: center'>There's nobody nearby :(</div>";
		}
		document.getElementById("player-div").innerHTML = users;
	}

	function userLiLink(usr) {
		var inviteClass = usr.isBusy || committedToGame ? "busybutton" : "activebutton";
		var waitVis = usr.inviteSent ? "visible" : "hidden";
		var responseVis = usr.inviteRcv || usr.declined ? "visible" : "hidden";
		var waitGif =  "<img height='10' width='30' src='img/wait.gif' style='visibility: " + waitVis + "; margin-right:5px'>";
		var inviteButton = "<button class='" + inviteClass + "' onclick=\"findPlayers.invite('" + usr.Id + "')\">Invite</button>";
		var leftPad = "<span style='margin-left:35'>";
		var acceptButton = "<button class='activebutton' onclick=\"findPlayers.accept('" + usr.Id + "')\">Accept</button>";
		var declineButton = "<button class='activebutton' onclick=\"findPlayers.decline('" + usr.Id + "')\">Decline</button>";
		var declineMsg = "Invitation Declined :(<button class='activeButton' style='visibility: hidden'></button>";
		var secondLine = usr.declined ? leftPad + declineMsg : leftPad + acceptButton + declineButton;
		return "<div class='player-column'><div>" + waitGif + inviteButton + usr.nick + "</div><div style='visibility: " + responseVis + "'>" + secondLine + "</div></div>";
	}

	function escapeGame() {
		committedToGame = false;
		selectionUsers.forEach(function(u) {connect.sendMsg(u.Id, new BusyMsg(false));});
		findPlayersState();
		refreshUsers();
	}

	// Public functions
	return {
		main: function() {
			      nickname = document.getElementById('nick-input').value;
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
				if (committedToGame) {
					return;
				}
				// Flush game q
				turnQHandler.q.clear();
				idYou = otherId;
				var pair = positionPair(terrainWidth);
				xPosMe = pair[0];
				xPosYou = pair[1];
				divs = genDivisors();
				committedToGame = true;
				selectionUsers.forEach(function(u) {if (u.Id == idYou) u.inviteSent = true;});
				selectionUsers.forEach(function(u) {connect.sendMsg(u.Id, new BusyMsg(true));});
				refreshUsers();
				connect.sendMsg(idYou, mkInvite({divs: divs, xPosMe: xPosMe, xPosYou: xPosYou}));
			},

		accept: function(otherId) {
				// Flush game q
				turnQHandler.q.clear();
				idYou = otherId;
				committedToGame = true;
				connect.sendMsg(otherId, mkAccept());
				playGameState();
				selectionUsers.forEach(function(u) {connect.sendMsg(u.Id, new BusyMsg(true));});
				var nickYou;
				selectionUsers.forEach(function(u) {if (u.Id == idYou) {nickYou = u.nick}});
				selectionUsers.forEach(function(u) {if (u.Id == otherId) u.inviteRcv = false;});
				tankGame = mkTankGame();
				tankGame.init(idMe, idYou, nickname, nickYou, xPosMe, xPosYou, connect, divs, turnQHandler, escapeGame);
			},

		decline: function(otherId) {
				 connect.sendMsg(otherId, mkDecline());
				 selectionUsers.forEach(function(u) {if (u.Id == otherId) u.inviteRcv = false;});
				 refreshUsers();
			 }
	}
})();
