function getId() {
	var idReq = new XMLHttpRequest();
	var url = "http://" + nakedURL();
	idReq.open("GET", url+"/id", false);
	idReq.send();
	idMsg = JSON.parse(idReq.responseText);
	return idMsg.id;
}
