var selectionUsers = new LinkedList();
var connect;
var gameStarted = false;
var myPos;
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
		alert(msg.Op + msg.Msg.From + msg.Msg.Content);
		var from = msg.Msg.From;
		var content = JSON.parse(msg.Msg.Content);
		if (content.op == "start") {
			if (gameStarted) {
				connect.sendMsg(new Msg(from, JSON.stringify({op:"engaged"})));
			} else {
				connect.sendMsg(new Msg(from, JSON.stringify({op:"accepted"})));
				myPos = content.defs.pos;
				divs = content.defs.divs;
				gameStarted = true;
				initGame(myPos, divs);
			}
		}
		if (content.op == "engaged") {
			gameStarted = false;
		}
		if (content.op == "accepted") {
			initGame(myPos, divs);	
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
	connect = new Connect([svcHandler], [svcHandler]);
}

function userLiLink(user) {
	return "<li><a href=\"javascript:void(0)\" onclick=\"startGame('"+user.Id+"')\">"+JSON.stringify(user)+"</a></li>";
}

function startGame(id) {
	var pair = positionPair(canvasWidth);
	myPos = pair[0];
	var oPos = pair[1];
	divs = genDivisors();
	var msg = new Msg(id, JSON.stringify({op: "start", defs: {divs: divs, pos: oPos}}));
	connect.sendMsg(msg);
	gameStarted = true;
}
