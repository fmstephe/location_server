function disconnectText() {
	return "<div class='comment-heading'>Websocket Connection Closed</div>"
		+
		"<p>This may be because your proxy server doesn't support websockets or your internet connection may be struggling.</p>"
		+
		"<p><a href='javascript:location.reload(true)'>Refresh</a></p>";
}
