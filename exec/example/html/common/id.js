function getId() {
	var idReq = new XMLHttpRequest();
	idReq.open("GET", "http://battlewith.me.uk/id", false);
	idReq.send();
	idMsg = JSON.parse(idReq.responseText);
	return idMsg.id;
}
