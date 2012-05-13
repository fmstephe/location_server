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
	$('intro-div').style.display='block';
	$('nick-div').style.display='block';
	$('player-div').style.display='none';
	$('game-div').style.display='none';
	$('control-div').style.display='none';
	$('nick-input').focus();
}

function findPlayersState() {
	document.onkeypress = null;
	$('intro-div').style.display='block';
	$('nick-div').style.display='none';
	$('player-div').style.display='block';
	$('game-div').style.display='none';
	$('control-div').style.display='none';
}

function playGameState() {
	document.onkeypress = null;
	$('intro-div').style.display='none';
	$('nick-div').style.display='none';
	$('player-div').style.display='none';
	$('game-div').style.display='block';
	$('control-div').style.display='block';
}
