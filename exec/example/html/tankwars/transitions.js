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
	if (navigator.appName == 'Microsoft Internet Explorer') {
		unsupportedState();
		return;
	}
	document.onkeypress = nickEnter;
	$('even-columns').style.display='block';
	$('game-columns').style.display='none';
	$('intro-div').style.display='block';
	$('player-div').style.display='none';
	$('game-div').style.display='none';
	$('disconnect-div').style.display='none';
	$('nick-input').focus();
	$('comment-div').innerHTML = nicknameText();
}

function findPlayersState() {
	document.onkeypress = null;
	$('even-columns').style.display='block';
	$('game-columns').style.display='none';
	$('intro-div').style.display='none';
	$('player-div').style.display='block';
	$('game-div').style.display='none';
	$('disconnect-div').style.display='none';
	$('comment-div').innerHTML = nearbyText();
}

function playGameState() {
	document.onkeypress = null;
	$('even-columns').style.display='none';
	$('game-columns').style.display='block';
	$('intro-div').style.display='none';
	$('player-div').style.display='none';
	$('game-div').style.display='block';
	$('disconnect-div').style.display='none';
	$('small-comment-div').innerHTML = gameText();
}

function disconnectState() {
	document.onkeypress = null;
	$('even-columns').style.display='none';
	$('game-columns').style.display='none';
	$('intro-div').style.display='none';
	$('player-div').style.display='none';
	$('game-div').style.display='none';
	$('disconnect-div').style.display='block';
	$('disconnect-div').innerHTML = disconnectText();
}

function unsupportedState() {
	document.onkeypress = null;
	$('even-columns').style.display='none';
	$('game-columns').style.display='none';
	$('intro-div').style.display='none';
	$('player-div').style.display='none';
	$('game-div').style.display='none';
	$('disconnect-div').style.display='block';
	$('disconnect-div').innerHTML = unsupportedText();
}
