function getId() {
	var idReq = new XMLHttpRequest();
	idReq.open("GET", "http://178.79.176.206:8001/common/id", false);
	idReq.send();
	idMsg = JSON.parse(idReq.responseText);
	return idMsg.Id;
}
