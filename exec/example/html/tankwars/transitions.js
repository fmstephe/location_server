function $() {
	return document.getElementById(arguments[0]);
}

function nickEnter(e) {
	if (e.keyCode == 13) {
		findPlayersState();
		findPlayers.main();
	}
}

function enterNickState() {
	if (!("WebSocket" in window) && !("MozWebsocket" in window)) {
		unsupportedState();
		return;
	}
	document.onkeypress = nickEnter;
	$('even-columns').style.display='block';
	$('game-columns').style.display='none';
	$('error-column').style.display='none';
	$('nick-input').focus();
	$('comment-div').innerHTML = nicknameText();
}

function findPlayersState() {
	document.onkeypress = null;
	$('even-columns').style.display='block';
	$('game-columns').style.display='none';
	$('error-column').style.display='none';
	$('comment-div').innerHTML = nearbyText();
}

function playGameState() {
	document.onkeypress = null;
	$('even-columns').style.display='none';
	$('game-columns').style.display='block';
	$('error-column').style.display='none';
	$('game-comment').innerHTML = gameText();
}

function disconnectState() {
	document.onkeypress = null;
	$('even-columns').style.display='none';
	$('game-columns').style.display='none';
	$('error-column').style.display='block';
	$('error-comment').innerHTML = disconnectText();
}

function unsupportedState() {
	document.onkeypress = null;
	$('even-columns').style.display='none';
	$('game-columns').style.display='none';
	$('error-column').style.display='block';
	$('error-comment').innerHTML = unsupportedText();
}
