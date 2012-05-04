var findPlayers = (function() {

	var nearbyUsers = new LinkedList();
	var connect;
	var committedToGame = false;
	var xPosMe, xPosYou;
	var initWind;
	var divs;
	var nickname = "anonymous";
	var tankGame;
	var height = 640;
	var width = 960;

	var locHandler = {
		handleLoc: function(loc) {
				   var op = loc.op;
				   var usrInfo = loc;
				   if (op == "sVisible") {
					   usrInfo.isBusy = false;
					   nearbyUsers.append(usrInfo);
					   connect.sendMsg(usrInfo.id, new NameReq());
					   connect.sendMsg(usrInfo.id, new BusyReq());
				   } else if (op == "sNotVisible") {
					   if (tankGame && idYou == usrInfo.id) {
						   tankGame.kill();
						   escapeGame();
					   }
					   var filtered = nearbyUsers.filter(function(u) {return usrInfo.id == u.id});
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

	var startHandler = {
		handleMsg: function(msg) {
				   if (msg.content.isStartMsg) {
					   var from = msg.from;
					   var startOp = msg.content.startOp;
					   if (startOp == "invite") {
						   if (committedToGame) {
							   // If the start msg is from the same person we are currently inviting this will cause deadlock
							   // Need to break the deadlock by ordering user-ids and breaking the tie
							   connect.sendMsg(from, mkEnaged());
						   } else {
							   xPosMe = msg.content.xPosYou;
							   xPosYou = msg.content.xPosMe;
							   divs = msg.content.divs;
							   initWind = msg.content.initWind;
							   nearbyUsers.forEach(function(u) {if (u.id == from) u.inviteRcv = true});
							   refreshUsers();
						   }
					   }
					   if (startOp == "engaged") {
						   uncommitFromGame();
						   nearbyUsers.forEach(function(u) {if (u.id == from) u.isBusy = true; u.inviteSent = false;});
						   refreshUsers();
					   }
					   if (startOp == "decline") {
						   uncommitFromGame();
						   nearbyUsers.forEach(function(u) {if (u.id == from) {u.inviteSent = false; u.declined = true;}});
						   setTimeout(function(){nearbyUsers.forEach(function(u) {if (u.id == from) {u.declined = false}}); refreshUsers();}, 2000);
						   refreshUsers();
					   }
					   if (startOp == "accept") {
						   playGameState();
						   var nickYou;
						   nearbyUsers.forEach(function(u) {if (u.id == from) u.inviteSent = false;});
						   nearbyUsers.forEach(function(u) {if (u.id == from) {nickYou = u.nick;}});
						   clearRequests();
						   tankGame = mkTankGame();
						   tankGame.init(height, width, idMe, from, nickname, nickYou, xPosMe, xPosYou, initWind, connect, divs, escapeGame);
					   }
				   }
			   }
	}

	var busyMsgHandler = {
		handleMsg: function(msg) {
				   if (msg.content.isBusyMsg) {
					   var from = msg.from;
					   nearbyUsers.forEach(function(u) {if (u.id == from) u.isBusy = msg.content.isBusy});
					   refreshUsers();
				   }
			   }
	}

	var busyReqHandler = {
		handleMsg: function(msg) {
				   if (msg.content.isBusyReq) {
					   var from = msg.from;
					   connect.sendMsg(from, new BusyMsg(committedToGame));
				   }
			   }
	}

	var nameRespHandler = {
		handleMsg: function(msg) {
				   var from = msg.from;
				   if (msg.content.isNameResp) {
					   nearbyUsers.forEach(function(u) {if (u.id == from) u.nick = msg.content.nick;});
					   refreshUsers();
				   }
			   }
	}

	var nameReqHandler = {
		handleMsg: function(msg) {
				   var from = msg.from;
				   if (msg.content.isNameReq) {
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
		var inviteFunc = usr.isBusy || committedToGame ? "function() {return 0;}" : "findPlayers.invite('" + usr.id + "');";
		var waitVis = usr.inviteSent ? "visible" : "hidden";
		var responseVis = usr.inviteRcv || usr.declined ? "visible" : "hidden";
		var waitGif =  "<img height='10' width='30' src='img/wait.gif' style='visibility: " + waitVis + "; margin-right:5px'>";
		var inviteButton = "<button class='" + inviteClass + "'onclick=\""+inviteFunc+"\">Invite</button>";
		var leftPad = "<span style='margin-left:35'>";
		var respondClass = committedToGame ? "busybutton" : "activebutton";
		var acceptFunc = committedToGame ? "function() {return 0;}" : "findPlayers.accept('" + usr.id + "');";
		var declineFunc = committedToGame ? "function() {return 0;}" : "findPlayers.decline('" + usr.id + "');";
		var acceptButton = "<button class='" + respondClass +"' onclick=\"" + acceptFunc + "\">Accept</button>";
		var declineButton = "<button class='" + respondClass +"' onclick=\"" + declineFunc + "\">Decline</button>";
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
		nearbyUsers.forEach(function(u) {connect.sendMsg(u.id, new BusyMsg(true));});
	}

	function uncommitFromGame() {
		committedToGame = false;
		nearbyUsers.forEach(function(u) {connect.sendMsg(u.id, new BusyMsg(false));});
	}

	function clearRequests() {
		nearbyUsers.forEach(function(u) {if (u.inviteRcv) { connect.sendMsg(u.id, mkDecline()); u.inviteRcv = false;}});
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
				idYou = otherId;
				var pair = positionPair(width);
				xPosMe = pair[0];
				xPosYou = pair[1];
				divs = genDivisors();
				initWind = windValue();
				commitToGame();
				nearbyUsers.forEach(function(u) {if (u.id == idYou) u.inviteSent = true;});
				refreshUsers();
				connect.sendMsg(idYou, mkInvite(divs, xPosMe, xPosYou, initWind));
			},

		accept: function(otherId) {
				idYou = otherId;
				connect.sendMsg(otherId, mkAccept());
				commitToGame();
				playGameState();
				var nickYou;
				nearbyUsers.forEach(function(u) {if (u.id == idYou) {nickYou = u.nick}});
				nearbyUsers.forEach(function(u) {if (u.id == otherId) u.inviteRcv = false;});
				clearRequests();
				tankGame = mkTankGame();
				tankGame.init(height, width, idMe, idYou, nickname, nickYou, xPosMe, xPosYou, initWind, connect, divs, escapeGame);
			},

		decline: function(otherId) {
				 connect.sendMsg(otherId, mkDecline());
				 nearbyUsers.forEach(function(u) {if (u.id == otherId) u.inviteRcv = false;});
				 refreshUsers();
			 }
	}
})();
