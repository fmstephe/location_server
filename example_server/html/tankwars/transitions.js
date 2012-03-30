function $() {
	return document.getElementById(arguments[0]);
}

function enterNickState() {
	$('intro-div').style.display='inline';
	$('nick-div').style.display='inline';
	$('player-div').style.display='none';
	$('game-div').style.display='none';
}

function findPlayersState() {
	$('intro-div').style.display='inline';
	$('nick-div').style.display='none';
	$('player-div').style.display='inline';
	$('game-div').style.display='none';
}

function playGameState() {
	$('intro-div').style.display='none';
	$('nick-div').style.display='none';
	$('player-div').style.display='none';
	$('game-div').style.display='inline';
}
