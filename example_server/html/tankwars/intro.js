var selectionUsers = new LinkedList();
var connect;

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
	}
}

function main() {
	connect = new Connect([svcHandler], [svcHandler]);
}

function userLiLink(user) {
	return "<li><a href=\"javascript:void(0)\" onclick=\"startGame('"+user.Id+"')\">"+JSON.stringify(user)+"</a></li>";
}

function startGame(id) {
	var msg = new Msg(id, "start!");
	connect.sendMsg(msg);
}
