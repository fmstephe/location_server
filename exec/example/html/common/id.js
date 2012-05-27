function getId(port) {
	var idReq = new XMLHttpRequest();
	idReq.open("GET", "http://battlewith.me.uk:"+port+"/id", false);
	idReq.send();
	idMsg = JSON.parse(idReq.responseText);
	return idMsg.id;
}
