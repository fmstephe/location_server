function $() {
	return document.getElementById(arguments[0]);
}

function nickEnter(e) {
	if (e.keyCode == 13) {
		tryNick();
	}
}

function nickSelect() {
	var input = $("nick-input");
	if (input.value == "Enter a nickname here...") {
		input.value = "";
		input.className = "nick-input";
	}
}

function tryNick() {
	var input = $("nick-input");
	if (input.value == "Enter a nickname here..." || input.value == "") {
		if (document.activeElement != input) {
			input.value = "Enter a nickname here...";
			input.className = "nick-input-virgin";
		}
		return;
	} else {
		shareLocationState();
		findPlayers.main();
	}
}

function enterNickState() {
	if (!("WebSocket" in window) && !("MozWebsocket" in window)) {
		unsupportedState();
		return;
	}
	document.onkeypress = nickEnter;
	$('three-columns').style.display='block';
	$('game-columns').style.display='none';
	$('error-column').style.display='none';
}

function shareLocationState() {
	document.onkeypress = null;
	$('three-columns').style.display='block';
	$('game-columns').style.display='none';
	$('error-column').style.display='none';
	$('comment-div').innerHTML = shareLocationText();
}

function findPlayersState() {
	document.onkeypress = null;
	$('three-columns').style.display='block';
	$('game-columns').style.display='none';
	$('error-column').style.display='none';
	$('comment-div').innerHTML = nearbyText();
}

function playGameState() {
	document.onkeypress = null;
	$('three-columns').style.display='none';
	$('game-columns').style.display='block';
	$('error-column').style.display='none';
	$('game-comment').innerHTML = gameText();
}

function disconnectState() {
	document.onkeypress = null;
	$('three-columns').style.display='none';
	$('game-columns').style.display='none';
	$('error-column').style.display='block';
	$('error-comment').innerHTML = disconnectText();
}

function unsupportedState() {
	document.onkeypress = null;
	$('three-columns').style.display='none';
	$('game-columns').style.display='none';
	$('error-column').style.display='block';
	$('error-comment').innerHTML = unsupportedText();
}
