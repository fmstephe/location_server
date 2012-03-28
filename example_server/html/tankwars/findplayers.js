var findPlayers = (function() {

	var phrases = ["Kill", "Destroy", "Annihilate", "Devastate", "Massacre", "Defeat", "Bludgen", "Explode", "Defile", "Humiliate", "Crush", "Murder", "Smash", "Assassinate"];
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
					   usrInfo.buttonClass = "activebutton";
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
							   connect.sendMsg(from, mkStartEnaged());
						   } else {
							   connect.sendMsg(from, mkStartAccept());
							   idYou = from;
							   xPosMe = msg.Content.defs.xPosYou;
							   xPosYou = msg.Content.defs.xPosMe;
							   divs = msg.Content.defs.divs;
							   gameStarted = true;
							   playGameState();
							   selectionUsers.forEach(function(u) {if (u.Id != from){connect.sendMsg(u.Id, new BusyMsg(true));}});
							   tankGame().initGame(idMe, from, xPosMe, xPosYou, connect, divs, turnQHandler);
						   }
					   }
					   if (startOp == "engaged") {
						   gameStarted = false;
						   selectionUsers.forEach(function(u) {if (u.Id == from) {u.buttonClass = 'busyButton'}});
						   refreshPlayers();
					   }
					   if (startOp == "accept") {
						   playGameState();
						   selectionUsers.forEach(function(u) {if (u.Id != from){connect.sendMsg(u.Id, new BusyMsg(true));}});
						   tankGame().initGame(idMe, from, xPosMe, xPosYou, connect, divs, turnQHandler);
					   }
				   }
			   }
	}

	var busyHandler = {
		handleMsg: function(msg) {
				   if (msg.Content.isBusyMsg) {
					   var from = msg.From;
					   if (msg.Content.isBusy) {
						   selectionUsers.forEach(function(u) {if (u.Id == from) {u.buttonClass = "busybutton"}});
					   } else {
						   selectionUsers.forEach(function(u) {if (u.Id == from) {u.buttonClass = "activebutton"}});
					   }
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
					   selectionUsers.forEach(function(u) {if (u.Id == from) {u.nick = msg.Content.nick; u.phrase = phrases[r(phrases.length-1)];}});
				   }
				   refreshUsers();
			   }
	}

	var refreshUsers = function() {
		users = "";
		selectionUsers.forEach(function(u) {users += userLiLink(u)});
		document.getElementById("player-list").innerHTML = users;
	}

	var userLiLink = function(usr) {
		return "<li><button class='"+usr.buttonClass+" left-text-button' width='50px' type='button' onclick=\"findPlayers.startGame('"+usr.Id+"')\">"+formatNick(usr.nick)+"</button></li>";
	}

	function formatNick(nick) {
		var lim = 13;
		if (nick.length > lim) {
			return nick.substring(0,lim-3) + "..."
		} else {
			return nick;
		}
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
	}
})();
