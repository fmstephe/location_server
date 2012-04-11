var findPlayers = (function() {

	var nearbyUsers = new LinkedList();
	var connect;
	var committedToGame = false;
	var xPosMe, xPosYou;
	var initWind;
	var divs;
	var nickname = "anonymous";
	var tankGame;

	var locHandler = {
		handleLoc: function(loc) {
				   var op = loc.Op;
				   var usrInfo = loc;
				   if (op == "sAdd" || op == "sNearby" || op == "sVisible") {
					   usrInfo.isBusy = false;
					   nearbyUsers.append(usrInfo);
					   connect.sendMsg(usrInfo.Id, new NameReq());
					   connect.sendMsg(usrInfo.Id, new BusyReq());
				   } else if (op == "sRemove" || op == "sNotVisible") {
					   if (tankGame && idYou == usrInfo.Id) {
						   tankGame.kill();
						   escapeGame();
					   }
					   var filtered = nearbyUsers.filter(function(u) {return usrInfo.Id == u.Id});
					   console.log("filtered " + filtered.size);
					   if (filtered.satOne(function(u) {return u.inviteSent;})) {
						   uncommitFromGame();
					   }
					   if (filtered.satOne(function(u) {return u.inviteRcv;})) {
						   uncommitFromGame();
					   }
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
					   if (startOp == "invite") {
						   if (committedToGame) {
							   // If the start msg is from the same person we are currently inviting this will cause deadlock
							   // Need to break the deadlock by ordering user-ids and breaking the tie
							   connect.sendMsg(from, mkEnaged());
						   } else {
							   xPosMe = msg.Content.xPosYou;
							   xPosYou = msg.Content.xPosMe;
							   divs = msg.Content.divs;
							   initWind = msg.Content.initWind;
							   nearbyUsers.forEach(function(u) {if (u.Id == from) u.inviteRcv = true});
							   refreshUsers();
						   }
					   }
					   if (startOp == "engaged") {
						   uncommitFromGame();
						   nearbyUsers.forEach(function(u) {if (u.Id == from) u.isBusy = true; u.inviteSent = false;});
						   refreshUsers();
					   }
					   if (startOp == "decline") {
						   uncommitFromGame();
						   nearbyUsers.forEach(function(u) {if (u.Id == from) {u.inviteSent = false; u.declined = true;}});
						   setTimeout(function(){nearbyUsers.forEach(function(u) {if (u.Id == from) {u.declined = false}}); refreshUsers();}, 2000);
						   refreshUsers();
					   }
					   if (startOp == "accept") {
						   playGameState();
						   var nickYou;
						   nearbyUsers.forEach(function(u) {if (u.Id == from) u.inviteSent = false;});
						   nearbyUsers.forEach(function(u) {if (u.Id == from) {nickYou = u.nick;}});
						   clearRequests();
						   tankGame = mkTankGame();
						   tankGame.init(idMe, from, nickname, nickYou, xPosMe, xPosYou, initWind, connect, divs, turnQHandler, escapeGame);
					   }
				   }
			   }
	}

	var busyMsgHandler = {
		handleMsg: function(msg) {
				   if (msg.Content.isBusyMsg) {
					   var from = msg.From;
					   nearbyUsers.forEach(function(u) {if (u.Id == from) u.isBusy = msg.Content.isBusy});
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

	var nameRespHandler = {
		handleMsg: function(msg) {
				   var from = msg.From;
				   if (msg.Content.isNameResp) {
					   nearbyUsers.forEach(function(u) {if (u.Id == from) u.nick = msg.Content.nick;});
					   refreshUsers();
				   }
			   }
	}

	var nameReqHandler = {
		handleMsg: function(msg) {
				   var from = msg.From;
				   if (msg.Content.isNameReq) {
					   connect.sendMsg(from, new NameResp(nickname));
				   }
			   }
	}

	function refreshUsers() {
		console.log("refresh users");
		var users = "";
		var count = 0;
		nearbyUsers.forEach(function(u) {if (u.nick) count++;});
		if (count > 0) {
			nearbyUsers.forEach(function(u) {if (u.nick) users += userLiLink(u);});
		} else {
			users = "<div class='player-column'><button class='activeButton' style='visibility: hidden'></button>There's nobody nearby :(</div>";
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
		uncommitFromGame();
		findPlayersState();
		refreshUsers();
	}

	function commitToGame() {
		committedToGame = true;
		nearbyUsers.forEach(function(u) {connect.sendMsg(u.Id, new BusyMsg(true));});
	}

	function uncommitFromGame() {
		committedToGame = false;
		nearbyUsers.forEach(function(u) {connect.sendMsg(u.Id, new BusyMsg(false));});
	}

	function clearRequests() {
		nearbyUsers.forEach(function(u) {if (u.inviteRcv) connect.sendMsg(u.Id, mkDecline());});
	}

	// Public functions
	return {
		main: function() {
			      nickname = document.getElementById('nick-input').value;
			      var locHandlers = new LinkedList();
			      var msgHandlers = new LinkedList();
			      locHandlers.append(locHandler);
			      msgHandlers.append(nameReqHandler);
			      msgHandlers.append(nameRespHandler);
			      msgHandlers.append(turnQHandler);
			      msgHandlers.append(startHandler);
			      msgHandlers.append(busyMsgHandler);
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
				initWind = windChange();
				commitToGame();
				nearbyUsers.forEach(function(u) {if (u.Id == idYou) u.inviteSent = true;});
				refreshUsers();
				connect.sendMsg(idYou, mkInvite(divs, xPosMe, xPosYou, initWind));
			},

		accept: function(otherId) {
				// Flush game q
				turnQHandler.q.clear();
				idYou = otherId;
				connect.sendMsg(otherId, mkAccept());
				commitToGame();
				playGameState();
				var nickYou;
				nearbyUsers.forEach(function(u) {if (u.Id == idYou) {nickYou = u.nick}});
				nearbyUsers.forEach(function(u) {if (u.Id == otherId) u.inviteRcv = false;});
				clearRequests();
				tankGame = mkTankGame();
				tankGame.init(idMe, idYou, nickname, nickYou, xPosMe, xPosYou, initWind, connect, divs, turnQHandler, escapeGame);
			},

		decline: function(otherId) {
				 connect.sendMsg(otherId, mkDecline());
				 nearbyUsers.forEach(function(u) {if (u.Id == otherId) u.inviteRcv = false;});
				 refreshUsers();
			 }
	}
})();
