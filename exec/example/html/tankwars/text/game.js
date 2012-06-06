function gameText() {
	return "<div class='comment-heading'>Battle!</div>"
		+
		"<p>The game itself is peer 2 peer, where player moves are communicated between browsers via the 'message' websocket.</p>"
		+
		"<p>Each browser waits for both turns to be sent and then plays them simultanously.</p>"
		+
		"<p>This means that I can run the game server with very little memory (because I have very little money).</p>"
		+
		"<p class='tech-text'>This was a lot of fun to make - race conditions everywhere I looked and plenty to learn</p>";
}
