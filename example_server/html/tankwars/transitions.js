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
	document.onkeypress = nickEnter;
	$('intro-div').style.display='inline';
	$('nick-div').style.display='inline';
	$('player-div').style.display='none';
	$('game-div').style.display='none';
	$('nick-input').focus();
}

function findPlayersState() {
	document.onkeypress = null;
	$('intro-div').style.display='inline';
	$('nick-div').style.display='none';
	$('player-div').style.display='inline';
	$('game-div').style.display='none';
}

function playGameState() {
	document.onkeypress = null;
	$('intro-div').style.display='none';
	$('nick-div').style.display='none';
	$('player-div').style.display='none';
	$('game-div').style.display='inline';
}
