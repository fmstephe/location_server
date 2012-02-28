var selectionUsers = new LinkedList();
var connect;
var gameStarted = false;
var idMe, idYou;
var xPosMe, xPosYou;
var divs;
// Constant Globals
var canvasHeight;
var canvasWidth;
var fgCtxt;
var terrainCtxt;
var bgCtxt;

var svcHandler = {
	handleLoc: function(loc) {
		var op = loc.Op;
		var usrInfo = loc.Msg;
		if (op == "sAdd" || op == "sNearby" || op == "sVisible") {
			selectionUsers.append(usrInfo);
		} else if (op == "sRemove" || op == "sNotVisible") {
			selectionUsers.filter(function(u) {return usrInfo.Id == u.Id});
		}
		users = "";
		selectionUsers.forEach(function(u) {users += userLiLink(u)});
		document.getElementById("player-list").innerHTML = users;
	},
	handleMsg: function(msg) {
		var from = msg.Msg.From;
		var content = msg.Msg.Content;
		if (content.op == "start") {
			if (gameStarted) {
				connect.sendMsg(new Msg(from, {op: "engaged"}));
			} else {
				connect.sendMsg(new Msg(from, {op: "accepted"}));
				idYou = from;
				xPosMe = content.defs.pos[1];
				xPosYou = content.defs.pos[0];
				divs = content.defs.divs;
				gameStarted = true;
				initGame(xPosMe, xPosYou, divs);
			}
		}
		if (content.op == "engaged") {
			gameStarted = false;
		}
		if (content.op == "accepted") {
			initGame(xPosMe, xPosYou, divs);
		}
		if (content.op == "fire") {
			launchList.append({from: from, params: content.params});
		}
	}
}

function main() {
	var fgCanvas = document.getElementById("foreground");
	var terrainCanvas = document.getElementById("terrain");
	var bgCanvas = document.getElementById("background");
	fgCtxt = fgCanvas.getContext("2d");
	terrainCtxt = terrainCanvas.getContext("2d");
	bgCtxt = bgCanvas.getContext("2d");
	canvasHeight = fgCanvas.height;
	canvasWidth = fgCanvas.width;
	var handlers = new LinkedList();
	handlers.append(svcHandler);
	connect = new Connect(handlers, handlers);
	idMe = connect.id;
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
	var msg = new Msg(idYou, {op: "start", defs: {divs: divs, pos: pair}});
	connect.sendMsg(msg);
	gameStarted = true;
}
