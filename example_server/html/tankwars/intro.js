var selectionUsers = new LinkedList();
var connecter = new Connecter();

var locHandler = {
	handleLoc: function(loc) {
		var op = loc.Op;
		var usrInfo = loc.Msg;
		if (op == "sAdd" || op == "sNearby" || op == "sVisible") {
			userList.append(usrInfo);
		} else if (op == "sRemove" || op == "sNotVisible") {
			selectionUsers.filter(function(u) {return usrInfo.Id == u.Id});
		}
		users = "";
		userList.forEach(function(u) {users += "<li>"+JSON.stringify(u)+"</li>"});
		document.getElementById("player-list").innerHTML = users;
	}
	handleMsg: function(msg) {
		alert(msg);
	}
}
