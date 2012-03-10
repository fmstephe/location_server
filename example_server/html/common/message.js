function ClientMessage(op, msg) {
	this.Op = op;
	this.Msg = msg;
}

function Add(id) {
	return new ClientMessage("cAdd", {Id: id});
}

function Msg(to, id, content) {
	return new ClientMessage("cMsg", {To: to, Id: id, Sends: 1, Content: content});
}

function ResendMsg(msg) {
	return new ClientMessage("cMsg", {To: msg.Msg.To, Id: msg.Msg.Id, Sends: msg.Msg.Sends+1, Content: msg.Msg.Content});
}

function Move(lat, lng) {
	return new ClientMessage("cMove", {Lat: lat, Lng: lng});
}

function InitLoc(lat, lng) {
	return new ClientMessage("cInitLoc", {Lat: lat, Lng: lng});
}

function Nearby(lat, lng) {
	return new ClientMessage("cNearby", {Lat: lat, Lng: lng});
}
